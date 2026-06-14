# Architecture

This document describes the system's structure: how requests flow through
it, how components are layered and connected, and the key design decisions
behind that structure.

## Table of Contents

1. [System Overview](#1-system-overview)
2. [Dependency Injection: the `application` Struct](#2-dependency-injection-the-application-struct)
3. [Request Pipeline](#3-request-pipeline)
4. [Data Layer: Repository Pattern](#4-data-layer-repository-pattern)
5. [Caching Layer: Cache-Aside](#5-caching-layer-cache-aside)
6. [Security: AuthN & AuthZ](#6-security-authn--authz)
7. [Resilience: Rate Limiting & Graceful Shutdown](#7-resilience-rate-limiting--graceful-shutdown)
8. [Database Design Decisions](#8-database-design-decisions)
9. [Observability](#9-observability)
10. [Testing Architecture](#10-testing-architecture)

---

## 1. System Overview

```
                         ┌─────────────────────────────────────┐
                         │              cmd/api                 │
                         │  (HTTP transport / composition root) │
                         └───────────────┬───────────────────────┘
                                          │
              ┌───────────────────────────┼───────────────────────────┐
              │                            │                            │
       ┌──────▼──────┐            ┌────────▼────────┐           ┌──────▼──────┐
       │ internal/    │            │ internal/store/  │           │ internal/   │
       │ store        │◄──cache────┤ cache             │           │ auth,       │
       │ (Postgres)   │   aside    │ (Redis)           │           │ ratelimiter,│
       └──────┬───────┘            └────────┬────────┘           │ mailer      │
              │                              │                     └─────────────┘
        ┌─────▼─────┐                  ┌─────▼─────┐
        │ PostgreSQL │                  │   Redis    │
        └────────────┘                  └────────────┘
```

The system is split into two layers:

- **`cmd/api`** — the HTTP transport layer (`package main`). Owns
  configuration, routing, middleware, and handlers. This is the
  "composition root": the only place that wires concrete implementations
  (Postgres store, Redis cache, JWT authenticator, rate limiter) together.
- **`internal/*`** — independently testable packages, each owning one
  concern (data access, caching, auth, rate limiting, email). Go's
  `internal/` mechanism enforces at compile time that these can only be
  imported from within this module, keeping them implementation details
  rather than a public API.

Every package boundary above is an **interface**. `cmd/api` never depends on
`*sql.DB` or `*redis.Client` directly — it depends on `store.Storage`,
`cache.CacheStorage`, `auth.Authenticator`, `ratelimiter.Limiter`. This is
what makes the next section possible.

---

## 2. Dependency Injection: the `application` Struct

```go
type application struct {
    config        config
    store         store.Storage
    cacheStorage  cache.CacheStorage
    logger        *slog.Logger
    mailer        mailer.Client
    authenticator auth.Authenticator
    rateLimiter   ratelimiter.Limiter
}
```

Every handler and middleware is a method on `*application`, giving it access
to every dependency through one receiver — no globals, no service locators.

`main.go` is the only place concrete types are constructed
(`store.NewPostgresStorage(db)`, `cache.NewRedisStorage(rdb)`,
`auth.NewJWTAuthenticator(...)`, `ratelimiter.NewFixedWindowLimiter(...)`).
Because `application`'s fields are interfaces, **tests build the same struct
with in-memory/mock implementations** instead
(`store.NewMockStore()`, `cache.NewMockStore()`, `auth.TestAuthenticator{}`)
and exercise the _real_ router (`app.mount()`) against them — see
[Testing Architecture](#10-testing-architecture).

---

## 3. Request Pipeline

`mount()` ([`cmd/api/api.go`](cmd/api/api.go)) builds the router as a fixed
pipeline of global middleware, followed by route-specific middleware:

```
Request
  │
  ▼
RequestID / RealIP        — tag the request, resolve real client IP
  │
  ▼
CORS                       — answer preflight OPTIONS, set CORS headers
  │
  ▼
Logger                     — structured access log
  │
  ▼
Recoverer                  — catch panics → 500 instead of crash
  │
  ▼
Timeout (60s)              — bound total request time
  │
  ▼
Rate Limiter (if enabled)  — 429 + Retry-After if over limit
  │
  ▼
Route-specific middleware  — e.g. authTokenMiddleware, checkPostOwnership
  │
  ▼
Handler
```

---

## 4. Data Layer: Repository Pattern

[`internal/store/store.go`](internal/store/store.go) defines the data access
contract as a set of interfaces grouped into `Storage`:

```go
type Storage struct {
    Posts     PostRepository
    Users     UserRepository
    Comments  CommentRepository
    Followers FollowerRepository
    Roles     RolesRepository
}
```

Handlers depend on these interfaces, not on Postgres-specific types. The
Postgres implementation lives in `postgres_*.go` files, using raw SQL (no
ORM) for full control over query shape, indexes, and joins.

A shared `DBTX` interface — satisfied by both `*sql.DB` and `*sql.Tx` — lets
the same query code run standalone _or_ inside a transaction. This is what
makes the [registration flow](#8-database-design-decisions) possible: user
creation and invitation-token creation run as one atomic unit, reusing the
same insert logic either way.

---

## 5. Caching Layer: Cache-Aside

```
                 ┌─────────┐   miss   ┌──────────┐
  getUser(id) ──►│  Redis  │─────────►│ Postgres │
                 └────┬────┘          └────┬─────┘
                      │ hit                │
                      │◄───────── Set ─────┘
                      ▼
                   response
```

`getUser()` ([`cmd/api/middleware.go`](cmd/api/middleware.go)) implements
cache-aside for user lookups: check Redis first; on a miss, read Postgres and
populate Redis (1-minute TTL) before returning. This sits in front of _every_
authenticated request, since `authTokenMiddleware` resolves the caller's user
record on each call.

Caching is controlled by a single `REDIS_ENABLED` flag. When disabled, the
cache layer is bypassed entirely and reads go straight to Postgres — the
feature can be turned on/off per environment (and off entirely in tests)
without code changes.

---

## 6. Security: AuthN & AuthZ

- **Authentication** is stateless JWT (HS256). `authTokenMiddleware`
  validates the bearer token's signature, issuer, audience, and expiration,
  then loads the user via the cache layer above. No server-side session
  store — any instance can validate any request, which is what makes
  horizontal scaling "just run more instances."
- **Authorization** is role-based: `user` / `moderator` / `admin`, each with
  a numeric precedence level. Post mutation handlers check two things —
  ownership (`post.UserID == user.ID`, always allowed) and role precedence
  (`user.Role.Level >= requiredRole.Level`, for non-owners).
- **CORS** (`go-chi/cors`) gates which browser origins may call the API at
  all, driven by `FRONTEND_URL`. Operational endpoints (`/v1/debug/vars`,
  `/v1/health`) use HTTP Basic Auth instead — a separate, simpler trust
  boundary for operators/monitoring tools rather than end users.

---

## 7. Resilience: Rate Limiting & Graceful Shutdown

**Rate limiting** protects the API from being overwhelmed by a single
client. It's defined as an interface —

```go
type Limiter interface {
    Allow(ip string) (bool, time.Duration)
}
```

— so the algorithm is swappable independent of the middleware that uses it.
The current implementation is a fixed-window counter per client IP; requests
over the limit get `429 Too Many Requests` with a `Retry-After` header and
never reach a handler.

**Graceful shutdown** protects in-flight requests during deploys/restarts.
`run()` ([`cmd/api/api.go`](cmd/api/api.go)) listens for `SIGINT`/`SIGTERM` in
a background goroutine; on receipt, it calls `http.Server.Shutdown` with a
10-second grace period — no new connections are accepted, but in-flight
requests are allowed to finish before the process exits. This is what lets
instances be spun up/down (scaling, rolling deploys) without dropping active
requests.

Both mechanisms exist for the same reason: **the API should degrade and
recover predictably under load and during operational changes**, rather than
failing abruptly.

---

## 8. Database Design Decisions

- **Optimistic concurrency control** — posts carry a `version` column.
  Updates are conditioned on `WHERE id = $1 AND version = $2`; if another
  update already bumped the version, zero rows match and the update is
  rejected (`ErrNotFound`) rather than silently overwriting a concurrent
  change. This avoids the "lost update" problem without holding row locks.
- **Feed algorithm** — a single joined query (`posts` ⋈ `followers` ⋈
  `users` ⋈ `comments`) returns "my posts OR posts from people I follow,"
  with the author's username and comment count attached — avoiding N+1
  queries.
- **Pagination & filtering** — offset/limit pagination with optional
  search (`ILIKE`), tag filtering (GIN-indexed array containment), and
  date-range filters, all parsed from query parameters.
- **Context timeouts** — every query gets a 5-second `context.WithTimeout`,
  guaranteeing a hung query can't hold a connection (or a request) forever.
- **Registration flow** — account creation and activation-token creation run
  in one transaction (via `DBTX`, see [Data Layer](#4-data-layer-repository-pattern)):
  an account should never exist without a way to activate it. The activation
  email is currently sent synchronously after commit — a known scaling
  limit; moving it to an async worker is a natural next step.

---

## 9. Observability

- **Structured logging** — all logs go through `slog` with a JSON handler,
  emitting consistent key-value fields (`method`, `path`, `error`) instead of
  free-form strings, so they're queryable in production.
- **Runtime metrics** via `expvar`, exposed at `/v1/debug/vars`: build
  version, `sql.DB.Stats()` (connection pool health), and goroutine count
  (a cheap leak indicator). Together with structured logs, these are the
  signals an operator would check first when something looks wrong.

---

## 10. Testing Architecture

The dependency-injection design from [§2](#2-dependency-injection-the-application-struct)
is what makes this possible: `newTestApplication(t, cfg)`
([`cmd/api/testutils_test.go`](cmd/api/testutils_test.go)) builds a real
`*application` backed by in-memory mocks (`store.NewMockStore()`,
`cache.NewMockStore()`, a test JWT authenticator), then calls `app.mount()` to
get the _actual_ router.

Tests fire real HTTP requests through `httptest` and assert on status codes
and mock call counts — exercising the full middleware pipeline (auth, rate
limiting, caching behavior) with no real database or Redis required.
`make test` runs the full suite (`go test ./... -v`).

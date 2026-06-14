# Social API (Go)

A social media backend API built in Go.

For a breakdown of architectural decisions ([ARCHITECTURE.md](ARCHITECTURE.md)).

---

## Features

- Users, Posts, Comments
- Follow / Unfollow
- Personalized feed with pagination, search, and tag filtering
- JWT authentication + role-based authorization (user / moderator / admin)
- Redis caching (cache-aside, toggleable via env flag)
- Fixed-window rate limiting with `Retry-After`
- CORS
- Graceful shutdown
- Structured logging (`slog`)
- Runtime metrics via `expvar`
- Swagger/OpenAPI docs

## Tech Stack

- Go + [chi](https://github.com/go-chi/chi) router
- PostgreSQL (raw SQL via `lib/pq`)
- Redis (`go-redis/v8`)
- JWT (`golang-jwt/jwt/v5`)
- SendGrid for transactional email
- `swaggo` for API docs
- `testify` for tests

---

## Getting Started

### 1. Start dependencies

```bash
docker-compose up -d
```

Starts Postgres, Redis, and Redis Commander (UI at `http://localhost:8081`).

### 2. Configure environment

Fill in `.env` as needed — defaults work for local dev with docker-compose.
A few flags worth knowing:

- `REDIS_ENABLED` — toggle Redis caching on/off
- `RATE_LIMITER_ENABLED` / `RATELIMITER_REQUESTS_COUNT` — toggle/tune rate limiting
- `FRONTEND_URL` — allowed CORS origin

### 3. Run migrations & seed data

```bash
make migrate-up
make seed
```

### 4. Run the API

```bash
go run cmd/api/main.go
# or, with live reload:
air
```

API docs: `http://localhost:8080/v1/swagger/index.html`

---

## Commands

| Command                                 | Description                        |
| --------------------------------------- | ---------------------------------- |
| `make test`                             | Run all tests (`go test ./... -v`) |
| `make migrate-up` / `make migrate-down` | Run / rollback DB migrations       |
| `make seed`                             | Seed the database with sample data |
| `make gen-docs`                         | Regenerate Swagger docs            |

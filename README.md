# Social API (Go)

A social media backend API built in Go to practice backend system design,
with key architectural decisions outlined below.

---

## Features

- Users, Posts, Comments
- Follow / Unfollow
- Personalized Feed

---

## Design Highlights

### Chi Router

- Uses chi router for lightweight routing
- Route grouping and middleware support
- Clean, modular endpoint structure

---

### Application Struct Pattern

- Central `application` struct holds dependencies (store, config)
- Handlers implemented as methods on `application`
- Enables dependency injection and cleaner code organization

---

### Repository Layer

- Abstracted database logic behind interfaces
- Separation: handlers → store → database
- Easily extendable to other databases

---

### PostgreSQL Database

- Uses raw SQL (no ORM)
- Fine-grained control over queries and performance
- Relational schema models relationships between users, posts, and comments
- Optimized with indexes and joins (e.g. indexed `username` for efficient user queries)

---

### Optimistic Concurrency Control

- `version` field on posts
- Prevents overwriting updates during concurrent writes

---

### Feed Algorithm

- Returns:
  - User’s own posts
  - Posts from followed users

- Implemented via SQL joins on `followers`

---

### Pagination & Filtering

- Offset/limit pagination
- Search via `ILIKE`
- Tag filtering with GIN indexes

---

### Context & Timeouts

- All DB operations use context timeouts
- Prevents hanging queries and improves reliability

---

## Planned Improvements

- Authentication / Authorization
- Redis caching (feed, hot data)
- Metrics & observability
- Rate limiting

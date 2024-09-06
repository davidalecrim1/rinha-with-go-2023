# Alternative (With Cache)
This branch is focused on testing new stack and compare the performance in Go.

## New Components
- [`fiber`](https://gofiber.io/) for REST API.
- Redis for Caching
- `ruedis` for Redis Driver

## Same Old
- Nginx as Reserve Proxy for Load Balacing
- Postgres as Database
- `pgx and pgxpool` for Postgres Driver
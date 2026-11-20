# zaps

A minimal URL shortener built with Go, PostgreSQL, and `sqlc` — usable from both the terminal (CLI) and a lightweight embedded web UI.

Live demo: [s.justinsep.me](https://s.justinsep.me)

## Features
- Shorten URLs with custom or auto-generated 6-character slugs
- Thread-safe in-memory redirect cache (`sync.RWMutex`)
- CLI client with automatic local machine key (`~/.zaps.json`)
- Embedded web UI (`go:embed`) bundled inside a single server binary
- PostgreSQL persistence with `sqlc` type-safe queries

## Install CLI
```bash
go install github.com/JUSTINS88/zaps/cmd/zaps@latest
```

## Project Structure
- `cmd/zapsd` — HTTP server (`POST /api/links`, `GET /{slug}`, `GET /stats`, `GET /`)
- `cmd/zaps` — CLI client
- `internal/web` — Embedded frontend (`index.html`, `stats.html`, `zaps-mark.svg`)
- `internal/store` — `sqlc` generated database layer
- `db/` — SQL migrations and queries
- `deploy/` — Dockerfile and Docker Compose config

## Running Locally

### With Docker Compose
```bash
docker compose -f deploy/docker-compose.yml up --build
```

### Manual
1. Start PostgreSQL and run migrations in `db/migrations/`.
2. Copy `.env.example` to `.env` and set `DB_URL`.
3. Run the server:
   ```bash
   go run ./cmd/zapsd
   ```
4. Use the CLI:
   ```bash
   zaps https://example.com -s my-slug
   ```

## API
- `POST /api/links` — Create a short link
  ```json
  {
    "original_link": "https://example.com",
    "slug": "optional-slug",
    "machine_key": "your-key"
  }
  ```
- `GET /{slug}` — Redirect to the original URL

## Benchmark (Local)

Tested locally against PostgreSQL using [`hey`](https://github.com/rakyll/hey):

```bash
hey -n 1000 -c 50 -disable-redirects http://localhost:8080/<slug>
```

| Metric | Before Cache (DB) | After Cache (`sync.RWMutex`) |
| --- | --- | --- |
| Throughput | ~2,210 req/s | **~9,110 req/s** |
| Median Latency (p50) | 2.0 ms | **1.3 ms** |
| p95 Latency | 78.0 ms | **22.1 ms** |
| Status | 1000/1000 (`302 Found`, 0 errors) | 1000/1000 (`302 Found`, 0 errors) |

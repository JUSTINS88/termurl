# termurl

A minimal URL shortener built with Go, PostgreSQL, and `sqlc` — usable from both the terminal (CLI) and a lightweight embedded web UI.

## Features
- Shorten URLs with custom or auto-generated 6-character slugs
- CLI client with automatic local machine key (`~/.termurl.json`)
- Embedded web UI (`go:embed`) bundled inside a single server binary
- PostgreSQL persistence with `sqlc` type-safe queries

## Project Structure
- `cmd/termurld` — HTTP server (`POST /api/links`, `GET /{slug}`, `GET /`)
- `cmd/termurl` — CLI client
- `internal/web` — Embedded frontend (`index.html`)
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
   go run ./cmd/termurld
   ```
4. Use the CLI:
   ```bash
   go run ./cmd/termurl https://example.com -s my-slug
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
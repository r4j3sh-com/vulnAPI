# Harbor Cloud Demo API (Go)

A deliberately vulnerable API lab written in Go, wrapped in a realistic storefront/ops landing page. Use it to practice API security testing and learn OWASP API Security Top 10 (2019 and 2023) issues in a more authentic setting.

> Warning: This app is intentionally insecure. Run it only in isolated, local environments.

## Quick start

```bash
docker compose up --build
```

Open the landing page at `http://localhost:9000`.

SQLite data is stored at `data/vulnapi.db` by default (or `/data/vulnapi.db` in Docker).

## Demo users

- alice / password123 (admin)
- bob / bobpass
- charlie / charlie

## Authentication (intentionally weak)

```bash
curl -s -X POST http://localhost:9000/api/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"bob","password":"bobpass"}'
```

Use the token:

```bash
curl -s http://localhost:9000/api/orders/1 \
  -H 'Authorization: Bearer <token>'
```

## Core endpoints

- `POST /api/auth/login`
- `POST /api/auth/reset`
- `GET /api/auth/dashboard`
- `GET /api/users`
- `GET /api/users/{id}`
- `PUT /api/users/{id}`
- `GET /api/users/{id}/orders`
- `GET /api/orders`
- `GET /api/orders/{id}`
- `POST /api/orders`
- `GET /api/catalog`
- `GET /api/catalog/search?q=...`
- `GET /api/reports?size=...`
- `GET /api/ops/metrics`
- `POST /api/ops/roles/promote/{id}`
- `GET /api/ops/ping?host=...`
- `GET /api/ops/files?path=...`
- `GET /api/v1/users`
- `GET /api/v2/users`
- `GET /ops/config`

## Training guide

The OWASP API Top 10 mapping (2019 and 2023) and sample attack flows are documented in `LAB_GUIDE.md`.

## Configuration

Environment variables:

- `PORT` (default: 9000)
- `APP_ENV` (default: dev)
- `DEBUG` (default: true)
- `TOKEN_SECRET` (unused, kept for misconfiguration training)
- `DB_PATH` (default: data/vulnapi.db)

## Project structure

- `cmd/vulnapi/main.go` - entry point
- `internal/app` - handlers, auth helpers, sqlite store
- `internal/app/web/` - landing page templates and static assets (embedded)
- `Dockerfile`, `docker-compose.yml` - container setup

## Disclaimer

This project is for educational and authorized testing only. Do not deploy or expose it to untrusted networks.

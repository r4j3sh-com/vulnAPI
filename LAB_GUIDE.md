# Harbor Cloud API Lab Guide

This guide maps the OWASP API Security Top 10 (2019 and 2023) to endpoints in this lab. Use it for structured practice.

## OWASP API Top 10 (2019)

### API1: Broken Object Level Authorization
- `GET /api/users/{id}`
- `GET /api/orders/{id}`
- `GET /api/users/{id}/orders`

### API2: Broken Authentication
- `POST /api/auth/login` (weak tokens, no signature or expiry)
- `POST /api/auth/reset` (no verification or rate limit)

### API3: Excessive Data Exposure
- `GET /api/users`
- `GET /api/v1/users`
- `GET /api/catalog`

### API4: Lack of Resources & Rate Limiting
- `GET /api/reports?size=1000000`

### API5: Broken Function Level Authorization
- `GET /api/ops/metrics`
- `POST /api/ops/roles/promote/{id}`

### API6: Mass Assignment
- `POST /api/users`
- `PUT /api/users/{id}`

### API7: Security Misconfiguration
- `GET /ops/config`
- `GET /api/ops/files?path=/etc/hosts`

### API8: Injection
- `GET /api/catalog/search?q=' OR 1=1 --`
- `GET /api/ops/ping?host=127.0.0.1;id`

### API9: Improper Assets Management
- `GET /api/v1/users`

### API10: Insufficient Logging & Monitoring
- Minimal logging and no audit trails for auth/privileged actions

## OWASP API Top 10 (2023)

### API1: Broken Object Level Authorization
- `GET /api/users/{id}`
- `GET /api/orders/{id}`
- `GET /api/users/{id}/orders`

### API2: Broken Authentication
- `POST /api/auth/login` (weak tokens, no signature or expiry)
- `POST /api/auth/reset` (no verification or rate limit)

### API3: Broken Object Property Level Authorization
- `GET /api/users` (returns password, api_key)
- `GET /api/v1/users` (legacy response exposes secrets)
- `GET /api/catalog` (internal cost, supplier, tags)

### API4: Unrestricted Resource Consumption
- `GET /api/reports?size=1000000`

### API5: Broken Function Level Authorization
- `GET /api/ops/metrics`
- `POST /api/ops/roles/promote/{id}`

### API6: Unrestricted Access to Sensitive Business Flows
- `POST /api/auth/reset` (no friction or verification)
- `POST /api/orders` (no anti-automation checks)

### API7: Server Side Request Forgery
- `GET /api/ops/ping?host=127.0.0.1` (host control, can target internal hosts)

### API8: Security Misconfiguration
- `GET /ops/config`
- `GET /api/ops/files?path=/etc/hosts`

### API9: Improper Inventory Management
- `GET /api/v1/users` (legacy version still deployed)

### API10: Unsafe Consumption of APIs
- API responses are trusted and surfaced directly without validation in multiple endpoints (e.g., `/api/catalog/search`), mirroring unsafe downstream consumption patterns.

## Sample requests

```bash
# BOLA: read another user
curl -s http://localhost:9000/api/users/1

# Mass assignment: create admin user
curl -s -X POST http://localhost:9000/api/users \
  -H 'Content-Type: application/json' \
  -d '{"username":"evil","password":"p@ss","email":"evil@harbor.example","isAdmin":true}'

# BFLA: promote without admin role
curl -s -X POST http://localhost:9000/api/ops/roles/promote/2 \
  -H 'Authorization: Bearer <token>'

# Injection: command injection
curl -s 'http://localhost:9000/api/ops/ping?host=127.0.0.1;id'

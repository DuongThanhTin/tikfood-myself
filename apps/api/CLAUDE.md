# CLAUDE.md — apps/api (Go Backend)

Read `/CLAUDE.md` first. This file describes the backend as it is built today.

## Stack Reality

- Go + **Gin**.
- **`database/sql` via the pgx stdlib driver** — NOT GORM.
- **slog** JSON logging; manual dependency injection via `internal/app/Container`.
- PostGIS-backed Postgres; an in-memory fallback repository runs when
  `DATABASE_URL` is unset.

## Layering

`http (Handler) → discovery (Service) → Repository → postgres`. Dependencies
point one way only.

Forbidden: Handler → repository, Handler → SQL, Service → `gin.Context`,
Repository → business logic.

## Handler (`internal/http`)

Handlers are named `*Handler` (e.g. `VenueHandler`) — not "Controller". They:
parse and validate the request (`venue_request.go`), call the service, and
respond with the `respondWith*` helpers. Handlers never contain business logic,
touch the database, or manage transactions.

## Service (`internal/discovery`)

Services hold all business logic and input normalization and call repositories.
Services never import or reference HTTP / `gin.Context`.

## Repository

`VenueRepository` is an interface with two implementations: `DiscoveryRepository`
in `storage/postgres` and the in-memory fallback in `discovery`. Repositories do
persistence only — no business rules, no HTTP.

## Response Envelope (real)

`{ "data": ..., "error": { "code", "message", "details?" } }`.
Use `respondWithData`, `respondWithError`, `respondWithBadRequest`,
`respondWithInternalServerError`. Do NOT invent a `{ "success": ... }` format,
and do not add `meta` or `trace_id` fields that do not exist yet.

## Error Handling

Use the predefined codes in `internal/http/errors.go` (`invalid_request`,
`not_found`, `internal_error`). Map domain errors (e.g. `ErrVenueNotFound`) to
HTTP status. Never leak SQL or raw error strings to the client.

## Validation

Query/DTO validation lives in the handler (`venue_request.go`). Business rules
live in the service.

## Config & DI

Configuration comes from `internal/config` (env). Never hardcode URLs, secrets,
or credentials. Use constructor injection through `Container`; no globals.

## Testing

Table-driven tests; follow `internal/http/router_test.go`. Mock
`VenueRepository` for service-level tests.

## Future Direction (approved, do not do preemptively)

- Add `trace_id` to the response envelope (`requestIDMiddleware` already exists).
- Consolidate error codes into a documented catalog as they grow.
- Transactions live in the Service layer once write endpoints exist (discovery
  is read-only today).

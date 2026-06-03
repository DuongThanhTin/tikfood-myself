# Backend Architecture Standard

Backend code lives in `apps/api`.

TikFood backend is a Go service for realtime social food discovery. It must stay focused on map discovery, dish discovery, trend scoring, social proof, AI summaries, hidden gems, and geo intelligence.

## Current Backend Baseline

- Language: Go
- HTTP: standard library `net/http`
- Entrypoint: `apps/api/cmd/server/main.go`
- Routes: `apps/api/internal/http`
- Discovery domain: `apps/api/internal/discovery`
- Persistence: in-memory MVP data for now

## Layering

Use this direction:

```text
cmd/server
-> internal/http
-> internal/<domain>
-> internal/storage
-> database or external systems
```

Responsibilities:

- `cmd/server`: process startup, port, server wiring.
- `internal/http`: routes, request parsing, response writing, status codes.
- `internal/<domain>`: domain models and business rules.
- `internal/storage`: database access once persistence exists.

HTTP handlers should stay thin. If a handler starts doing filtering, scoring, persistence, or summarization logic, move that behavior into a domain service.

## Domain Boundaries

Current and future domains:

- `discovery`: map venues, dish search, discovery feed
- `trend`: trend scoring and ranking
- `summary`: AI-generated venue/dish summaries
- `ingestion`: social signal ingestion
- `geo`: geospatial filtering and ranking

Keep these separate. Do not blend trend scoring, AI summary generation, ingestion, and HTTP transport code.

## API Rules

- Use JSON request and response bodies.
- Use explicit typed structs.
- Use `snake_case` JSON fields for API responses.
- Return a consistent envelope:

```json
{
  "data": {},
  "error": ""
}
```

- Omit `error` when there is no error.
- Use proper HTTP status codes.

## Error Handling

- Do not panic for request-level errors.
- Validate request input at the transport boundary.
- Return clear error messages without leaking internals.
- Do not include secrets, tokens, SQL details, or stack traces in API responses.

## Testing

Required for backend changes:

```bash
npm run api:test
```

Add tests for:

- New endpoints
- Filtering behavior
- Domain service behavior
- Error cases

## Persistence Direction

When persistence is added:

- Use PostgreSQL + PostGIS.
- Keep migrations explicit.
- Put database code under `internal/storage`.
- Keep SQL and storage models separate from HTTP response structs when complexity grows.

## Anti-Goals

Do not add backend logic for:

- Delivery
- Cart
- Orders
- Checkout
- Payment
- Booking or reservations
- In-app chat
- Social follow graph
- Creator monetization
- Livestream

# Service Contract — `apps/api` (Discovery API)

> **Why this doc exists:** a single at-a-glance contract for the Go discovery API — its
> surface, dependencies, consumers, and change rules — so a change here is made with
> full knowledge of who depends on it. Binding rules live in
> [`apps/api/CLAUDE.md`](../../apps/api/CLAUDE.md) and
> [`docs/standards/backend/`](../standards/backend/); this doc links, it does not restate.

## Purpose

Serve read-only food-discovery data (venues, dishes, social videos, geo) to the web
frontend. Discovery only — no writes, no anti-goal features. Entity vocabulary behind the
fields: [`docs/tikfood/domain-model.md`](../tikfood/domain-model.md).

## Interface

**Base:** `/api/v1` · **Envelope:** `{ data, error }`
([ADR-0003](../adr/0003-data-error-response-envelope.md)).

| Method & path | Handler | Purpose |
| --- | --- | --- |
| `GET /health` | inline | Liveness. Returns `{ "ok": true }` (does not ping DB yet) |
| `GET /api/v1/discovery/venues` | `Search` | Venue search/list with filters |
| `GET /api/v1/map/venues` | `Search` | **Alias** of the above (map surface) |
| `GET /api/v1/discovery/venues/:slug` | `Detail` | Single venue by slug |

**Search query parameters** (validated in `internal/http/venue_request.go`):

| Param | Type | Rules |
| --- | --- | --- |
| `q` | string | ≤ 120 chars |
| `city` | string | ≤ 80 chars |
| `district` | string | ≤ 80 chars |
| `dish` | string | free text |
| `tags` | CSV | comma-separated |
| `platform` | CSV | each of `tiktok\|instagram\|youtube\|facebook\|other` |
| `lat`,`lng` | float | must be provided together; lat ∈ [-90,90], lng ∈ [-180,180] |
| `radius_m` | int | 0–50000 |
| `min_price_vnd`,`max_price_vnd` | int | ≥ 0; `min` ≤ `max` |
| `open_now` | bool | parseable boolean |
| `sort` | enum | `trending`(default) `\| videos \| distance \| price`; `distance` requires `lat`+`lng` |
| `limit` | int | 1–100 (default 50) |

Invalid input → `400` with `error.code = invalid_request`. Codes are defined in
`internal/http/errors.go` (`invalid_request`, `not_found`, `internal_error`); raw SQL/
errors are never leaked. Full request/response detail:
[`docs/standards/api-contracts.md`](../standards/api-contracts.md),
[`docs/standards/backend/request-response.md`](../standards/backend/request-response.md).

## Dependencies

- **PostgreSQL/PostGIS** for the real path; **in-memory fallback** when `DATABASE_URL`
  is unset ([ADR-0004](../adr/0004-in-memory-fallback-repository.md)). Distance
  (`sort=distance`) and `dish` substring matching require the Postgres path — the
  fallback diverges (see ADR-0004 / `IMPROVEMENTS.md`).
- **Configuration** via `internal/config` (env only; no hardcoded URLs/secrets).
- Framework/driver: Gin + `database/sql`/pgx
  ([ADR-0001](../adr/0001-gin-and-pgx-not-gorm.md)).

## Consumers

- **`apps/web`** via [`apps/web/lib/api.ts`](../../apps/web/lib/api.ts) — the only
  client. Frontend types mirror the Go JSON tags (snake_case).
- **n8n / external** may call `GET /health`.

## Change rules

- **No breaking changes** to the envelope or existing fields ([ADR-0003](../adr/0003-data-error-response-envelope.md)).
- Any response-shape change **must** be mirrored in `apps/web/lib/api.ts`.
- Keep validation in the handler, business rules in the service, persistence in the
  repository ([ADR-0002](../adr/0002-handler-service-repository-layering.md)).
- New error codes: add to `errors.go` and map domain errors to HTTP status.
- Adding write endpoints, auth, or migrations → **human approval**
  ([`docs/ai/AI-CONTRACT.md`](../ai/AI-CONTRACT.md)).
- Update tests when behavior changes (`internal/http/router_test.go` is the pattern).

## Current state

Read-only; no auth; no pagination `meta`; no `trace_id` in the envelope
(`requestIDMiddleware` exists, wiring deferred); `GET /health` is shallow. Venue
ingestion (`internal/ingest`, `cmd/ingest`) exists as a CLI but is **not** exposed over
HTTP. Trend-score/AI-summary data is not populated yet. See
[`docs/architecture.md`](../architecture.md) §2 and [`IMPROVEMENTS.md`](../../IMPROVEMENTS.md).

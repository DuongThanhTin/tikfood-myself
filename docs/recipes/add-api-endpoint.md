# Recipe — Add / extend an API endpoint (`apps/api`)

**Goal:** add or extend a discovery endpoint without breaking the contract or the
layering.

**When to use:** new `GET` route, new query filter, or new field on an existing venue
response.

## Context to load
- [`apps/api/CLAUDE.md`](../../apps/api/CLAUDE.md) · [`docs/services/api.md`](../services/api.md)
- [`docs/standards/backend/patterns.md`](../standards/backend/patterns.md),
  [`request-response.md`](../standards/backend/request-response.md),
  [`errors.md`](../standards/backend/errors.md)
- ADRs [0002](../adr/0002-handler-service-repository-layering.md) (layering),
  [0003](../adr/0003-data-error-response-envelope.md) (envelope)

## Process skill
`brainstorming` (if non-trivial) → `writing-plans` → `test-driven-development`.

## Steps
1. **Confirm it's in scope** — discovery-only, not an anti-goal. Read-only unless writes
   are explicitly approved (writes/auth = human approval).
2. **Design the contract first** — path, query params (with validation rules + limits),
   response fields. Reuse the `{data,error}` envelope; do not invent new shapes.
3. **Handler** (`internal/http`) — parse/validate in the request parser (follow
   `venue_request.go`), call the service, respond with `respondWith*`. No SQL, no logic.
4. **Service** (`internal/discovery`) — business logic/normalization; call the repository
   interface. No `gin.Context`.
5. **Repository** — add the query to **both** implementations: Postgres
   (`storage/postgres`) and the in-memory fallback. Note any deliberate divergence
   ([ADR-0004](../adr/0004-in-memory-fallback-repository.md)).
6. **Register** the route in `RegisterRoutes`.
7. **Frontend sync** — if consumed by web, update `apps/web/lib/api.ts` types + calls in
   the same change.
8. **Docs** — update [`docs/services/api.md`](../services/api.md) and
   [`api-contracts.md`](../standards/api-contracts.md).

## Verification
Table-driven tests for the parser (valid + each invalid class + boundaries) and the
service (mocked repo). Run `make verify-api` (go vet + test + build). See
[Definition of Done → Backend](../verification/definition-of-done.md).

## Exit
Endpoint works on both repo paths, tests cover happy + error paths, envelope unchanged,
web types synced, docs updated.

## Common mistakes
- Business logic in the handler / SQL in the handler (forbidden edges).
- Updating only the Postgres path, not the fallback (or vice versa).
- Adding `meta`/`trace_id`/pagination speculatively (add only when justified —
  [`thinking/README.md`](../thinking/README.md) §3).
- Forgetting to sync `apps/web/lib/api.ts`.

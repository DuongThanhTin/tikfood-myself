---
name: add-api-endpoint
description: Use when adding or extending a discovery HTTP endpoint in apps/api — a new GET route, a new query filter, or a new field on a venue response. Keeps the handler->service->repo layering and the {data,error} envelope intact.
---

# Add / extend an API endpoint (apps/api)

Authoritative steps, context to load, and exit gate: **@docs/recipes/add-api-endpoint.md** — read it first.

## Quick checklist
1. Confirm in scope — discovery-only, not an anti-goal; read-only unless writes are approved.
2. Design the contract first — path, validated query params, response fields; reuse `{data,error}` (ADR-0003). Contract touched? -> use the `contract-first-change` skill.
3. handler (`internal/http`) -> service (`internal/discovery`) -> repository — no SQL/logic in the handler (ADR-0002).
4. Add the query to **both** repo impls (Postgres + in-memory fallback, ADR-0004); register the route.
5. Sync `apps/web/lib/api.ts` if the web consumes it; update `docs/services/api.md` + `docs/contracts/api.md`.
6. Table-driven tests (parser: valid + each invalid class + boundaries; service: mocked repo).

## Guard rails
- Writes / auth / migrations = **human approval** (CLAUDE.md; docs/ai/AI-CONTRACT.md §4).
- No breaking API changes; envelope shape unchanged.

## Done when
`make verify-api` (or `make verify`) is green and the surface meets docs/verification/definition-of-done.md.

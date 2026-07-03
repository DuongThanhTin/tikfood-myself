# Testing Standard

> **Why this doc exists:** testing guidance previously covered only the Go backend
> ([`backend/testing.md`](backend/testing.md)); `apps/web` has **zero tests** and there
> was no shared philosophy for *what to test, how much, and when*. This is the
> workspace-wide testing standard. It sets the philosophy and the frontend/automation
> expectations, and defers Go specifics to the backend doc.

**Related:** [`docs/standards/backend/testing.md`](backend/testing.md) (Go specifics) ·
[`docs/verification/definition-of-done.md`](../verification/definition-of-done.md)
(the completion gate, added in a later phase) · [`IMPROVEMENTS.md`](../../IMPROVEMENTS.md)
(current coverage gaps).

## Philosophy

- **Test behavior, not implementation.** Assert on observable outputs (responses,
  rendered results), not internal call order.
- **Prioritize by risk.** Cover logic that is easy to get wrong and expensive to break
  first: search/filter/normalization, price/geo math, envelope/error mapping,
  alias resolution. (These are exactly the current gaps in `IMPROVEMENTS.md`.)
- **Table-driven where inputs vary.** One case per meaningful input class + edge cases
  (empty, max length, out-of-range, boundary).
- **Deterministic.** No reliance on wall-clock, network, or ordering unless that is the
  thing under test. Prefer the in-memory fallback repo for API tests
  ([ADR-0004](../adr/0004-in-memory-fallback-repository.md)), but remember it diverges
  from Postgres on distance/dish — do not let a fallback-only test imply Postgres
  correctness.
- **Honesty.** Never claim tests passed unless they ran and passed; report skipped
  checks (per [`docs/ai/AI-CONTRACT.md`](../ai/AI-CONTRACT.md) §9).

## Coverage expectations

Coverage is risk-based, not a single percentage target. Every behavior change:

- **must** add or update at least one test that would fail without the change;
- **must** cover the error/invalid path for new validation or error codes;
- **should** add a boundary case for any new numeric/length constraint.

Pure logic (search normalization, mappers, parsers, price/geo rules) is expected to be
well covered; thin glue code is not required to be.

## Backend (Go)

Follow [`backend/testing.md`](backend/testing.md): `go test ./...`, table-driven tests,
mock `VenueRepository` for service tests, `httptest` for handler/router tests
(`internal/http/router_test.go` is the pattern). Priority gaps to close:
`normalizeSearch`/alias matching and the fallback repository filter.

## Frontend (Next.js) — target

`apps/web` has no test infrastructure yet. When adding it:

- Use **Vitest + React Testing Library** for component/unit tests; start with one test
  to bootstrap the harness, then cover `lib/api.ts` (envelope parsing, fallback
  behavior) and filter/format utilities as `DiscoveryExperience.tsx` is split.
- Test rendered behavior and accessibility (labels/`aria-label`), not styling details.
- Until the harness exists, `npm run web:typecheck` + `npm run web:build` are the
  minimum gate — state clearly that no unit tests ran.

> Setting up the frontend harness and CI touches build config / `.github/workflows`
> (a protected path). Track it in `IMPROVEMENTS.md`; it needs human approval.

## Automation runner

`packages/schemas/*` are the contract tests for runner I/O — any runner change must keep
responses valid against them. Prompt changes should include example inputs and expected
JSON shape (see the Prompt Standard: [`prompt-engineering.md`](prompt-engineering.md)).

## What to run before claiming done

| Area | Minimum commands |
| --- | --- |
| `apps/api` | `go test ./...` (and `go vet ./...`, `go build ./...`) |
| `apps/web` | `npm run web:typecheck`, `npm run web:build` (+ unit tests once present) |
| runner | `npm run runner:test`; validate I/O against `packages/schemas/` |

See [`docs/verification/definition-of-done.md`](../verification/definition-of-done.md)
for the full per-surface completion gate.

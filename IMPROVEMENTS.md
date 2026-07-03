# IMPROVEMENTS.md — TikFood Roadmap

Prioritized review findings. This is a roadmap, not a set of mandatory rules.
CLAUDE.md only references the relevant constraints (e.g. do not show fabricated
ratings).

## P0 — Product Credibility

### Fabricated UI data
`apps/web/components/DiscoveryExperience.tsx` hardcodes `mediaByVenue` with fake
ratings ("4.8"), review counts ("128 reviews"), badges, and expiring
`lh3.googleusercontent.com/aida-public/...` image URLs keyed to seed UUIDs. The
data model has no rating/review fields. `venues.photos[]` exists in the schema
but is not mapped in the Go `Venue` model or the repository query, and the UI
ignores it entirely — so surfacing real photos requires changes end to end:
model, SQL query, API response, and UI. Decision needed: either surface real
`photos` and add rating fields end to end, or remove the rating/review UI until
real data exists. Until then, fake data must be clearly marked as placeholder.

### Trend score & AI summary are static seed values
`trend_scores` and `ai_summaries` tables exist (migration 001) but no worker
populates them. The core differentiation — realtime trend scoring and AI "why
trending" summaries — does not exist yet. Roadmap: social ingestion →
trend-scoring worker → AI-summary worker.

## P1 — Engineering Robustness

### No CI
There is no `.github/workflows`. Add GitHub Actions: `go test` / `go vet` /
`go build` for `apps/api`, and `typecheck` / `lint` / `build` for `apps/web`.
This enforces the "tests updated" rule.

### Thin test coverage
Only `apps/api/internal/http/router_test.go` exists. Add table-driven tests for
`normalizeSearch` and alias matching (`apps/api/internal/discovery/search.go`)
and for the in-memory fallback repository filter. The long list SQL is a
regression risk.

### Location-alias logic duplicated in three places / two paradigms
`apps/api/internal/discovery/search.go` hardcodes `quan-1`/`quan-3`; the DB has
`locations` / `location_aliases` (migration 003) used by the SQL query; and
`apps/web/lib/api.ts` re-implements `normalizeLocationAlias`. Consolidate on the
DB as the source of truth; mark in-memory/fallback alias logic as dev-only.

### Fallback repository diverges from the Postgres path (verified at runtime)
Running the API with no `DATABASE_URL` (in-memory fallback) behaves differently
from the Postgres repository on two filters, confirmed by live calls:
- `distance_meters` is always `null` and `sort=distance` does not order by
  proximity — distance is only computed by the Postgres repo via PostGIS
  `ST_Distance`. (lat/lng validation still returns 400 when missing, correctly.)
- `dish=` matches `trending_dishes` by EXACT string in the fallback (`dish=pho`
  → 0 results, `dish=pho bo tai` → 1), whereas the Postgres path uses
  `like` / `normalized_name` substring matching.
Acceptable while the fallback is a dev-only approximation, but document it so the
divergence is not mistaken for a bug. Long term, either narrow the gap or gate
distance/dish features behind the Postgres path explicitly.

## P2 — API / Scaling Polish

### No pagination metadata
The list endpoint returns a bare array in `data` with only `limit`. This is
where a `meta` field (total / cursor) is legitimately useful — reserve it for
pagination. (Future.)

### OSRM public demo server
`apps/web/components/DiscoveryExperience.tsx` calls `router.project-osrm.org`
directly from the browser — rate-limited, no SLA, not production-safe. Replace
with self-hosted or managed routing.

### Shallow health check
`GET /health` returns `{ "ok": true }` without pinging the DB. Add a readiness
check that verifies the database before orchestration relies on it.

### Repo mixes automation tooling and product app
`apps/ai-code-runner` and the n8n workflow docs live beside the product app.
`starters/tikfood` hints the product graduates to its own repo. A conscious
decision for later.

## AI Operating System handoffs

The AI Operating System (`docs/ai/ROADMAP.md`) delivered documentation only. The
following are **code/infra** items it deliberately left to humans; each is governed by
an AI-OS doc so the eventual implementation has a standard to meet. Most already appear
above — this maps them to their governing doc/recipe.

- **CI** (`.github/workflows/`, protected path) — enables the completion gate in
  `docs/verification/definition-of-done.md`. See *No CI* (P1).
- **Frontend test harness** (Vitest + React Testing Library) — required by
  `docs/standards/testing.md` → *Frontend*; today `apps/web` has zero tests. Bootstrap
  one test, then cover `lib/api.ts` + filter/format utils. (Extends *Thin test
  coverage*, P1.)
- **Backend test gaps** — `normalizeSearch`/alias + fallback filter, per
  `docs/standards/testing.md`. See *Thin test coverage* (P1).
- **Location-alias consolidation** — DB as source of truth; recorded in
  `docs/adr/0004-in-memory-fallback-repository.md`. See *Location-alias logic
  duplicated* (P1).
- **Fabricated UI data** — governed by `docs/standards/application-security.md` §2 and
  `docs/services/web.md`. See *Fabricated UI data* (P0).
- **Trend-scoring / AI-summary workers** — architecture in `docs/architecture.md` §2.2;
  perf rules in `docs/standards/performance.md`. See *Trend score & AI summary* (P0).

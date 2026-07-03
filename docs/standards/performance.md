# Performance Standard

> **Why this doc exists:** performance guidance previously lived only inside
> [`backend/scaling.md`](backend/scaling.md) as a staged scaling philosophy — there was
> no objective bar for latency, query, or frontend cost that an AI agent could hold a
> change to. This doc sets pragmatic, MVP-appropriate targets and rules. It extends,
> and does not replace, the scaling doc.

**Related:** [`docs/standards/backend/scaling.md`](backend/scaling.md) (scaling stages) ·
[`docs/standards/backend/database.md`](backend/database.md) (schema/indexes) ·
[`docs/services/api.md`](../services/api.md).

## Principle: measure before optimizing

Do not add caches, indexes, or complexity speculatively. Optimize when a target below is
missed or a query is provably hot. Every optimization change should state the
before/after it targets. This mirrors the repo's "don't add X until justified" stance.

## Backend targets (discovery API, MVP)

Guideline budgets for a warm service at MVP data sizes (treat as review triggers, not
SLAs):

| Concern | Target |
| --- | --- |
| Read endpoint latency (p95) | < ~300 ms server-side |
| `GET /health` | < ~50 ms |
| Rows returned per list call | bounded by `limit` (≤ 100) |
| DB round-trips per request | 1 primary query; avoid N+1 |

Rules:

- **Bound every list query** (the `limit` cap of 100 exists for this reason). Add
  pagination `meta` only when an endpoint can realistically return large data — not
  before.
- **Index for the query, not the table.** Add indexes to support the actual filter/sort
  paths (city/district, geo via PostGIS, trend sort). Justify each index in the change
  and, if significant, an ADR. See [`backend/database.md`](backend/database.md).
- **No N+1.** Fetch related data in the primary query or a bounded batch.
- **Geo:** distance ordering uses PostGIS `ST_Distance` on the Postgres path; the
  in-memory fallback does not compute distance
  ([ADR-0004](../adr/0004-in-memory-fallback-repository.md)) — do not benchmark against
  the fallback.
- **Caching:** introduce only when a hot, stable query is demonstrated (scaling stage
  concern) — with an ADR describing invalidation.

## Frontend budgets (`apps/web`)

| Concern | Guideline |
| --- | --- |
| Map library | lazy-import MapLibre GL (already done) — keep it out of the initial bundle |
| Initial data | server component passes `initialVenues`; avoid a client refetch on first paint |
| Re-renders | memoize expensive list/map work; avoid rerendering the whole map on filter change |
| Images | use real, sized images once wired; avoid shipping large placeholder assets |

Rules:

- Keep heavy/optional dependencies lazy. Do not add a data/query/form library until a
  real caching/mutation/form need exists (per [`apps/web/CLAUDE.md`](../../apps/web/CLAUDE.md)).
- Do a bundle sanity check (`npm run web:build`) when adding a dependency; call out size
  impact.

## Workers / ingestion (roadmap)

Trend-scoring and AI-summary workers do not exist yet. When built: process in bounded
batches, make jobs idempotent, and keep per-item model cost visible (model-cost changes
require human approval, per [`AI-CONTRACT.md`](../ai/AI-CONTRACT.md)).

## Checklist for a performance-relevant change

- [ ] States the target/metric it addresses (before/after).
- [ ] List queries bounded; no N+1.
- [ ] New index/cache justified (and ADR'd if significant).
- [ ] Frontend: heavy deps lazy; bundle impact noted.
- [ ] Not benchmarking correctness-divergent paths (fallback vs Postgres).

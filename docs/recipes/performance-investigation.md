# Recipe — Performance investigation

**Goal:** find and fix a real performance problem, measured before and after.

**When to use:** a target in [`performance.md`](../standards/performance.md) is missed, or
a query/page is provably slow.

## Context to load
- [`docs/standards/performance.md`](../standards/performance.md) (targets/rules) ·
  [`docs/standards/backend/scaling.md`](../standards/backend/scaling.md) ·
  [`backend/database.md`](../standards/backend/database.md).
- [`docs/services/api.md`](../services/api.md) for the endpoint under scrutiny.

## Process skill
`systematic-debugging` applied to performance — measure, hypothesize, isolate, verify.

## Steps
1. **Measure first** — capture the current metric (latency/p95, query time, bundle size).
   No optimizing without a baseline.
2. **Locate the cost** — is it the query (missing index, N+1), the payload (unbounded
   `limit`), the app layer, or the frontend (bundle/re-render)? Measure on the **Postgres**
   path, not the fallback ([ADR-0004](../adr/0004-in-memory-fallback-repository.md)).
3. **Form one hypothesis** and test it in isolation.
4. **Apply the smallest fix** — bound the query, add a justified index, remove an N+1,
   lazy-load a heavy dep, memoize a hot render. Avoid speculative caching.
5. **Re-measure** — confirm the target is met; record before/after.
6. **ADR if significant** — a cache or a notable indexing/architecture change warrants
   an [ADR](../adr/) (esp. cache invalidation).

## Verification
Before/after numbers stated; tests still pass; no correctness regression on either repo
path. See [Definition of Done](../verification/definition-of-done.md).

## Exit
Target met with measured evidence, minimal change, no regressions, decision recorded if
significant.

## Common mistakes
- Optimizing without measuring; "it feels faster".
- Benchmarking the in-memory fallback and inferring production behavior.
- Adding a cache/index speculatively (complexity + invalidation risk).
- Fixing a micro-cost while ignoring the dominant one.

# ADR-0004: In-memory fallback repository selected by `DATABASE_URL`

- **Status:** Accepted
- **Date:** 2026-07-03 (retroactively recorded; decision predates this ADR)
- **Deciders:** Backend owners
- **Type:** retroactive (documents an existing decision)

## Context

Developers and CI need to run and exercise `apps/api` without provisioning a
PostgreSQL/PostGIS instance. At the same time production behavior depends on PostGIS
(distance) and SQL string matching that an in-memory store cannot fully replicate.

## Decision

`VenueRepository` is an interface with two implementations, chosen at startup by the
`app.Container` based on the `DATABASE_URL` env var:

- **`DiscoveryRepository`** (`storage/postgres`) — the real PostgreSQL/PostGIS path.
- **In-memory fallback** (`discovery`) — a seed-based approximation used when
  `DATABASE_URL` is unset.

The fallback is explicitly a **dev-only approximation**, not a second production path.

## Consequences

- **Positive:** the API runs with zero external dependencies for local dev and tests;
  the interface keeps the service layer storage-agnostic.
- **Negative / cost — known divergences (verified at runtime, see `IMPROVEMENTS.md`):**
  - `distance_meters` is always `null` in the fallback and `sort=distance` does not
    order by proximity (distance needs PostGIS `ST_Distance`).
  - `dish=` matches by **exact** string in the fallback vs. `like`/`normalized_name`
    substring matching in Postgres.
- **Follow-ups:** document the divergence so it is not mistaken for a bug (done in
  `IMPROVEMENTS.md`); long-term, either narrow the gap or gate distance/dish features
  behind the Postgres path explicitly. Prefer the DB as the source of truth for
  location-alias logic.

## Alternatives considered

- **Require Postgres for all runs** — rejected: heavier local/CI setup, slower feedback.
- **Testcontainers/ephemeral Postgres in tests** — not adopted yet; a reasonable future
  step to close the divergence (would warrant a superseding ADR).

## Sources / links

- [`apps/api/CLAUDE.md`](../../apps/api/CLAUDE.md) → *Repository*.
- [`IMPROVEMENTS.md`](../../IMPROVEMENTS.md) → *Fallback repository diverges from the Postgres path*.
- Related: [ADR-0001](0001-gin-and-pgx-not-gorm.md), [ADR-0002](0002-handler-service-repository-layering.md).

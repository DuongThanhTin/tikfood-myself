# ADR-0001: Gin + `database/sql` (pgx), not GORM

- **Status:** Accepted
- **Date:** 2026-07-03 (retroactively recorded; decision predates this ADR)
- **Deciders:** Backend owners
- **Type:** retroactive (documents an existing decision)

## Context

`apps/api` is a read-heavy discovery API over PostgreSQL/PostGIS. It needs precise,
geospatial SQL (e.g. `ST_Distance` for distance sorting) and predictable query
behavior. A choice was needed for the HTTP framework and the database access layer.

## Decision

We use **Gin** as the HTTP framework and **`database/sql` via the pgx stdlib driver**
for database access. We do **not** use an ORM (GORM).

Rules that follow: repositories write explicit SQL; PostGIS features are used directly;
no ORM abstractions are introduced.

## Consequences

- **Positive:** full control over SQL and PostGIS; no ORM impedance mismatch for
  geospatial queries; smaller dependency surface; behavior is easy to reason about.
- **Negative / cost:** more hand-written SQL and mapping code; no automatic
  migrations/associations from an ORM.
- **Follow-ups:** query correctness relies on tests (see thin-coverage note in
  `IMPROVEMENTS.md`); the in-memory fallback repository must approximate SQL behavior
  and diverges in places (see [ADR-0004](0004-in-memory-fallback-repository.md)).

## Alternatives considered

- **GORM / an ORM** — rejected: weaker fit for PostGIS geospatial SQL and adds an
  abstraction over the exact queries this API depends on.
- **A query builder (e.g. squirrel)** — not adopted; explicit SQL is preferred while
  the query set is small.

## Sources / links

- [`apps/api/CLAUDE.md`](../../apps/api/CLAUDE.md) → *Stack Reality* (states pgx, "NOT GORM").
- [`docs/standards/backend/framework.md`](../standards/backend/framework.md), [`docs/standards/backend/database.md`](../standards/backend/database.md).
- Related: [ADR-0002](0002-handler-service-repository-layering.md), [ADR-0004](0004-in-memory-fallback-repository.md).

# Backend Architecture Standard

> **This page is a pointer.** The detailed, authoritative backend standards live under
> [`docs/standards/backend/`](backend/). This file previously restated backend
> architecture and had drifted from the code (it named `net/http`, an `error`-string
> envelope, and in-memory-only persistence). To avoid a second, conflicting source it no
> longer describes the architecture — read the canonical docs below. Code wins on
> conflict.

Read these for backend work, in order:

- [`backend/README.md`](backend/README.md) — index + current backend state
- [`backend/structure.md`](backend/structure.md) — directory layout & ownership
- [`backend/architecture.md`](backend/architecture.md) — layering & responsibilities
- [`backend/patterns.md`](backend/patterns.md) — handler/service/repository patterns
- [`backend/request-response.md`](backend/request-response.md) — the real `{ data, error }` envelope
- [`backend/errors.md`](backend/errors.md) · [`backend/database.md`](backend/database.md) ·
  [`backend/logging.md`](backend/logging.md) · [`backend/testing.md`](backend/testing.md) ·
  [`backend/dependencies.md`](backend/dependencies.md) ·
  [`backend/search-filtering.md`](backend/search-filtering.md) ·
  [`backend/scaling.md`](backend/scaling.md) ·
  [`backend/framework.md`](backend/framework.md)

The reasoning behind key backend choices (Gin + pgx, layering, envelope) is recorded in
[`docs/adr/`](../adr/) (ADR-0001, 0002, 0003). App-specific rules "as built today" are in
[`apps/api/CLAUDE.md`](../../apps/api/CLAUDE.md).

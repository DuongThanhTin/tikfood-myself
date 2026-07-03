# ADR-0002: Handler → Service → Repository layering

- **Status:** Accepted
- **Date:** 2026-07-03 (retroactively recorded; decision predates this ADR)
- **Deciders:** Backend owners
- **Type:** retroactive (documents an existing decision)

## Context

`apps/api` needs a clear separation between HTTP concerns, business logic, and
persistence so that logic is testable without HTTP or a database, and so contributors
(human or AI) know exactly where a given change belongs.

## Decision

We enforce a one-way layering: **`http (Handler) → discovery (Service) → Repository →
postgres`**. Dependencies point one direction only.

Forbidden edges: Handler → Repository, Handler → SQL, Service → `gin.Context`,
Repository → business logic.

- **Handlers** (`internal/http`, named `*Handler`) parse/validate requests and format
  responses; no business logic, no DB.
- **Services** (`internal/discovery`) hold business logic and input normalization; no
  HTTP/`gin.Context` imports.
- **Repositories** do persistence only, behind an interface.

## Consequences

- **Positive:** business logic is unit-testable with a mocked `VenueRepository`; clear
  placement for new code; boundaries are checkable by inspection.
- **Negative / cost:** some boilerplate (DTOs, mapping across layers); requires
  discipline to not shortcut a layer.
- **Follow-ups:** validation lives in the handler, business rules in the service;
  transactions will live in the service layer once write endpoints exist.

## Alternatives considered

- **Fat handlers (logic in HTTP layer)** — rejected: untestable without HTTP, couples
  transport to domain.
- **Active-record / repository-in-domain** — rejected: mixes persistence with business
  rules, contrary to the repository-interface approach.

## Sources / links

- [`apps/api/CLAUDE.md`](../../apps/api/CLAUDE.md) → *Layering* (lists the forbidden edges).
- [`docs/standards/backend/architecture.md`](../standards/backend/architecture.md), [`docs/standards/backend/patterns.md`](../standards/backend/patterns.md).
- Related: [ADR-0003](0003-data-error-response-envelope.md), [ADR-0004](0004-in-memory-fallback-repository.md).

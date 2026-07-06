# Architecture Decision Records (ADRs)

> **Why this directory exists:** several significant decisions in this repo (why Gin
> and not GORM, why plain CSS, why the `{data,error}` envelope, why an in-memory
> fallback repository) lived only implicitly inside `CLAUDE.md` files and code
> comments. An AI agent had no durable "why" to consult, so it risked re-deriving or
> silently contradicting settled choices. ADRs are the **decision memory**: one short
> record per architecturally significant decision, capturing context, the decision,
> and its consequences.

**Related:** [`docs/architecture.md`](../architecture.md) (the resulting shape) ·
[`docs/ai/AI-CONTRACT.md`](../ai/AI-CONTRACT.md) · the handbook forecast of this dir in
[`docs/handbook/02-repository-design.md`](../handbook/02-repository-design.md) §2.5.

## What is an ADR

A lightweight (MADR-style) record of one decision that is **architecturally
significant** — it is expensive to reverse, constrains future work, or would surprise
a newcomer. Use the [`TEMPLATE.md`](TEMPLATE.md).

Record a decision when you:
- choose a framework, library, storage engine, or protocol;
- set a boundary or layering rule;
- fix a public contract shape (API envelope, error codes);
- deliberately defer or reject a common approach ("no X until Y").

Do **not** write an ADR for routine, easily-reversible implementation details.

## Conventions

- **File name:** `NNNN-short-slug.md` (zero-padded, e.g. `0007-use-redis-cache.md`).
- **Numbering:** monotonically increasing; never reuse a number.
- **Status:** `Proposed` → `Accepted` → optionally `Superseded by ADR-XXXX` or
  `Deprecated`. Never edit an accepted ADR's decision — write a new ADR that
  supersedes it and link both ways.
- **Honesty:** ADRs 0001–0006 are **retroactive** — they record decisions already
  live in the code/`CLAUDE.md`. That is stated in each so no false timeline is implied.
- **Source wins:** if an ADR ever disagrees with the code, the code wins — supersede
  the ADR (per `CLAUDE.md` → *Source of Truth*).

## When to open one (workflow)

1. During brainstorming/planning, if a decision meets the "significant" bar, draft an
   ADR as `Proposed`.
2. On approval, set `Accepted` and reference it from the PR.
3. If a later change reverses it, add a new ADR and mark the old one `Superseded`.

The core-workflow retrospective step (handbook Part 3) should ask: "did this task make
a decision worth an ADR?"

## Index

| ADR | Title | Status |
| --- | --- | --- |
| [0001](0001-gin-and-pgx-not-gorm.md) | Gin + `database/sql` (pgx), not GORM | Accepted (retroactive) |
| [0002](0002-handler-service-repository-layering.md) | Handler → Service → Repository layering | Accepted (retroactive) |
| [0003](0003-data-error-response-envelope.md) | `{ data, error }` response envelope | Accepted (retroactive) |
| [0004](0004-in-memory-fallback-repository.md) | In-memory fallback repository selected by `DATABASE_URL` | Accepted (retroactive) |
| [0005](0005-plain-css-frontend.md) | Plain CSS for `apps/web`, no framework yet | Accepted (retroactive) |
| [0006](0006-bilingual-handbook.md) | Bilingual handbook (VN prose + EN technical) | Superseded by [ADR-0008](0008-documentation-language-policy.md) |
| 0007 | *Reserved:* authentication approach — proposed in [`docs/features/authentication/`](../features/authentication/overview.md), ADR not yet written | Proposed |
| [0008](0008-documentation-language-policy.md) | Documentation language policy (one canonical language per doc) | Accepted |

# Recipe — Add a database migration (**human approval required**)

**Goal:** change the schema safely and reversibly, with approval.

**When to use:** new table/column/index, or altering `apps/api/migrations`.

> ⚠️ **Gate:** migrations are a human-approval area
> ([`docs/ai/AI-CONTRACT.md`](../ai/AI-CONTRACT.md) §4, [`docs/security.md`](../security.md)).
> Do not apply a migration autonomously. Propose it, get approval, then proceed.

## Context to load
- [`docs/standards/backend/database.md`](../standards/backend/database.md) (schema/entity/
  index rules) · existing files in `apps/api/migrations`.
- [`docs/standards/performance.md`](../standards/performance.md) (index for the query).
- ADR [0001](../adr/0001-gin-and-pgx-not-gorm.md) (explicit SQL, PostGIS).

## Process skill
`writing-plans` — migrations are a one-way-ish door; plan and get sign-off first.

## Steps
1. **Justify** the schema change against a real requirement; consider whether an
   [ADR](../adr/) is warranted (usually yes for new tables or contract-affecting columns).
2. **Request human approval** with the proposed DDL and rollback plan. **Wait.**
3. **Write the migration** following the existing file convention; make it additive and
   backward-compatible where possible (add nullable columns; avoid destructive drops in
   the same step).
4. **Indexes** — add only those the actual query paths need; justify each.
5. **Wire it through** — model, repository query, service, API response, and
   `apps/web/lib/api.ts` if the field is surfaced. Keep the `{data,error}` envelope.
6. **Seeds/fallback** — update seeds and the in-memory fallback so both repo paths stay
   coherent ([ADR-0004](../adr/0004-in-memory-fallback-repository.md)).
7. **Docs/ADR** — update `database.md` notes if conventions change; record the decision.

## Verification
Migration applies cleanly and is reversible; `make verify-api` passes; new field visible
end-to-end. See [Definition of Done](../verification/definition-of-done.md).

## Exit
Approved, applied cleanly, reversible, wired end-to-end (incl. fallback), ADR/docs
updated.

## Common mistakes
- Applying without approval (hard violation).
- Destructive change with no rollback.
- Schema change not reflected in the fallback repo or web types.
- Adding unused indexes / speculative columns.

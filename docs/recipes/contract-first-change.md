# Recipe — Contract-first change (shared shapes across services)

**Goal:** change a shape that crosses a service boundary — the `{data,error}` envelope,
a response field consumed by the web, `apps/web/lib/api.ts`, or `packages/schemas` —
without breaking a consumer.

**When to use:** any edit to a field/type that more than one side reads or writes. If the
change is purely internal to one app, use the app's recipe instead
([add-api-endpoint](add-api-endpoint.md) / [frontend-component](frontend-component.md)).

## Context to load
- [`docs/contracts/api.md`](../contracts/api.md) + [`docs/standards/backend/request-response.md`](../standards/backend/request-response.md)
- [ADR-0003](../adr/0003-data-error-response-envelope.md) (envelope) · [`docs/services/api.md`](../services/api.md) · [`docs/services/web.md`](../services/web.md)
- Schema changes: [`docs/contracts/runner.md`](../contracts/runner.md) — `packages/schemas/**` is a **protected path**.

## Process skill
`brainstorming` (enumerate consumers) -> `writing-plans` (if multi-step) -> `test-driven-development`.

## Steps
1. **Enumerate every consumer** of the shape *before* editing it (Go JSON tags,
   `apps/web/lib/api.ts`, `packages/schemas`, runtime prompts). List them explicitly
   ([thinking §5](../thinking/README.md#5-scope--risk-heuristics)).
2. **Prefer additive, backward-compatible changes** — add a field, don't rename/remove.
   A breaking change requires human approval and a migration plan for consumers.
3. **Change the producer and all consumers in the same change** — never leave a
   consumer reading a field that no longer exists.
4. **Schema (`packages/schemas`) touched? -> STOP for human approval** (protected path,
   [AI-CONTRACT §4](../ai/AI-CONTRACT.md)).
5. Update [`docs/contracts/api.md`](../contracts/api.md) and the relevant
   [`docs/services/`](../services/) doc.

## Verification
Tests on both sides of the boundary (producer + each consumer). Run `make verify` (all
apps), not just one. See [Definition of Done](../verification/definition-of-done.md).

## Exit
Every enumerated consumer updated and tested; envelope compatible; contracts/services
docs synced; no protected path changed without approval.

## Common mistakes
- Editing the producer but not `apps/web/lib/api.ts` (or vice versa).
- Renaming/removing a field where an additive change would do.
- Touching `packages/schemas` without human approval.

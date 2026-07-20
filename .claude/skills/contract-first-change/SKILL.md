---
name: contract-first-change
description: Use when a change touches a shape shared across services — the {data,error} envelope, a response field the web reads, apps/web/lib/api.ts, or packages/schemas. Forces enumerating every consumer before editing the contract.
---

# Contract-first change

Authoritative steps and exit gate: **@docs/recipes/contract-first-change.md** — read it first.

## Quick checklist
1. **Enumerate every consumer** of the shape before editing (Go JSON tags, `apps/web/lib/api.ts`, `packages/schemas`, runtime prompts).
2. Prefer **additive, backward-compatible** changes; a breaking change needs human approval + a consumer migration plan.
3. Change the **producer and all consumers in the same change**.
4. `packages/schemas` touched? -> **STOP for human approval** (protected path).
5. Update `docs/contracts/api.md` + the relevant `docs/services/` doc.

## Guard rails
- `packages/schemas/**` is a protected path — human approval required (@CLAUDE.md; @docs/ai/AI-CONTRACT.md §4).
- No breaking API changes; keep the `{data,error}` envelope compatible (ADR-0003).

## Done when
`make verify` (all apps, both sides of the boundary) is green and the change meets @docs/verification/definition-of-done.md.

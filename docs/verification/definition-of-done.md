# Definition of Done

> **Why this doc exists:** the `verification-before-completion` skill and the
> `Claude_Code_Work_Session_Playbook` both demand real verification, but there was no
> concrete, per-surface checklist of what "done" means. This is that gate. An agent must
> satisfy the relevant surface(s) below — and have the evidence — **before** claiming a
> task is complete, committing, or opening a PR.

**Related:** [`docs/standards/testing.md`](../standards/testing.md) ·
[`docs/verification/review-guide.md`](review-guide.md) ·
[`docs/ai/AI-CONTRACT.md`](../ai/AI-CONTRACT.md) §9 (honesty) ·
`Claude_Code_Work_Session_Playbook.md`.

## The rule

**Evidence before assertions.** Never say "done", "fixed", or "tests pass" without
having run the command and seen the result. If a check was skipped (e.g. no DB, no
frontend harness), say so explicitly. Claiming an unrun check is a hard violation
([`AI-CONTRACT.md`](../ai/AI-CONTRACT.md) §9).

## How to run the gates

The commands below are wired into the root [`Makefile`](../../Makefile) so there is one
canonical way to run each surface — use it instead of retyping loose commands:

| Surface | Command | Runs |
| --- | --- | --- |
| Backend | `make verify-api` | `go vet` + `go test` + `go build` (apps/api) |
| Frontend | `make verify-web` | typecheck + test + build (apps/web) |
| Runner | `make verify-runner` | `tsc --noEmit` + test (apps/ai-code-runner) |
| Everything | `make verify` | all three |

The per-surface checklists below describe **what** each gate proves; `make verify` is
**how** you run them.

## Universal (every task)

- [ ] Change is scoped to the requirement; no unrelated edits.
- [ ] Reused existing code where possible; no needless new helper/component.
- [ ] Not in a gated/anti-goal area without human approval.
- [ ] Behavior change is covered by a new/updated test (or skip is justified).
- [ ] Docs updated if architecture/behavior/contract changed (see below).
- [ ] Decision worth an [ADR](../adr/)? If yes, written.

## Backend — `apps/api`

- [ ] `go build ./...` succeeds.
- [ ] `go vet ./...` clean.
- [ ] `go test ./...` passes (report which packages ran).
- [ ] Response envelope unchanged / still `{data,error}` ([ADR-0003](../adr/0003-data-error-response-envelope.md)); no breaking field changes.
- [ ] New validation/error paths return correct codes and leak no SQL/internal detail.
- [ ] Layering intact — no forbidden edges ([ADR-0002](../adr/0002-handler-service-repository-layering.md)).
- [ ] If search/geo/dish touched: reasoned about **both** Postgres and fallback paths.
- [ ] Structured logging present for new paths; no secrets logged.

## Frontend — `apps/web`

- [ ] `make verify-web` passes (typecheck + Vitest unit tests + production build).
- [ ] New/changed component behavior is covered by a Vitest test (harness lives in
      `apps/web`; see `components/VenueList.test.tsx` for the pattern).
- [ ] Types in `lib/api.ts` still match the Go JSON tags (snake_case).
- [ ] API access still routed through `lib/api.ts`; envelope respected.
- [ ] Accessibility preserved (labels / `aria-label` on icon-only buttons).
- [ ] No fabricated data presented as real (placeholders clearly marked).

## Runtime prompt — `packages/prompts`

- [ ] Follows the [Prompt Standard](../standards/prompt-engineering.md) section order.
- [ ] JSON-only; success **and** failure shapes defined.
- [ ] Valid against the relevant `packages/schemas/` contract.
- [ ] Product context + anti-goals + security posture present.
- [ ] Example input + expected output shape provided.
- [ ] Paired `docs/agents/` role doc updated.

## Automation runner — `apps/ai-code-runner`

- [ ] `make verify-runner` passes (`tsc --noEmit` + Vitest; see `src/guards/commandGuard.test.ts`).
- [ ] Responses validate against `packages/schemas/*`.
- [ ] Hard contract honored: input validated first; `ai/*` branch only; no `main`/`master`
      push; no secrets read; no success claim without verified+committed+pushed branch.
- [ ] Unfinished stages clearly marked `TODO`; no production claims.

## Docs to update when relevant

| If you changed… | Update |
| --- | --- |
| Directory layout | [`docs/REPOSITORY-MAP.md`](../REPOSITORY-MAP.md) |
| API surface / envelope | [`docs/services/api.md`](../services/api.md), `apps/web/lib/api.ts`, [`api-contracts.md`](../standards/api-contracts.md) |
| A significant decision | a new [ADR](../adr/) |
| Runner contract/schema | [`docs/runner-contract.md`](../runner-contract.md) (schema = protected path) |
| A rule/standard | the owning standard + its index |

## Before opening a PR

- [ ] Branch is under `ai/`.
- [ ] Diff reviewed against [`review-guide.md`](review-guide.md).
- [ ] PR body states what was verified and what was skipped — honestly.
- [ ] No merge, no push to `main`/`master`, no force-push.

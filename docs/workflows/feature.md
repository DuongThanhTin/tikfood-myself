# Workflow — Feature (request → PR)

**Goal:** take a feature from idea to a clean PR that a later session can continue.

**When to use:** "Add / build / support X" — a behavior/UX change (not a bug or refactor).

> Touches auth/migration/infra/external/anti-goal? → [gated-change.md](gated-change.md) **first**, then return here.

## Context to load
[CONTEXT-LOADING §1 core + §2 matrix](../ai/CONTEXT-LOADING.md) + the owning `apps/*/CLAUDE.md`.

## Process skill
`brainstorming` → (large) `writing-plans` → `test-driven-development` → `verification-before-completion`.

## Scope — which steps you may skip
- **Small** (skip steps 5–6; brief = a PR comment): one app only; no contract change
  ({data,error} / `lib/api.ts` / `packages/schemas`); no new domain; reversible; not
  gated; fits one recipe/playbook.
- **Large** (steps 5–6 **required**, written to files): touches both `apps/api`+`apps/web`;
  changes a contract/schema; adds a domain; multi-step/milestone; or any gated area.
- **Required at every size:** steps 1, 2, 4, 7, 8, 9, 10, 12.

## Steps
> Input of each step = the previous step's artifact, unless noted.

| # | Step | Goal | Action (doc) | Artifact | Exit | 🚦 Gate |
|---|------|------|--------------|----------|------|---------|
| 1 | Classify *(req)* | Name task type; spot gated areas | [thinking §1](../thinking/README.md#1-classify-the-task-first) | Task-type label | Classified | Gated → [gated-change](gated-change.md) |
| 2 | Load context *(req)* | Read the right, minimal docs | [CONTEXT-LOADING §1+§2](../ai/CONTEXT-LOADING.md) | Docs read | Enough to act | — |
| 3 | Brief | Agree scope/acceptance | [new-feature §2](../getting-started/new-feature.md) + [feature-brief](../getting-started/templates/feature-brief.md) | Brief (small: comment / large: → spec) | 1-line goal + acceptance | — |
| 4 | Brainstorm *(req)* | Clarify intent before code | skill `brainstorming` | Agreed approach | No ambiguity | **Stop for human if scope/anti-goal unclear** |
| 5 | Spec *(large)* | Design | → `docs/superpowers/specs/<date>-<feature>-design.md` | Spec | Design agreed | **Stop for spec approval** |
| 6 | Plan *(large)* | Break into milestones/PRs | skill `writing-plans` → `docs/superpowers/plans/` | Plan | Plan agreed | **Stop for plan approval** |
| 7 | Branch *(req)* | Clean workspace | [new-feature §1](../getting-started/new-feature.md) — `git checkout -b ai/<feature>` | `ai/*` branch | On new branch | Never branch off `main` |
| 8 | Implement *(req)* | Smallest correct change, TDD | UI → [design/playbooks](../design/playbooks/README.md); backend/domain → [recipes](../recipes/) | Code + tests | Recipe/playbook done | New gated area → [gated-change](gated-change.md) |
| 9 | Test *(req)* | Cover new behavior | [testing.md](../standards/testing.md) | Green tests | Suite passes | — |
| 10 | Verify *(req)* | Evidence of "done" | skill `verification-before-completion` + [DoD](../verification/definition-of-done.md) | Command logs | DoD satisfied | — |
| 11 | Docs | Keep docs in sync | [DoD → Docs to update](../verification/definition-of-done.md#docs-to-update-when-relevant) | Updated docs/ADR | No drift | Schema/protected → [gated-change](gated-change.md) |
| 12 | PR *(req)* | Hand off for review | [new-feature §6](../getting-started/new-feature.md) + [pr template](../getting-started/templates/pr-description.md) | PR | Body states verified/skipped | **No merge / no push to main / no force-push** |

## Verification
Follow the [Definition of Done](../verification/definition-of-done.md) for the surface(s) touched (backend/frontend); every command actually run, with logs.

## Exit
Acceptance met, tests green, contracts intact, docs/ADR updated, PR open under `ai/` — not merged.

## Common mistakes
- Jumping into code, skipping brainstorm (violates process-first, [handbook/01](../handbook/01-ai-operating-system.md)).
- A large feature with no spec/plan on disk → the next session loses context.
- Changing a contract without enumerating consumers ([thinking §5](../thinking/README.md#5-scope--risk-heuristics)).
- Claiming done before running verification ([AI-CONTRACT §9](../ai/AI-CONTRACT.md)).
- Doing gated work without approval ([gated-change](gated-change.md)).

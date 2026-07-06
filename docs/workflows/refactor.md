# Workflow — Refactor (request → PR)

**Goal:** improve structure/readability/size while keeping behavior identical.

**When to use:** a file/function carries multiple responsibilities, duplication accumulated, or an approved split is due (e.g. `DiscoveryExperience.tsx`).

> If the "refactor" actually changes a contract/schema/behavior → it is not a refactor → [feature](feature.md) / [gated-change](gated-change.md).
> A measured **performance** problem (not just structure) → [recipe performance-investigation](../recipes/performance-investigation.md).

## Context to load
The owning `apps/*/CLAUDE.md` + its "Future Direction" section; [thinking §3](../thinking/README.md#3-reading-when-justified-standards) (is the refactor justified yet?).

## Process skill
`writing-plans` for anything non-trivial; small steps, each independently verifiable.

## Steps
> Input of each step = the previous step's artifact, unless noted.

| # | Step | Goal | Action (doc) | Artifact | Exit | 🚦 Gate |
|---|------|------|--------------|----------|------|---------|
| 1 | Classify + justify | Confirm refactor & that it's needed *now* | [thinking §1](../thinking/README.md#1-classify-the-task-first) + §3 | Justification | Justified | Contract/schema change → [feature](feature.md) / [gated-change](gated-change.md) |
| 2 | Branch | Clean workspace | [new-feature §1](../getting-started/new-feature.md) — `git checkout -b ai/refactor-<area>` | `ai/*` branch | On new branch | Never branch off `main` |
| 3 | Safety net | Pin current behavior | ensure tests cover it; add characterization tests first if thin ([recipe refactor](../recipes/refactor.md)) | Green baseline tests | Suite green before editing | — |
| 4 | Refactor in steps | Move one responsibility at a time, stay green | [recipe refactor](../recipes/refactor.md); keep contracts ({data,error}, `lib/api.ts` types, exported signatures) | "Movement, not behavior change" diff | Each step green | Found a bug → split [bug](bug.md) |
| 5 | Verify no behavior change | Prove identical behavior | full suite green *identically* before/after + `vet`/`build` or `web:typecheck`/`web:build` ([DoD](../verification/definition-of-done.md)) | Logs | DoD satisfied | — |
| 6 | Docs + PR | Hand off | [DoD → Docs](../verification/definition-of-done.md#docs-to-update-when-relevant) + [pr template](../getting-started/templates/pr-description.md) | PR | PR open | **No merge / no push to main / no force-push** |

## Verification
Suite passes identically before & after; the diff is movement only, no behavior change. See [Definition of Done](../verification/definition-of-done.md).

## Exit
Behavior identical, tests green throughout, contracts intact, structure genuinely improved.

## Common mistakes
- Mixing a refactor with a feature/bug (hard to review, hides regressions).
- Refactoring without a test safety net.
- Introducing an abstraction before it's justified ([thinking §3](../thinking/README.md#3-reading-when-justified-standards)).
- Breaking a public contract "while in there".

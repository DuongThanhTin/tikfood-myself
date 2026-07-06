# Workflow — Bug (request → PR)

**Goal:** find the real root cause and fix it with a regression test — not a symptom patch.

**When to use:** something is broken / fails / behaves incorrectly.

> If the fix must touch auth/migration/infra/external → [gated-change.md](gated-change.md) first.

## Context to load
The owning `apps/*/CLAUDE.md` + [CONTEXT-LOADING §2](../ai/CONTEXT-LOADING.md); search/geo/dish bugs: [ADR-0004](../adr/0004-in-memory-fallback-repository.md) (may be a known fallback-vs-Postgres divergence, not a defect).

## Process skill
`systematic-debugging` — **do not propose a fix before the root cause is found**.

## Steps
> Input of each step = the previous step's artifact, unless noted.

| # | Step | Goal | Action (doc) | Artifact | Exit | 🚦 Gate |
|---|------|------|--------------|----------|------|---------|
| 1 | Classify + gate | Confirm it's a bug & not gated | [thinking §1](../thinking/README.md#1-classify-the-task-first) | Label | Classified | Fix touches gated → [gated-change](gated-change.md) |
| 2 | Load context | Narrow to layer/input | owner doc + [CONTEXT-LOADING](../ai/CONTEXT-LOADING.md) | Suspect area | Enough to reproduce | — |
| 3 | Branch | Clean workspace | [new-feature §1](../getting-started/new-feature.md) — `git checkout -b ai/fix-<bug>` | `ai/*` branch | On new branch | Never branch off `main` |
| 4 | Investigate + fix | Root cause → smallest fix | skill `systematic-debugging` via [recipe fix-bug](../recipes/fix-bug.md) (reproduce → localize → root cause → failing test → fix → siblings) | Code + regression test | fix-bug recipe done | Bigger work → split [feature](feature.md) / [refactor](refactor.md) |
| 5 | Verify | Evidence | skill `verification-before-completion` + [DoD](../verification/definition-of-done.md); regression test red-before / green-after | Logs | DoD satisfied | — |
| 6 | Docs + PR | Hand off | [DoD → Docs](../verification/definition-of-done.md#docs-to-update-when-relevant) + [pr template](../getting-started/templates/pr-description.md) | PR stating root cause | PR open | **No merge / no push to main / no force-push** |

## Verification
`go test ./...` (or `web:typecheck`+`web:build`); the regression test fails before / passes after. See [Definition of Done](../verification/definition-of-done.md).

## Exit
Root cause stated, regression test added, fix minimal, no regressions, siblings checked.

## Common mistakes
- Patching the symptom without understanding the cause.
- "Fixing" a fallback divergence that is actually expected ([ADR-0004](../adr/0004-in-memory-fallback-repository.md)).
- No regression test (the bug silently returns).
- Refactoring while fixing (split it out → [refactor](refactor.md)).

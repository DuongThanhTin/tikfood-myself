# Workflow — Gated change (auth · migration · infra · external calls)

**Goal:** pass through an approval-gated area safely: prepare artifacts → request approval → wait → execute.

**When to use:** the task touches auth/authorization, a DB migration, infra/deploy, an external network call, a protected path, or any anti-goal.

> ⚠️ **Gates are hard rules:** [AI-CONTRACT §4](../ai/AI-CONTRACT.md) · [security.md](../security.md) · [thinking §6](../thinking/README.md#6-when-to-stop-and-ask). Never execute autonomously.

## Context to load
[AI-CONTRACT §1 + §4](../ai/AI-CONTRACT.md), [tikfood/anti-goals](../tikfood/anti-goals.md), the owning doc; migration → [recipe add-migration](../recipes/add-migration.md).

## Process skill
`writing-plans` — a one-way door; plan and get sign-off first.

## Steps
> Input of each step = the previous step's artifact, unless noted.

| # | Step | Goal | Action (doc) | Artifact | Exit | 🚦 Gate |
|---|------|------|--------------|----------|------|---------|
| 1 | Identify the gate | Name which area(s) it touches | [thinking §1](../thinking/README.md#1-classify-the-task-first) + [AI-CONTRACT §4](../ai/AI-CONTRACT.md) | List of gates hit | Listed | — |
| 2 | Prepare approval artifacts | Give the human enough to decide | draft: justification, design/[ADR](../adr/), (migration) DDL + rollback, (external) endpoint + data sent, (infra) change + blast radius | Proposal pack | Enough to approve | — |
| 3 | Request approval & STOP | Get sign-off | present the pack, **wait for the human** | Recorded decision | Explicit approval | **STOP — no code until approved** |
| 4 | Branch + execute via recipe | Stay on the right rail | `git checkout -b ai/<gated>`; migration → [add-migration](../recipes/add-migration.md); else the matching recipe/workflow; wire the fallback too ([ADR-0004](../adr/0004-in-memory-fallback-repository.md)) | Code + tests | Recipe done | New gate appears → back to step 1 |
| 5 | Verify + ADR/docs | Evidence + record the decision | [DoD](../verification/definition-of-done.md) + [DoD → Docs](../verification/definition-of-done.md#docs-to-update-when-relevant); write the [ADR](../adr/) | Logs + ADR/docs | DoD satisfied | — |
| 6 | PR | Hand off for review | [pr template](../getting-started/templates/pr-description.md); body states approval + what was verified | PR | PR open | **No merge / no push to main / no force-push; no auto-applied migration** |

## Verification
Approval obtained before editing; migration applies cleanly & is reversible; `go test ./...` passes; the ADR is written. See [Definition of Done](../verification/definition-of-done.md).

## Exit
Approved first, executed via the right recipe, reversible, fallback wired, ADR/docs updated.

## Common mistakes
- Writing code before approval (a hard violation of [AI-CONTRACT §4](../ai/AI-CONTRACT.md)).
- A destructive migration with no rollback.
- Changing a schema but forgetting the fallback repo / web types.
- Skipping an external network call — it is also "gated".
- Not writing an ADR for a binding decision.

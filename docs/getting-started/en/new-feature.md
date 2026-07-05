# Starting a new feature

> 🇻🇳 Tiếng Việt: [`../new-feature.md`](../new-feature.md)
>
> Use when: you have an idea/feature and want to go from scratch to PR cleanly, so any
> new session can understand and continue.

## The 6-step flow

### 1. Create a clean branch (don't work on an in-progress branch)
```bash
git checkout dev && git pull            # base = dev (integration). Target main for the PR.
git checkout -b ai/<feature-name>        # PR branches ALWAYS under the ai/ prefix
```
> Rules from `CLAUDE.md`: PR branches under `ai/`; **never** push `main`/`master`, no
> force-push, no auto-merge. To work in parallel with other work → use
> [worktrees](worktrees.md).

### 2. Write a feature brief (required — so future sessions understand)
Use [`templates/feature-brief.md`](templates/feature-brief.md). Keep it short: **one-line
goal · scope · out-of-scope · acceptance · constraints**. Store it:
- small task → at the top of the PR or a comment;
- large/multi-step → a **spec** at `docs/superpowers/specs/<date>-<feature>-design.md`
  (the repo's spec-driven convention).

### 3. Kick off Claude
Paste [`templates/session-kickoff.md`](templates/session-kickoff.md) (the "new feature"
variant) and **require a brainstorm/plan before coding**. For multi-step work, have Claude
write the plan to a file first.

### 4. Pick the right "track"
- **UI** (`apps/web`) → follow [`docs/design/`](../../design/00-overview.md) + the matching
  playbook in [`docs/design/playbooks/`](../../design/playbooks/README.md) (`create-page`,
  `create-component`, `create-form`, `create-modal`…).
- **Backend/domain work** → [`docs/recipes/`](../../recipes).
- Needs **human approval**: auth, migration, infra, external network calls, or anything
  touching an anti-goal (`CLAUDE.md`).

### 5. Build — reuse, tokens, tests
Reuse first; use tokens not hardcoded values; data through `lib/api.ts`; Vietnamese copy;
update tests when behavior changes.

### 6. Finish & PR
Check against [`docs/verification/definition-of-done.md`](../../verification/definition-of-done.md).
Commit (message ends with the repo's `Co-Authored-By:` line). Open the PR using
[`templates/pr-description.md`](templates/pr-description.md) (body ends with the
`🤖 Generated with Claude Code` line).

## Worked example — "clear-filters button" in the venue rail empty-state (UI)

```bash
git checkout dev && git pull
git checkout -b ai/rail-empty-state-reset
```

**Brief (paste for Claude):**
> Feature: when the rail has no venues, show an empty-state with a "Clear filters" button.
> Scope: `apps/web` rail only. Out-of-scope: API changes, pagination.
> Acceptance: empty → message + button; click → calls `clearFilters` and re-searches.
> Constraints: discovery-only; follow docs/design; use tokens; no new component if reusable.

**Track:** read [`docs/design/patterns/empty-state.md`](../../design/patterns/empty-state.md)
+ [`docs/design/playbooks/update-component.md`](../../design/playbooks/update-component.md)
(this extends an existing pattern, not a new component).

**Finish:** check against DoD → commit → PR from the template.

## Worked example — big feature (backend + web) → split tabs

If the feature touches both `apps/api` and `apps/web`, consider **two worktrees/branches**
so two tabs work in parallel without file collisions (see [`worktrees.md`](worktrees.md)),
then land it via two PRs.

## Related

- [`new-terminal.md`](new-terminal.md) · [`worktrees.md`](worktrees.md) · [`templates/`](templates/)
- [`docs/design/playbooks/`](../../design/playbooks/README.md) · [`docs/recipes/`](../../recipes) · [`docs/verification/definition-of-done.md`](../../verification/definition-of-done.md)

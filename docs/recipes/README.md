# Recipes (Domain Playbooks)

> **Why this directory exists:** the repo's principles (load the right context, process
> skill first, smallest correct change, verify honestly) are general — but doing a
> *specific* task well means applying them in a specific order with the right docs. These
> recipes are reusable, step-by-step playbooks for the common task types. They embody the
> roadmap principle **"reusable docs over large prompts"**: instead of one giant prompt
> per task, a recipe names the exact skill, docs, steps, and exit gate.

**Related:** [`docs/thinking/README.md`](../thinking/README.md) (how to classify/reason) ·
[`docs/ai/CONTEXT-LOADING.md`](../ai/CONTEXT-LOADING.md) (what to read) ·
[`docs/verification/definition-of-done.md`](../verification/definition-of-done.md) (the gate) ·
[`docs/workflows/`](../workflows/) (end-to-end request→PR flows that call these recipes).
These fill the handbook Part 4 "Domain Workflows" target; the handbook links here rather
than duplicating.

## How to use

1. Classify the task ([`thinking/README.md`](../thinking/README.md) §1).
2. Open the matching recipe and follow it top to bottom.
3. Stop at any **human-approval gate** the recipe marks.
4. Finish against the recipe's exit criteria + the [Definition of Done](../verification/definition-of-done.md).

## Recipe format

Every recipe has: **Goal · When to use · Context to load · Process skill · Steps ·
Verification · Exit · Common mistakes**.

## Index

| Recipe | Use when |
| --- | --- |
| [add-api-endpoint](add-api-endpoint.md) | Adding/extending an `apps/api` endpoint |
| [frontend-component](frontend-component.md) | Adding/changing `apps/web` UI |
| [fix-bug](fix-bug.md) | Something is broken or behaves wrong |
| [refactor](refactor.md) | Improving structure without changing behavior |
| [add-runtime-prompt](add-runtime-prompt.md) | Adding/editing a `packages/prompts` prompt |
| [add-migration](add-migration.md) | Schema/migration change (**human approval**) |
| [performance-investigation](performance-investigation.md) | A latency/cost target is missed |
| [security-review](security-review.md) | Reviewing a change for security |

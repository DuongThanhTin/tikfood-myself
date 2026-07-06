# Workflows — runnable end-to-end flows for AI agents

> **Why this directory exists:** [recipes](../recipes/) are *per-domain* playbooks
> (add-api, fix-bug…); [handbook/03](../handbook/03-core-workflow.md) is the 15-step
> *map*; [getting-started/new-feature](../getting-started/new-feature.md) is a human
> quickstart. This directory stitches them into **runnable request → PR flows per task
> type**. Principle: **orchestrate by link, do not copy content** (like [AI-CONTRACT](../ai/AI-CONTRACT.md)).

**Related:** [thinking §1](../thinking/README.md#1-classify-the-task-first) (classify) ·
[CONTEXT-LOADING](../ai/CONTEXT-LOADING.md) (what to read) · [recipes](../recipes/) (what to do per domain) ·
[definition-of-done](../verification/definition-of-done.md) (the "done" gate).

## Step fields

Every workflow lists its steps as a table with these columns: **Goal · Action (doc) ·
Artifact · Exit · Gate**. The **input** of a step is the **artifact of the previous
step**, unless the step says otherwise. **Gate** = when to stop and wait for human approval.

## Choose a workflow by task type (mapped from [thinking §1](../thinking/README.md#1-classify-the-task-first))

| Signal | Task type | Workflow |
| --- | --- | --- |
| "Add / build / support X" | Feature | [feature.md](feature.md) |
| "Broken / wrong / fails" | Bug | [bug.md](bug.md) |
| "Cleaner / more readable" (no behavior change) | Refactor | [refactor.md](refactor.md) |
| "Slow / over latency-cost budget" | Performance | [recipe: performance-investigation](../recipes/performance-investigation.md) |
| "Is it safe / is X exposed" | Security review | [recipe: security-review](../recipes/security-review.md) |
| "How does X work / where is Y" | Research | read-only, delegate if broad — [CONTEXT-LOADING §4](../ai/CONTEXT-LOADING.md#4-when-to-drop-context-or-delegate) |
| Touches auth / migration / infra / external / anti-goal | **Gated** | [gated-change.md](gated-change.md) **first** |

## How to use

1. Classify the task ([thinking §1](../thinking/README.md#1-classify-the-task-first)).
2. Open the matching workflow and run it top to bottom.
3. Stop at every 🚦 gate.
4. Finish against the [Definition of Done](../verification/definition-of-done.md).

> **Gated always wins:** if a task touches a gated area, run [gated-change.md](gated-change.md)
> first, then return to the original workflow to implement.

# Recipe — <<TASK_TYPE>>

> **Why this file exists:** a per-domain playbook for a task that recurs and that agents
> tend to do inconsistently. Write one **only when** a class of work has happened ≥2 times
> and gone wrong/divergent — not preemptively. Recipes are the "how"; [workflows](../workflows/README.md)
> are the "when/order".

## When to use
<!-- FILL: dấu hiệu nhận ra loại việc này. -->
<<Signal: e.g. "add a new API endpoint", "add a DB migration".>>

## Context to load
<!-- FILL: đọc gì trước — app CLAUDE.md, contract, standard liên quan. -->
[CONTEXT-LOADING](../ai/CONTEXT-LOADING.md) + <<relevant app CLAUDE.md / contract>>.

## Steps
<!-- FILL: các bước cụ thể của repo này. Ngắn, theo thứ tự, nêu file/layer chạm vào. -->
1. <<step — which file/layer>>
2. <<step>>
3. **Contract/tests** — keep the contract stable ([contracts](../contracts/README.md)); add/adjust tests.

## Verification
Run `make verify-<<app>>`. See [Definition of Done](../verification/definition-of-done.md).

## Exit
<<What "done" looks like for this task type — reused code, contract intact, tests cover happy + error.>>

## Common mistakes
- <<recurring mistake this recipe exists to prevent>>

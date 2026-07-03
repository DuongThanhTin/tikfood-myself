# Recipe — Add / edit a runtime prompt (`packages/prompts`)

**Goal:** add or change a Layer-3 runtime prompt that stays consistent, safe, and
JSON-valid.

**When to use:** creating a new pipeline prompt (e.g. `feature-analyzer`,
`pull-request-writer`) or editing an existing one.

## Context to load
- [`docs/standards/prompt-engineering.md`](../standards/prompt-engineering.md) (the Prompt
  Standard — the template) · existing prompts in [`packages/prompts/`](../../packages/prompts/).
- The paired role doc in [`docs/agents/`](../agents/).
- [`docs/runner-contract.md`](../runner-contract.md) + [`packages/schemas/`](../../packages/schemas/).

## Process skill
`brainstorming` to pin the role's single responsibility and its I/O; then author.

## Steps
1. **One role, one prompt** — define the single job in the pipeline
   (analyzer → context-reader → coding-agent → reviewer → PR-writer).
2. **Use the template** — sections in the standard's order: Product Context → Required
   Inputs → Required Behavior → Verification → Output Success (JSON) → Output Failure
   (JSON) → Rules.
3. **Bind to a schema** — reference the matching `packages/schemas/*` where one exists;
   define both success and failure JSON.
4. **Embed guardrails** — product context + anti-goals, forbidden paths (`.env*`,
   secrets), `ai/*`-only branches, no `main`/`master` push, prompt-injection posture.
5. **JSON-only** — state "valid JSON only, no Markdown/prose/trailing commas".
6. **Provide an example** — sample input + expected output shape.
7. **Sync the role doc** in [`docs/agents/`](../agents/); if the pipeline order changes,
   update [`docs/architecture.md`](../architecture.md) §1 and
   [`docs/services/ai-code-runner.md`](../services/ai-code-runner.md).

## Verification
Output validates against the schema; example round-trips; role doc and prompt agree. See
[Definition of Done → Runtime prompt](../verification/definition-of-done.md).

## Exit
Prompt follows the standard, JSON-valid against its schema, guardrails present, example
included, paired docs synced.

## Common mistakes
- Merging two roles into one prompt.
- Free-form/Markdown output; missing failure shape.
- Dropping product context / anti-goals / security rules.
- Editing the prompt but not its `docs/agents/` role doc.
- Changing a schema without human approval (protected path).

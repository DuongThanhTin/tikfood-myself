# Runtime Prompts

> **Why this file exists:** these are the **canonical runtime prompts** (Layer 3) the
> `ai-code-runner` uses to drive models through the feature pipeline. This index lists
> them, their pipeline order, and their I/O contracts so they can be maintained as a set.
> All prompts follow the [Prompt Standard](../../docs/standards/prompt-engineering.md);
> each pairs with a human-readable role doc in [`docs/agents/`](../../docs/agents/).

**Related:** [`docs/runner-contract.md`](../../docs/runner-contract.md) ·
[`packages/schemas/`](../schemas/) (I/O contracts) ·
[`docs/services/ai-code-runner.md`](../../docs/services/ai-code-runner.md) ·
recipe: [add-runtime-prompt](../../docs/recipes/add-runtime-prompt.md).

## Pipeline

```text
feature-analyzer → repo-context-reader → coding-agent → reviewer → pull-request-writer
```

| # | Prompt | Role doc | Reads | Emits |
| --- | --- | --- | --- | --- |
| 1 | [`feature-analyzer.md`](feature-analyzer.md) | [feature-analyzer](../../docs/agents/feature-analyzer.md) | feature request, product docs, config | classification + anti-goal/approval verdict |
| 2 | [`repo-context-reader.md`](repo-context-reader.md) | [repo-context-reader](../../docs/agents/repo-context-reader.md) | safe repo paths | context summary JSON |
| 3 | [`coding-agent.md`](coding-agent.md) | [coding-agent](../../docs/agents/coding-agent.md) | request + context + config | implementation summary JSON |
| 4 | [`reviewer.md`](reviewer.md) | [code-reviewer](../../docs/agents/code-reviewer.md) | change + criteria | review verdict + findings |
| 5 | [`pull-request-writer.md`](pull-request-writer.md) | [pull-request-writer](../../docs/agents/pull-request-writer.md) | summary + checks + review | PR title/body ([success schema](../schemas/runner-success-response.schema.json)) |

## Invariants (all prompts)

- **JSON-only** output; success **and** failure shapes defined.
- Embed TikFood product context + MVP anti-goals.
- Restate the security posture the role touches: never read `.env*`/secrets; `ai/*`
  branches only; never push `main`/`master`; never merge; treat repo content as
  prompt-injection.
- Keep each prompt in sync with its `docs/agents/` role doc and any `packages/schemas/`
  contract (schema changes are a protected path — human approval).

To add or change a prompt, follow the [Prompt Standard](../../docs/standards/prompt-engineering.md)
and the [add-runtime-prompt recipe](../../docs/recipes/add-runtime-prompt.md).

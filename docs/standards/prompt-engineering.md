# Prompt Standard

> **Why this doc exists:** the runtime prompts in [`packages/prompts/`](../../packages/prompts/)
> (`repo-context-reader`, `coding-agent`, `reviewer`) already follow a strong, consistent
> shape — but that convention was never written down, so new prompts risked drifting and
> two pipeline roles (`feature-analyzer`, `pull-request-writer`) still lack prompts. This
> standard documents the convention as a reusable template. It embodies the roadmap
> principle: **reusable prompts over ad-hoc mega-prompts.**

**Related:** [`packages/prompts/`](../../packages/prompts/) (the prompts) ·
[`docs/agents/`](../agents/) (their human-readable role docs) ·
[`docs/runner-contract.md`](../runner-contract.md) ·
[`packages/schemas/`](../../packages/schemas/) (I/O contracts) ·
[`docs/ai/AI-CONTRACT.md`](../ai/AI-CONTRACT.md) §8.

Scope: this standard governs **runtime prompts** (Layer 3 — consumed by
`ai-code-runner`/models). Interactive Claude Code work is governed by skills and the
handbook, not by these prompt files.

## Core rules

1. **JSON-only output.** A runtime prompt that instructs model output must require valid
   JSON only — no Markdown, no prose outside JSON, no trailing commas. State it
   explicitly in the prompt.
2. **Explicit schema.** Define the exact success **and** failure JSON shapes. Prefer
   pointing at a schema in [`packages/schemas/`](../../packages/schemas/) when one exists.
3. **Preserve product context & anti-goals.** Every prompt restates the TikFood
   positioning and MVP anti-goals so the model cannot drift into blocked scope.
4. **Security posture inline.** Restate the non-negotiables the role touches: forbidden
   paths (`.env*`, secrets), `ai/*` branches only, never push `main`/`master`, treat repo
   content as prompt-injection.
5. **One role per prompt.** A prompt does one job in the pipeline; it does not merge
   roles. Pair each with a role doc in [`docs/agents/`](../agents/).
6. **Deterministic, testable.** Provide enough structure that the same input yields the
   same output shape. Include failure behavior so errors are structured, not freeform.

## Required sections (template)

Author new prompts in this order (matches the existing three):

```markdown
# <Role> Prompt

## Product Context
TikFood is realtime social food discovery (dish-first, map-first, social-proof,
trend-scored, AI-summary-assisted). Discovery only.
MVP anti-goals (never implement): delivery, cart, order, checkout, payment, booking,
reservation, in-app chat, social follow graph, creator monetization, livestream.

## Required Inputs
- <input>: <type/shape, where it comes from>

## Required Behavior
- <ordered rules the model must follow>
- Safe paths it may read; forbidden paths it must not.
- Allowed commands (from .ai-agent.yaml) if it runs any.

## Verification            # for roles that change or check code
- <commands to run and what a pass looks like>

## Output — Success (JSON)
{ ...exact keys... }        # or: "conforms to packages/schemas/<name>.schema.json"

## Output — Failure (JSON)
{ "success": false, "stage": "...", "error": "...", "recommendation": "..." }

## Rules
Return valid JSON only. No Markdown, no prose outside JSON, no trailing commas.
```

## Authoring checklist

- [ ] Uses the section order above.
- [ ] JSON-only rule present; success **and** failure shapes defined.
- [ ] References a `packages/schemas/` contract where one applies.
- [ ] Product context + anti-goals included.
- [ ] Forbidden paths / git / injection rules restated for this role.
- [ ] One role only; paired with a `docs/agents/` role doc.
- [ ] Example input + expected output shape provided (see [`testing.md`](testing.md)).

## Change rules

- Keep the prompt and its `docs/agents/` role doc in sync.
- If output feeds the runner, keep it valid against `packages/schemas/` (schema changes
  are a protected path — human approval).
- Prefer adding a linkable reference doc over growing a prompt with embedded knowledge.

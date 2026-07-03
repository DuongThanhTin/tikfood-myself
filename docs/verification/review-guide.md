# Review Guide

> **Why this doc exists:** the repo has a `reviewer` runtime prompt
> ([`packages/prompts/reviewer.md`](../../packages/prompts/reviewer.md)) and a
> `code-reviewer` agent role, but no shared **human/Claude review guide** — what to look
> for, and how to rank what you find. This guide defines the review lens and a single
> severity rubric used everywhere (interactive review, the runner's reviewer, and PR
> feedback), so reviews are consistent and actionable.

**Related:** [`docs/verification/definition-of-done.md`](definition-of-done.md) (the
completion gate) · [`packages/prompts/reviewer.md`](../../packages/prompts/reviewer.md) ·
[`docs/agents/code-reviewer.md`](../agents/code-reviewer.md) ·
[`superpowers:requesting-code-review` / `receiving-code-review`].

## Review lenses (what to check)

Review against these, in priority order:

1. **Requirements & acceptance** — does it do what was asked, and only that?
2. **Product alignment & anti-goals** — no delivery/cart/checkout/booking/chat/etc.;
   stays discovery-only ([`docs/tikfood/anti-goals.md`](../tikfood/anti-goals.md)).
3. **Correctness** — logic, edge cases, error paths; both Postgres and fallback paths for
   search/geo/dish ([ADR-0004](../adr/0004-in-memory-fallback-repository.md)).
4. **Contracts** — `{data,error}` envelope intact; `lib/api.ts` in sync; schemas valid;
   no breaking API changes ([ADR-0003](../adr/0003-data-error-response-envelope.md)).
5. **Architecture** — layering respected, no forbidden edges
   ([ADR-0002](../adr/0002-handler-service-repository-layering.md)); change in the right
   place; reuse over new.
6. **Security** — input validation, parameterized SQL, no leaked SQL/secrets, gated areas
   have approval ([`application-security.md`](../standards/application-security.md)).
7. **Tests** — behavior change covered; tests meaningful, not tautological
   ([`testing.md`](../standards/testing.md)).
8. **Performance** — bounded queries, no N+1, no speculative complexity
   ([`performance.md`](../standards/performance.md)).
9. **Readability & naming** — clear verb names; matches surrounding style; no dead code.
10. **Honesty of claims** — PR/summary states what was actually verified vs. skipped.

## Severity rubric

Rank every finding (same scale as the Work Session Playbook):

| Severity | Meaning | Action |
| --- | --- | --- |
| **Critical** | Broken/unsafe: data loss, security hole, anti-goal, breaks a contract, wrong result | Must fix before merge |
| **Major** | Real bug or standard violation under plausible conditions; missing tests for risky logic | Should fix before merge |
| **Minor** | Limited-impact issue, small correctness/clarity gap | Fix now or file follow-up |
| **Nit** | Style/preference, non-blocking | Optional |

State severity, the file:line, why it matters, and a concrete suggestion. Do **not**
rewrite the author's code unless asked — propose the change.

## Reviewer conduct

- Be specific and evidence-based; cite the standard/ADR a finding violates.
- Separate "must" (Critical/Major) from "nice" (Minor/Nit) clearly.
- Verify claims rather than trusting the description; if you can't verify, say so.
- Reviewing ≠ implementing — keep review separate from the change (Golden Rule 4 of the
  Playbook).

## Receiving review

- Verify each point technically before acting; don't perform agreement or blindly apply
  suggestions (see `superpowers:receiving-code-review`).
- Push back with reasoning when a suggestion is wrong; otherwise fix and note how.

## For the automated reviewer

The runner's `reviewer` prompt must emit structured findings (boolean pass +
findings/fixes arrays) that map onto this rubric, and must keep output JSON-only per the
[Prompt Standard](../standards/prompt-engineering.md).

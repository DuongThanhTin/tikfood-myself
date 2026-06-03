# AI Implementation Rules

These rules are for Codex, Cursor, Claude, and future automation agents.

## Required Reading Before Coding

Read these first:

- `README.md`
- `.ai-agent.yaml`
- `docs/architecture.md`
- `docs/tikfood/**`
- `docs/standards/README.md`
- `docs/standards/project-structure.md`
- `docs/standards/dependency-management.md`
- `docs/standards/backend-architecture.md`
- `docs/standards/frontend-architecture.md`
- `docs/standards/api-contracts.md`
- `docs/standards/logging-observability.md`

If the request touches automation, also read:

- `docs/runner-contract.md`
- `docs/openai-integration-plan.md`
- `apps/ai-code-runner/README.md`

## Implementation Workflow

```text
understand request
-> check TikFood anti-goals
-> inspect relevant code/docs
-> choose smallest scoped change
-> update contracts/types/docs when needed
-> run checks
-> summarize changes, checks, risks, assumptions
```

## Request Handling

For every request, identify:

- target app: `apps/api`, `apps/web`, `apps/ai-code-runner`, docs, or workflow
- feature area: discovery, venue, search, trend, summary, ingestion, automation
- affected contract: API response, frontend type, runner schema, or no contract
- required checks
- anti-goal risk

## Backend Change Rules

- Keep handlers thin.
- Put domain logic under `internal/<domain>`.
- Keep response structs explicit.
- Use `snake_case` JSON.
- Add or update tests for endpoint/domain behavior.
- Run `npm run api:test`.

## Frontend Change Rules

- Keep route files focused on composition.
- Put reusable UI in `components`.
- Put API calls and shared types in `lib`.
- Keep API types in sync with backend responses.
- Run `npm run web:typecheck`.
- Run `npm run web:build` for production-facing UI changes.

## Dependency Rules

- Do not add dependencies casually.
- Explain why a new dependency is needed.
- Update lockfiles through package managers only.
- Do not run force dependency fixes without understanding breakage.

## Contract Rules

If changing backend response shape:

- Update frontend types.
- Update docs.
- Update schemas if the contract is public or runner-facing.
- Update examples if relevant.

## Logging Rules

- Add logs only when they improve debugging.
- Never log secrets.
- Prefer structured, small, non-sensitive context.

## Blocked MVP Scope

Do not implement:

- Delivery
- Cart
- Orders
- Checkout
- Payment
- Booking or reservations
- In-app chat
- Social follow graph
- Creator monetization
- Livestream

If the user asks for these in MVP, explain that it conflicts with TikFood scope and ask for explicit approval or defer it.

## Output Rules

When finishing a coding task, report:

- What changed
- Files touched
- Checks run
- Any skipped checks
- Risks or assumptions
- Next step if useful

# Engineering Standards

These standards are the source of truth for TikFood backend, frontend, automation, and AI-assisted implementation work.

Any AI coding assistant, including Codex, Cursor, Claude, or a future `ai-code-runner`, must read these files before changing application code:

- `docs/standards/project-structure.md`
- `docs/standards/dependency-management.md`
- `docs/standards/backend-architecture.md`
- `docs/standards/backend/README.md`
- `docs/standards/frontend-architecture.md`
- `docs/standards/api-contracts.md`
- `docs/standards/logging-observability.md`
- `docs/standards/ai-implementation-rules.md`

## Non-Negotiables

- TikFood is realtime social food discovery, not food delivery.
- Do not implement delivery, cart, order, checkout, payment, booking, reservation, in-app chat, social follow graph, creator monetization, or livestream logic in the MVP.
- Read `README.md`, `.ai-agent.yaml`, `docs/architecture.md`, `docs/tikfood/**`, and these standards before coding.
- Keep backend and frontend changes scoped to the relevant app.
- Keep request and response contracts consistent across API, web, schemas, and docs.
- Never read `.env`, `.env.*`, secrets, credentials, private keys, or tokens.
- Never claim tests passed unless they ran and passed.

## Current Monorepo Apps

```text
apps/api             Go backend for TikFood discovery
apps/web             Next.js frontend for TikFood discovery
apps/ai-code-runner  AI automation runner
```

## Backend-Specific Standards

Backend work must also read:

- `docs/standards/backend/structure.md`
- `docs/standards/backend/architecture.md`
- `docs/standards/backend/dependencies.md`
- `docs/standards/backend/request-response.md`
- `docs/standards/backend/errors.md`
- `docs/standards/backend/logging.md`
- `docs/standards/backend/testing.md`
- `docs/standards/backend/scaling.md`

## Expected AI Workflow

```text
read standards
-> understand feature request
-> inspect safe context
-> propose smallest implementation
-> change scoped files
-> run relevant checks
-> report assumptions, risks, and skipped checks honestly
```

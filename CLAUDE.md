# CLAUDE.md — TikFood Workspace

TikFood is realtime social food discovery (TikTok + Google Maps for food):
dish-first, map-first, social-proof-driven, trend-scored, geo-aware,
AI-summary-assisted. Discovery only.

## MVP Anti-Goals (block or require human approval)

Do not implement: delivery, cart, order, checkout, payment, booking,
reservation, in-app chat, social follow graph, creator monetization, livestream.

## Monorepo Map

- `apps/api` — Go backend (discovery API). See `apps/api/CLAUDE.md`.
- `apps/web` — Next.js frontend (discovery UX). See `apps/web/CLAUDE.md`.
- `apps/ai-code-runner` — automation runner skeleton (n8n → runner → PR).
- `packages/` — prompts, schemas, config policies.
- `docs/standards/` — engineering standards (background reference).
- `IMPROVEMENTS.md` — prioritized improvement roadmap.

## Source of Truth

When `docs/standards/*` conflicts with the actual code, the code wins. Flag the
mismatch instead of silently following stale docs. Read the relevant
`apps/*/CLAUDE.md` before changing that app.

## Security & Git

- Never read `.env`, `.env.*`, secrets, credentials, private keys, or tokens.
- Never push to `main` or `master`. Never force-push. Never auto-merge.
- Create PR branches under `ai/`. Active working branch is `dev`.
- Require human approval for auth, migrations, infra, external network calls,
  and any MVP anti-goal area.

## Naming (shared)

Use clear verb names: `CreateUser`, `FindByEmail`, `VenueDetail`.
Avoid `Do`, `Handle`, `Process`, `ExecuteX`, `Component1`.

## Before Finishing (shared)

- Reuse existing code; search before adding a new helper/service/component.
- No breaking API changes. Keep the response format unchanged.
- Update tests when you change behavior.

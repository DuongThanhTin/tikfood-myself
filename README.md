# TikFood AI Automation Starter

This repository is the workspace root for TikFood AI automation. It contains the local n8n workflow docs, an honest MVP `ai-code-runner` skeleton, runtime prompts, policy configs, schemas, examples, and a reusable starter template for a real TikFood product repo.

TikFood is realtime social food discovery: TikTok + Google Maps for food discovery. It is dish-first, map-first, social-proof-driven, trend-scored, geo-aware, and AI-summary-assisted.

TikFood is not a delivery, ordering, checkout, payment, booking, reservation, chat, creator monetization, livestream, or social follow graph app for the MVP.

## Workspace Layout

```text
apps/
  api/                  Go backend for TikFood discovery MVP
  web/                  Next.js frontend for TikFood discovery MVP
  ai-code-runner/       MVP runner skeleton and service README
packages/
  prompts/              Canonical runtime prompts used by the runner
  schemas/              OpenAPI and JSON Schema contracts
  config/               Workspace and TikFood agent policies
workflows/
  n8n/                  n8n setup, mapping, and example workflow docs
examples/
  feature-requests/     TikFood feature request payloads
starters/
  tikfood/              Docs/config starter for a real TikFood repo
docs/
  agents/               Human-readable agent docs
  tikfood/              Product scope, vision, and anti-goals
```

## Run Locally

Create `.env` from `.env.example`, then start the stack:

```bash
docker compose up --build
```

Local URLs:

- web: `http://localhost:3000`
- api: `http://localhost:18081`
- n8n: `http://localhost:5678`
- ai-code-runner: `http://localhost:8080`

Inside the Docker network, n8n calls:

```text
http://ai-code-runner:8080/jobs/feature
```

`AI_AGENT_URL` is optional for the MVP. The runner is expected to call OpenAI directly with `OPENAI_API_KEY` once the real implementation is added.

## Run Monorepo Apps

Backend:

```bash
npm run api:dev
```

Frontend:

```bash
npm run web:dev
```

Checks:

```bash
npm test
```

## Test A Feature Request

The runner currently exposes an MVP skeleton endpoint. It validates the shape enough to show the contract, then returns an honest failure response because clone/model/edit/push behavior is still TODO.

```bash
curl -X POST http://localhost:8080/jobs/feature \
  -H "Content-Type: application/json" \
  -d @examples/feature-requests/tikfood-realtime-map.json
```

## n8n To Runner Flow

```text
Manual Trigger
-> Set feature request JSON
-> HTTP Request POST /jobs/feature
-> IF success
-> GitHub Create Pull Request
-> Notification
```

n8n should create PRs from runner output. It must not merge PRs.

## TikFood App Skeleton

This repo now includes a real monorepo skeleton:

- `apps/api`: Go API with `GET /health` and `GET /api/v1/map/venues`
- `apps/web`: Next.js App Router frontend for discovery
- `apps/ai-code-runner`: automation runner

The current product app is still MVP scaffold code. It is discovery-only and intentionally avoids delivery, ordering, checkout, payment, booking, reservation, chat, creator monetization, livestream, and social follow graph scope.

## Security Rules

- Never read `.env`, `.env.*`, secrets, credentials, private keys, or tokens.
- Never push to `main` or `master`.
- Never force push.
- Never auto-merge.
- Create PR branches under `ai/`.
- Require human approval for protected paths, auth, migrations, infra, external network calls, and any MVP anti-goal area.

## MVP Anti-Goals

Do not implement delivery, cart, orders, checkout, payment, booking, reservations, in-app chat, social follow graph, creator monetization, or livestream features for the MVP.

# Local Development

## Prerequisites

- **Docker and Docker Compose** — the simplest path; `docker compose up` starts every
  service and auto-initializes the API database (no manual migrations needed).
- Only if you run an app **directly on the host** instead of via Docker:
  - **Go 1.23+** for `apps/api` (matches `apps/api/go.mod`).
  - **Node.js 20+** for `apps/web` / `apps/ai-code-runner`.
  - This repo is an **npm workspace** — run `npm install` **once at the repo root**
    before any `npm run ...`, `make verify-web`, or `make verify-runner`.
- A local `.env` created from `.env.example`
- Least-privilege GitHub credentials for push/PR automation later
- OpenAI credentials only when you are ready to enable model planning/editing

Do not commit `.env`.

## OpenAI Usage Status

OpenAI API usage is intentionally deferred for now to avoid API cost.

Current runner behavior does not require `OPENAI_API_KEY`. It can validate feature requests, clone a repo, create a local `ai/*` branch, and read safe repo context. It stops at `stage: "model"` before any paid model call.

When budget is available, follow `docs/openai-integration-plan.md` to enable planning and editing.

## Start The Stack

```bash
docker compose up --build
```

Services:

- n8n: `http://localhost:5678`
- ai-code-runner host URL: `http://localhost:8080`
- ai-code-runner Docker network URL: `http://ai-code-runner:8080`
- postgres internal service: `postgres:5432`
- TikFood API: `http://localhost:18081`
- TikFood API PostGIS: `localhost:15432`

Docker Compose wires the API container to `api-db:5432` by default. If you run `apps/api` directly on your host machine, use:

```bash
DATABASE_URL=postgres://tikfood:tikfood@localhost:15432/tikfood?sslmode=disable npm run api:dev
```

## Verify The Stack Is Up

```bash
curl http://localhost:18081/health                              # API    -> {"data":{"ok":true}}
curl http://localhost:8080/health                               # runner -> ok
curl http://localhost:3000                                      # web    -> HTML
curl "http://localhost:18081/api/v1/discovery/venues?limit=3"   # real discovery data (in the {data,error} envelope)
```

## Test The Runner Skeleton

```bash
curl -X POST http://localhost:8080/jobs/feature \
  -H "Content-Type: application/json" \
  -d @examples/feature-requests/tikfood-realtime-map.json
```

For a valid accessible repo, the MVP runner prepares a Git workspace, creates a local `ai/*` branch, reads safe repo context, then returns a failure response with `stage: "model"`. That is expected while OpenAI integration is deferred.

## Run ai-code-runner Outside Docker

```bash
cd apps/ai-code-runner
npm install
npm run build
npm start
```

Then call `http://localhost:8080/jobs/feature`.

## n8n Setup

Use `workflows/n8n/README.md` and `workflows/n8n/workflow.example.json` for the local MVP workflow.

## Troubleshooting

- **`npm run ...` fails with a missing-workspace or missing-module error** — run
  `npm install` at the **repo root** first (npm workspace; deps are hoisted there).
- **Port already in use** — the stack uses `3000` (web), `5678` (n8n), `8080` (runner),
  `15432` (API PostGIS), `18081` (API). Stop the conflicting process or run
  `docker compose down`, then retry.
- **`go: command not found`** when running `apps/api` on the host — install Go 1.23+,
  or just use `docker compose up` (no host Go needed).
- **API can't reach the database on the host** — use port **15432** (not 5432) with the
  `DATABASE_URL` shown above; `5432` belongs to the separate n8n Postgres.
- **Runner returns `stage: "model"`** — expected while OpenAI integration is deferred
  (see *OpenAI Usage Status* above).

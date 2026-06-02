# Local Development

## Prerequisites

- Docker and Docker Compose
- Node.js if you want to run `apps/ai-code-runner` outside Docker
- A local `.env` created from `.env.example`
- Least-privilege GitHub and OpenAI credentials for real runner work

Do not commit `.env`.

## Start The Stack

```bash
docker compose up --build
```

Services:

- n8n: `http://localhost:5678`
- ai-code-runner host URL: `http://localhost:8080`
- ai-code-runner Docker network URL: `http://ai-code-runner:8080`
- postgres internal service: `postgres:5432`

## Test The Runner Skeleton

```bash
curl -X POST http://localhost:8080/jobs/feature \
  -H "Content-Type: application/json" \
  -d @examples/feature-requests/tikfood-realtime-map.json
```

For a valid accessible repo, the MVP runner prepares a Git workspace, creates a local `ai/*` branch, then returns a failure response with `stage: "model"`. That is expected until context reading, OpenAI planning, guarded edits, review, commit, and push behavior are implemented.

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

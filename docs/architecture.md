# Architecture

> **Why this doc exists:** it describes both architectures in this repo — the
> **automation workspace** (turning scoped feature requests into reviewed PRs) and the
> **TikFood product** (discovery API + web). A previous version blended what is *built*
> with what is *planned*; this version separates **Current state** from **Roadmap** so
> the doc can be trusted. When it disagrees with code, code wins — flag the mismatch.

**Related:** [`docs/REPOSITORY-MAP.md`](REPOSITORY-MAP.md) (where everything lives) ·
[`docs/adr/`](adr/) (why the shape is what it is) · [`docs/standards/`](standards/)
(the binding implementation rules) · [`docs/contracts/runner.md`](contracts/runner.md).
Per-service contracts live in [`docs/services/`](services/).

---

## 1. Automation Workspace

Turns a scoped TikFood feature request into a reviewed GitHub pull request.

```text
n8n → ai-code-runner → repo clone → repo-context-reader → coding-agent
    → lint/test/build → reviewer → commit + push branch ai/* → n8n creates PR
```

- **n8n** owns workflow orchestration, PR creation, and notifications. It calls
  `POST /jobs/feature` at `http://ai-code-runner:8080/jobs/feature` inside Docker. It
  must **not** merge PRs.
- **ai-code-runner** owns input validation, repo cloning, branch creation
  (`ai/*` only), context reading, model calls, guarded file edits, allowlisted command
  execution, diff inspection, review, commit, and branch push.

**Current state (honest):** `ai-code-runner` is an **MVP skeleton**. The HTTP contract,
guards (`commandGuard`/`fileGuard`/`secretGuard`), and repo tools exist; model calls,
edits, verification, review, and push stages are **`TODO`**. Do not describe it as
production-ready. Contract: [`docs/services/ai-code-runner.md`](services/ai-code-runner.md) ·
[`docs/contracts/runner.md`](contracts/runner.md); schemas:
[`packages/schemas/`](../packages/schemas/).

**Roadmap:** implement the TODO stages (OpenAI integration —
see [`docs/openai-integration-plan.md`](openai-integration-plan.md)). All five pipeline
prompts now exist in [`packages/prompts/`](../packages/prompts/); wiring them into the
runner is the remaining work.

## 2. TikFood Product

Target product architecture the automation must respect:

```text
Social ingestion → PostgreSQL/PostGIS → Trend-scoring workers
    → AI-summary workers → Go API → Next.js frontend
```

### 2.1 Current state (built)

- **`apps/api`** — Go + Gin discovery API over PostgreSQL/PostGIS
  ([ADR-0001](adr/0001-gin-and-pgx-not-gorm.md)), layered
  `Handler → Service → Repository → postgres` ([ADR-0002](adr/0002-handler-service-repository-layering.md)),
  `{data,error}` envelope ([ADR-0003](adr/0003-data-error-response-envelope.md)), with an
  in-memory fallback repository for dev/CI
  ([ADR-0004](adr/0004-in-memory-fallback-repository.md)).
  Endpoints: `GET /health`, `GET /api/v1/discovery/venues`,
  `GET /api/v1/discovery/venues/:slug`, `GET /api/v1/map/venues`.
  Contract: [`docs/services/api.md`](services/api.md).
- **`apps/web`** — Next.js App Router + TypeScript discovery UX, MapLibre GL map, plain
  CSS ([ADR-0005](adr/0005-plain-css-frontend.md)), data via `lib/api.ts`.
  Contract: [`docs/services/web.md`](services/web.md).
- **Venue ingestion** — `apps/api/internal/ingest` + `cmd/ingest` provides an OSM/
  Overpass ingestion path (CLI). It is **not yet wired to an HTTP endpoint**.
- **Database** — PostgreSQL/PostGIS with SQL migrations in `apps/api/migrations`
  (including `locations`/`location_aliases`). Migrations require human approval.

### 2.2 Roadmap (not built — do not assume present)

- **Trend-scoring workers** and **AI-summary workers**: the `trend_scores` and
  `ai_summaries` tables exist, but **no worker populates them** — the realtime
  trend-scoring and AI "why trending" summaries are not implemented. There is no job
  scheduler/queue yet. (See [`IMPROVEMENTS.md`](../IMPROVEMENTS.md) P0.)
- **Social ingestion** beyond the OSM path.
- **Pagination `meta`**, `trace_id` in the envelope, readiness health check that pings
  the DB — all deferred (see `IMPROVEMENTS.md`).

## 3. Product guardrails

TikFood is realtime social food **discovery**, not delivery. Automation must block or
require human approval for delivery, cart, order, checkout, payment, booking,
reservation, in-app chat, social follow graph, creator monetization, and livestream
during the MVP.
→ [`docs/tikfood/anti-goals.md`](tikfood/anti-goals.md) · [`docs/ai/AI-CONTRACT.md`](ai/AI-CONTRACT.md).

## 4. Source of truth

Detailed implementation standards live in [`docs/standards/`](standards/) and are the
source of truth for dependency, logging, structure, backend/frontend architecture, and
API-contract decisions. The reasoning behind significant choices is in
[`docs/adr/`](adr/). This document describes shape; it does not restate the rules.

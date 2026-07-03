# Repository Map

> **Why this file exists:** the layout of this monorepo was previously described in
> pieces across `README.md`, `docs/standards/project-structure.md`,
> `docs/standards/backend/structure.md`, and three `CLAUDE.md` files. An AI agent (or
> a new human) had to trace several docs to answer "where does X live and what may I
> touch?". This is the **single directory → purpose inventory**. It does not restate
> rules — it points to the file that owns them. When the tree changes, update this map.

**Related:** [`docs/ai/AI-CONTRACT.md`](ai/AI-CONTRACT.md) (what you may/may not do) ·
[`docs/ai/CONTEXT-LOADING.md`](ai/CONTEXT-LOADING.md) (what to read, in what order) ·
[`docs/handbook/02-repository-design.md`](handbook/02-repository-design.md) (the *why*
behind the structure).

## Top level

| Path | Purpose | Read before editing |
| --- | --- | --- |
| `CLAUDE.md` | Root AI contract: anti-goals, security/git rules, naming, monorepo map | itself |
| `.ai-agent.yaml` | Machine-readable policy for `ai-code-runner` (protected paths, allowed/blocked commands, docs-required) | itself |
| `README.md` | Human entry point: what the workspace is, how to run it | — |
| `IMPROVEMENTS.md` | Prioritized roadmap of code/infra findings (a roadmap, not rules) | — |
| `docker-compose.yml` | Local stack (web, api, n8n, runner). **Protected path** | requires human approval |
| `.env.example` | Template for local env. Never read real `.env*` | never read `.env*` |
| `apps/` | Runnable applications (product + automation) | per-app `CLAUDE.md` |
| `packages/` | Shared, non-app assets: prompts, schemas, config | — |
| `docs/` | All documentation (rules, standards, handbook, this map) | — |
| `workflows/` | n8n workflow documentation | — |
| `examples/` | Sample feature-request payloads for the runner | — |
| `starters/` | Copy-out template for a *real* TikFood product repo | — |

## `apps/` — applications

### `apps/api` — Go discovery backend
Owner rules: [`apps/api/CLAUDE.md`](../apps/api/CLAUDE.md) +
[`docs/standards/backend/`](standards/backend/). Layering:
`http (Handler) → discovery (Service) → Repository → postgres`, one-way only.

| Path | Purpose |
| --- | --- |
| `cmd/server/` | Process entrypoint (`main.go`) for the API server |
| `cmd/ingest/` | CLI entrypoint for the OSM venue ingestion command |
| `internal/app/` | Dependency-injection root (`Container`); picks Postgres vs in-memory repo by `DATABASE_URL` |
| `internal/config/` | Environment-based configuration loading |
| `internal/http/` | Gin router, `*Handler`s, request parsing/validation, response envelopes, middleware, error codes |
| `internal/discovery/` | Discovery **domain**: models, `VenueService` (business logic), search/normalization, in-memory fallback repository |
| `internal/ingest/` | Venue **ingestion** domain: service, `SourcePlace→Venue` mapper, OSM/Overpass client, opening-hours parser |
| `internal/platform/textutil/` | Bounded infrastructure: text normalization helpers |
| `internal/storage/postgres/` | PostgreSQL/PostGIS repository implementations (`database/sql` via pgx stdlib) |
| `migrations/` | SQL schema migrations. **Migrations require human approval** |
| `seeds/` | Seed data |
| `db/` | Database support assets |

### `apps/web` — Next.js discovery frontend
Owner rules: [`apps/web/CLAUDE.md`](../apps/web/CLAUDE.md) +
[`docs/standards/frontend-architecture.md`](standards/frontend-architecture.md).

| Path | Purpose |
| --- | --- |
| `app/` | App Router: `layout.tsx`, `page.tsx` (server component passing `initialVenues`), `globals.css` (all styling, plain CSS) |
| `components/` | React components (`DiscoveryExperience.tsx` — large, split approved but not done; `VenueList.tsx`) |
| `lib/` | `api.ts` — the **only** place raw `fetch` calls live; frontend types kept in sync with Go JSON tags |
| `public/` | Static assets |

### `apps/ai-code-runner` — automation runner (MVP skeleton)
Owner: [`apps/ai-code-runner/README.md`](../apps/ai-code-runner/README.md) +
contract in [`docs/runner-contract.md`](runner-contract.md). Incomplete stages are
marked `TODO`; do not claim production-ready.

| Path | Purpose |
| --- | --- |
| `src/server.ts` | HTTP entry (`POST /jobs/feature`) |
| `src/jobs/runFeatureJob.ts` | Orchestrates the feature job pipeline |
| `src/guards/` | Policy guards: `commandGuard`, `fileGuard`, `secretGuard` |
| `src/tools/` | Repo tools: git workspace/diff, read/write file, search, run command, repo context |
| `src/utils/` | Pure helpers (e.g. `slugify`) |
| `prompts/` | Runner-local prompt notes (canonical prompts live in `packages/prompts/`) |

## `packages/` — shared assets

| Path | Purpose |
| --- | --- |
| `packages/prompts/` | Canonical **runtime** prompts: `coding-agent.md`, `repo-context-reader.md`, `reviewer.md` (JSON-only output). See [`docs/agents/`](agents/) for their human-readable role docs |
| `packages/schemas/` | Data contracts: `feature-request.schema.json`, `runner-success-response.schema.json`, `runner-failure-response.schema.json`, `openapi.yaml`. **Protected path** |
| `packages/config/` | Agent policies: `default.ai-agent.yaml` (generic baseline), `tikfood.ai-agent.yaml` (product policy + anti-goals) |

## `docs/` — documentation

| Path | Purpose |
| --- | --- |
| `docs/ai/` | AI-OS spine: `AI-CONTRACT.md`, `CONTEXT-LOADING.md`, `ROADMAP.md` |
| `docs/standards/` | Engineering standards. `backend/` is deep (13 files); frontend/api/deps/logging at root |
| `docs/agents/` | Human-readable agent role docs (pair with `packages/prompts/`) |
| `docs/tikfood/` | Product source of truth: `product-vision.md`, `mvp-scope.md`, `anti-goals.md` |
| `docs/handbook/` | Bilingual onboarding handbook (Parts 1–2 complete; 3–8 being filled) |
| `docs/superpowers/` | Real spec-driven-development trail: `specs/` (designs) + `plans/` (implementation plans) |
| `docs/architecture.md` | System architecture (automation + product) |
| `docs/security.md` | **Runner/automation** security policy (distinct from product app security) |
| `docs/runner-contract.md` | `POST /jobs/feature` request/response contract + runner rules |
| `docs/local-development.md` | Local dev setup |
| `docs/roadmap.md` / `docs/openai-integration-plan.md` / `docs/restructure-prompt.md` | Product roadmap / OpenAI plan / historical restructure spec |
| `docs/REPOSITORY-MAP.md` | This file |

**Future directories (forecast in handbook §2.5, created by the AI-OS roadmap):**
`docs/adr/`, `docs/services/`, `docs/thinking/`, `docs/verification/`,
`docs/recipes/`. Marked here so the map matches [`docs/ai/ROADMAP.md`](ai/ROADMAP.md);
each is created only when its phase runs, never as an empty placeholder.

## Other top-level dirs

| Path | Purpose |
| --- | --- |
| `workflows/n8n/` | n8n setup, mapping, example workflow (`ai-feature-to-pr-workflow.md`, `workflow.example.json`) |
| `examples/feature-requests/` | Sample runner payloads (`tikfood-ai-summary`, `tikfood-dish-search`, `tikfood-realtime-map`) |
| `starters/tikfood/` | Seed for a real product repo: `.ai-agent.yaml`, `project-docs/` (VISION/REQUIREMENTS/ARCHITECTURE/TASKBOARD), `automation/`. Ships **no app source** |
| `.superpowers/sdd/` | Spec-driven-development working state (task briefs/reports/progress) |

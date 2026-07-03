# TikFood AI Automation Workspace Restructure Prompt

You are a senior monorepo architect, AI automation engineer, and technical documentation engineer.

## Current Repository

I currently have a repository named `tikfood-ai-automation-starter`.

Inside it, there is a nested folder:

```text
ai-automation-starter/
|-- docs/
|-- prompts/
|-- .ai-agent.yaml
|-- .gitignore
|-- docker-compose.yml
|-- Dockerfile
|-- feature-request.example.json
|-- n8n-example-mapping.md
|-- README-AI-AUTOMATION-USAGE.md
`-- README.md
```

I want to turn `tikfood-ai-automation-starter` into the main workspace root.

This repository should become a clean, production-minded workspace for building and maintaining AI automation for my TikFood project.

## TikFood Product Context

TikFood - realtime social food discovery

Tagline: TikTok + Google Maps for food discovery

TikFood is:

- Dish-first discovery
- Map-first discovery
- Social proof aggregation
- Realtime trend scoring
- AI-generated summaries
- Hidden gem detection
- Food graph, creator graph, trend graph, and geo intelligence

TikFood is not a food delivery app.

MVP anti-goals:

- Food delivery
- Cart
- Orders
- Checkout
- Payment
- Booking or reservations
- In-app chat
- Social network follow graph
- Creator monetization
- Livestream

## Goal

Restructure the current repo into a workspace layout.

Requirements:

- Do not delete useful content.
- Use `git mv` where possible.
- If a file is outdated, update it.
- If a file is duplicated, consolidate it and leave clear references.
- Make the final repo easy to understand, easy to run locally, and ready for future expansion.

## Target Workspace Structure

Create or update the repository so the final structure is:

```text
tikfood-ai-automation-starter/
|-- apps/
|   `-- ai-code-runner/
|       |-- src/
|       |   |-- server.ts
|       |   |-- jobs/
|       |   |-- tools/
|       |   |-- guards/
|       |   `-- utils/
|       |-- prompts/
|       |-- package.json
|       |-- tsconfig.json
|       |-- Dockerfile
|       `-- README.md
|
|-- packages/
|   |-- prompts/
|   |   |-- coding-agent.md
|   |   |-- repo-context-reader.md
|   |   `-- reviewer.md
|   |
|   |-- schemas/
|   |   |-- openapi.yaml
|   |   |-- feature-request.schema.json
|   |   |-- runner-success-response.schema.json
|   |   `-- runner-failure-response.schema.json
|   |
|   `-- config/
|       |-- tikfood.ai-agent.yaml
|       `-- default.ai-agent.yaml
|
|-- workflows/
|   `-- n8n/
|       |-- README.md
|       |-- ai-feature-to-pr-workflow.md
|       |-- n8n-example-mapping.md
|       `-- workflow.example.json
|
|-- examples/
|   `-- feature-requests/
|       |-- tikfood-realtime-map.json
|       |-- tikfood-ai-summary.json
|       `-- tikfood-dish-search.json
|
|-- starters/
|   `-- tikfood/
|       |-- README.md
|       |-- .ai-agent.yaml
|       |-- project-docs/
|       |   |-- VISION.md
|       |   |-- REQUIREMENTS.md
|       |   |-- ARCHITECTURE.md
|       |   `-- TASKBOARD.md
|       `-- automation/
|           |-- prompts/
|           |-- feature-request.example.json
|           `-- n8n-example-mapping.md
|
|-- docs/
|   |-- architecture.md
|   |-- runner-contract.md
|   |-- security.md
|   |-- roadmap.md
|   |-- local-development.md
|   |-- agents/
|   |   |-- feature-analyzer.md
|   |   |-- repo-context-reader.md
|   |   |-- coding-agent.md
|   |   |-- code-reviewer.md
|   |   `-- pull-request-writer.md
|   `-- tikfood/
|       |-- product-vision.md
|       |-- mvp-scope.md
|       `-- anti-goals.md
|
|-- docker-compose.yml
|-- .env.example
|-- .gitignore
|-- .ai-agent.yaml
|-- package.json
`-- README.md
```

## Important Restructure Rules

- The current `ai-automation-starter/` folder should not remain as the main nested folder.
- Move its useful contents into the new workspace layout.
- After restructuring, the root `README.md` should be the main entry point.
- Keep runtime prompts in `packages/prompts/`.
- Keep human-readable agent docs in `docs/agents/`.
- Keep n8n workflow docs in `workflows/n8n/`.
- Keep feature request examples in `examples/feature-requests/`.
- Keep TikFood reusable starter files in `starters/tikfood/`.
- Keep runner contract and OpenAPI files in `packages/schemas/`.
- Reference schema files from `docs/runner-contract.md`.
- Keep root `docker-compose.yml` as the main local development compose file.

## Root README Requirements

Update root `README.md` so it explains:

- What this workspace is
- How it supports TikFood AI automation
- Workspace folder structure
- What is inside `apps/`
- What is inside `packages/`
- What is inside `workflows/`
- What is inside `examples/`
- What is inside `starters/`
- How to run locally
- How n8n calls `ai-code-runner`
- How to test with example feature requests
- How to use this starter inside a real TikFood repo
- Security rules
- MVP anti-goals

## Root Docker Compose Requirements

Update root `docker-compose.yml`.

It should include:

- `postgres`
- `n8n`
- `ai-code-runner`

Local URLs:

- n8n: `http://localhost:5678`
- ai-code-runner: `http://localhost:8080`

Inside the Docker network, n8n must call the runner using:

```text
http://ai-code-runner:8080/jobs/feature
```

Explain that `AI_AGENT_URL` is optional.

For the MVP:

```env
AI_AGENT_URL=
```

Reason: `ai-code-runner` calls OpenAI directly with `OPENAI_API_KEY`.

## Root .env.example

Create or update `.env.example` with:

```env
POSTGRES_USER=n8n
POSTGRES_PASSWORD=change_this_password
POSTGRES_DB=n8n

N8N_ENCRYPTION_KEY=change_this_to_a_long_random_string
N8N_HOST=localhost
N8N_PROTOCOL=http
WEBHOOK_URL=http://localhost:5678/
GENERIC_TIMEZONE=Asia/Ho_Chi_Minh

GITHUB_TOKEN=your_github_token_here

OPENAI_API_KEY=your_openai_api_key_here
OPENAI_MODEL=gpt-4.1

AI_AGENT_URL=
```

Do not create real secrets.

## Root .ai-agent.yaml

Update root `.ai-agent.yaml`.

It should describe this repository itself, not the TikFood app.

It should include:

```yaml
repo_name: tikfood-ai-automation-starter
repo_type: ai_automation_workspace
```

It should define:

- `docs_required`
- `protected_paths`
- `allowed_commands`
- `blocked_commands`
- `coding_rules`
- `testing_rules`
- `security_rules`

It should also say this workspace is for automating TikFood feature implementation, and TikFood anti-goals must be preserved in TikFood-specific prompts and configs.

## packages/config/tikfood.ai-agent.yaml

Create a TikFood-specific config.

It must include:

- TikFood product positioning
- North star
- Principles
- MVP anti-goals
- Backend stack
- Frontend stack
- `docs_required`
- `protected_paths`
- `approval_required_for`
- `allowed_commands`
- `blocked_commands`
- `coding_rules`
- `testing_rules`

It must explicitly prevent delivery, cart, order, payment, and booking logic in the MVP.

## packages/prompts/

Create the canonical runtime prompts:

- `packages/prompts/coding-agent.md`
- `packages/prompts/repo-context-reader.md`
- `packages/prompts/reviewer.md`

These are used by `ai-code-runner`.

They must be TikFood-aware but generic enough to work on a real TikFood repo.

They must enforce:

- Read docs before coding.
- Read `.ai-agent.yaml`.
- Understand TikFood product direction.
- Do not implement MVP anti-goals.
- Do not read secrets.
- Do not push `main` or `master`.
- Only create branches under `ai/`.
- Modify the smallest number of files possible.
- Run relevant checks.
- Return valid JSON only.

## docs/agents/

Move or rewrite these existing files:

- `docs/agents-code-reviewer.md`
- `docs/agents-feature-analyzer.md`
- `docs/agents-pull-request-writer.md`
- `docs/agents-repo-context-reader.md`
- `docs/agents-repository-aware-coding.md`

Into:

- `docs/agents/code-reviewer.md`
- `docs/agents/feature-analyzer.md`
- `docs/agents/pull-request-writer.md`
- `docs/agents/repo-context-reader.md`
- `docs/agents/coding-agent.md`

These docs are human-readable. They should explain:

- Agent purpose
- Inputs
- Outputs
- Runtime prompt location
- JSON contract
- Failure behavior
- Security rules

## packages/schemas/

Move or update `docs/openapi.yaml` into:

- `packages/schemas/openapi.yaml`

Also create JSON schemas:

- `packages/schemas/feature-request.schema.json`
- `packages/schemas/runner-success-response.schema.json`
- `packages/schemas/runner-failure-response.schema.json`

OpenAPI must define:

```text
POST /jobs/feature
```

Request fields:

- `feature_id`
- `repo`
- `base_branch`
- `title`
- `description`
- `acceptance_criteria`
- `mode`

Success response fields:

- `success`
- `feature_id`
- `repo`
- `branch`
- `commit_sha`
- `summary`
- `files_changed`
- `tests_run`
- `test_result`
- `risks`
- `pr_title`
- `pr_body`

Failure response fields:

- `success`
- `feature_id`
- `stage`
- `error`
- `logs`
- `recommendation`

## apps/ai-code-runner/

If the current repository only has a placeholder Dockerfile and no real runner source code, create a clear MVP skeleton for the runner.

It should include:

```text
apps/ai-code-runner/
|-- src/
|   |-- server.ts
|   |-- jobs/runFeatureJob.ts
|   |-- tools/readFile.ts
|   |-- tools/writeFile.ts
|   |-- tools/searchRepo.ts
|   |-- tools/runCommand.ts
|   |-- tools/gitDiff.ts
|   |-- guards/fileGuard.ts
|   |-- guards/commandGuard.ts
|   |-- guards/secretGuard.ts
|   `-- utils/slugify.ts
|-- prompts/
|-- package.json
|-- tsconfig.json
|-- Dockerfile
`-- README.md
```

The MVP runner can be a skeleton, but it must be honest.

Rules:

- If a function is not fully implemented, mark it clearly as `TODO`.
- Do not create fake production claims.

Runner must eventually implement:

- `POST /jobs/feature`
- Validate input.
- Clone repo.
- Create branch `ai/{feature_id}-{slug}`.
- Read repo context.
- Call model using OpenAI API.
- Run allowed commands only.
- Block protected paths.
- Inspect git diff.
- Commit and push `ai/*` branch only.
- Return JSON matching OpenAPI.

## workflows/n8n/

Move or rewrite:

- `n8n-example-mapping.md`
- `docs/n8n-workflow.md`
- `docs/n8n-ai-feature-to-pr-workflow.md`

Into:

- `workflows/n8n/README.md`
- `workflows/n8n/ai-feature-to-pr-workflow.md`
- `workflows/n8n/n8n-example-mapping.md`
- `workflows/n8n/workflow.example.json`

Docs must include:

- Manual Trigger setup
- Set node JSON
- HTTP Request node config
- IF condition
- GitHub PR node config
- Notification examples
- Error handling
- Docker network URL
- Local host URL
- How to test

## examples/feature-requests/

Move and update `feature-request.example.json`.

Create these examples:

- `examples/feature-requests/tikfood-realtime-map.json`
- `examples/feature-requests/tikfood-ai-summary.json`
- `examples/feature-requests/tikfood-dish-search.json`

Examples must be TikFood-specific and must include anti-goal criteria such as:

```text
No delivery, cart, order, payment, or booking logic is added.
```

## starters/tikfood/

Create a reusable starter template that can be copied into the real TikFood repo.

It should include:

```text
starters/tikfood/
|-- README.md
|-- .ai-agent.yaml
|-- project-docs/
|   |-- VISION.md
|   |-- REQUIREMENTS.md
|   |-- ARCHITECTURE.md
|   `-- TASKBOARD.md
`-- automation/
    |-- prompts/
    |-- feature-request.example.json
    `-- n8n-example-mapping.md
```

This folder should not contain full app source code.

It is a starter for docs, config, and prompts only.

## docs/architecture.md

Update it to explain both the automation workspace architecture and the TikFood application architecture.

Automation workspace architecture:

```text
n8n
-> ai-code-runner
-> repo clone
-> repo-context-reader
-> coding-agent
-> lint/test/build
-> reviewer
-> commit + push branch
-> n8n creates PR
```

TikFood application architecture:

```text
Social ingestion
-> PostgreSQL/PostGIS
-> Trend scoring workers
-> AI summary workers
-> Go API
-> Next.js frontend
```

## docs/security.md

Update with:

- GitHub token least privilege
- Prefer GitHub App later
- No secrets reading
- No `.env` reading
- No dangerous shell commands
- No push to `main` or `master`
- No force push
- No auto merge
- Human approval required for protected areas
- Sandbox execution
- Audit logs
- Prompt injection risks from repo docs and source files
- Social platform scraping and compliance caution

## docs/roadmap.md

Update roadmap.

Phase 1:

- Workspace restructure
- Local n8n
- Runner skeleton
- Manual Trigger
- Test example request

Phase 2:

- Real runner implementation
- GitHub branch push
- PR creation
- Notifications

Phase 3:

- GitHub Issue label `ai-ready`
- Improved context reading
- Reviewer hardening
- Cost tracking

Phase 4:

- Queue
- Metrics
- Audit logs
- GitHub App
- Policy engine
- Production hardening

## .gitignore

Update `.gitignore` with:

```gitignore
.env
.env.*
node_modules/
dist/
build/
coverage/
tmp/
logs/
*.log
.DS_Store
.vscode/
.idea/
runner_workspace/
.ai-code-runner/
```

## Consistency Requirements

After restructuring, make sure there are no confusing duplicate docs.

Every file should agree that:

- TikFood is realtime social food discovery.
- TikFood is not delivery, order, or payment.
- Backend is Go.
- Frontend is Next.js TypeScript.
- Automation is n8n + ai-code-runner.
- AI must read docs before coding.
- AI must return JSON.
- AI must never read secrets.
- AI must never push `main` or `master`.
- AI must create PRs, not merge.

## Definition of Done

The task is complete only when:

- The nested `ai-automation-starter/` folder is no longer the main workspace container.
- Useful existing content has been moved or consolidated.
- Root `README.md` is the clear entry point.
- Root `docker-compose.yml` runs postgres, n8n, and ai-code-runner locally.
- `.env.example` contains placeholders only.
- Runtime prompts, human-readable docs, schemas, workflows, examples, and starter files are in their target locations.
- The runner skeleton is clearly labeled as an MVP skeleton if it is not fully implemented.
- There are no fake production-readiness claims.
- No real secrets are created.
- The final file tree matches the target structure as closely as practical.

## Output After Completing

After updating files, print:

- Summary of changes
- Final folder tree
- Files moved
- Files created
- Files removed, if any
- Assumptions made
- Next commands to run locally

Do not create real secrets.

Do not claim the runner is production-ready if it is only a skeleton.

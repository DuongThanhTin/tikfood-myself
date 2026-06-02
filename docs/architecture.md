# Architecture

This workspace describes automation for turning scoped TikFood feature requests into reviewed GitHub pull requests. It also documents the target TikFood product architecture the automation must respect.

## Automation Workspace Architecture

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

n8n owns workflow orchestration, PR creation, and notifications. It calls `POST /jobs/feature` at `http://ai-code-runner:8080/jobs/feature` inside Docker.

`ai-code-runner` owns validation, repository cloning, branch creation, context reading, model calls, guarded file edits, allowlisted command execution, diff inspection, review, commit, and branch push. The current implementation is an MVP skeleton and marks incomplete stages as TODO.

## TikFood Application Architecture

```text
Social ingestion
-> PostgreSQL/PostGIS
-> Trend scoring workers
-> AI summary workers
-> Go API
-> Next.js frontend
```

The monorepo now includes an MVP app skeleton:

- `apps/api`: Go API for discovery endpoints.
- `apps/web`: Next.js TypeScript frontend for discovery UX.
- `apps/ai-code-runner`: automation runner.

The future product architecture should add PostgreSQL/PostGIS for MVP geo data, PostgreSQL full-text search for MVP search, and worker processes for ingestion, trend scoring, and AI summaries.

## Product Guardrails

TikFood is realtime social food discovery, not delivery. Automation must block or require human approval for delivery, cart, order, checkout, payment, booking, reservation, in-app chat, social follow graph, creator monetization, and livestream scope during the MVP.

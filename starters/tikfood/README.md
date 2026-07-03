# TikFood Repo Starter

Copy this folder into a real TikFood product repository to give AI automation enough context to work safely.

This starter contains docs, config, prompts references, and example request mapping only. It does not contain full application source code.

To make AI agents work like senior engineers in the target repo, follow `project-docs/AI-OS.md` — the checklist for standing up the AI Operating System (rules, context, decisions, standards, recipes), using the automation workspace `docs/` as the reference implementation.

Recommended target stack:

- Backend: Go, PostgreSQL/PostGIS, workers for ingestion, trend scoring, and AI summaries
- Frontend: Next.js App Router, TypeScript, Tailwind CSS, map UI, typed API client
- Automation: n8n + ai-code-runner

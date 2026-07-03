# Roadmap — Automation Build Phases

> **Scope:** the phased build-out of the automation workspace (runner, OpenAI, PR flow,
> production hardening). This is one of three planning docs — don't confuse them:
> **this file** = automation build phases · [`IMPROVEMENTS.md`](../IMPROVEMENTS.md) =
> prioritized code/infra findings · [`docs/ai/ROADMAP.md`](ai/ROADMAP.md) = the AI
> Operating System documentation build.

## Phase 1

- Workspace restructure
- Local n8n
- Runner skeleton
- Manual Trigger
- Test example request
- TikFood monorepo skeleton with Go API and Next.js web
- Git workspace preparation
- Safe repo context reading
- Keep OpenAI integration deferred to avoid API cost

## Phase 2

- Optional OpenAI planning step
- Optional OpenAI editing step
- Guarded file edits
- Verification commands
- GitHub branch push
- PR creation
- Notifications

## Phase 3

- GitHub Issue label `ai-ready`
- Improved context reading
- Reviewer hardening
- Cost tracking

## Phase 4

- Queue
- Metrics
- Audit logs
- GitHub App
- Policy engine
- Production hardening

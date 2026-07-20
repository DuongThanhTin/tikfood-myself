---
description: Drive a feature from request to a clean PR using the repo's feature workflow (classify -> context -> brainstorm -> spec/plan if large -> branch -> implement -> test -> verify -> docs -> PR), honoring every approval gate.
argument-hint: [short feature description]
---

# /feature — request to PR

Authoritative workflow (steps, scope rules, and every 🚦 gate): **@docs/workflows/feature.md** — follow it top to bottom.

Feature request: $ARGUMENTS

## Run the workflow
1. **Classify** the task; if it touches auth / migration / infra / external / anti-goal -> @docs/workflows/gated-change.md **first**.
2. **Load context** — @docs/ai/CONTEXT-LOADING.md §1 core + §2 matrix (or dispatch the `repo-context-reader` subagent).
3. **Brainstorm** to agree scope/acceptance (skill `brainstorming`). Large feature? Write a spec, then a plan, stopping for approval at each.
4. **Branch** `ai/<feature>` (never off `main`).
5. **Implement** via the matching skill/recipe; **test**; **verify** with `make verify`.
6. **Update docs/ADR**; open a PR under `ai/` — do not merge, do not push `main`, do not force-push.

## Guard rails
- Stop for human approval on any gated area (@CLAUDE.md; @docs/ai/AI-CONTRACT.md §3–§4).
- Claim done only after `make verify` actually ran — @docs/verification/definition-of-done.md.

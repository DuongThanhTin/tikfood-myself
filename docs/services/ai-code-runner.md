# Service Contract — `apps/ai-code-runner` (Automation Runner)

> **Why this doc exists:** a single at-a-glance contract for the automation runner — its
> job endpoint, the guarantees it must honor, its dependencies, and its **honest
> current state (MVP skeleton)** — so nobody mistakes the skeleton for a finished
> service. The wire contract lives in [`docs/runner-contract.md`](../runner-contract.md)
> and [`packages/schemas/`](../../packages/schemas/); policy lives in
> [`.ai-agent.yaml`](../../.ai-agent.yaml) and [`packages/config/`](../../packages/config/).

## Purpose

Turn a scoped, validated feature request into a reviewed PR branch (`ai/*`), driven by
n8n. It **must not** merge PRs and must never push `main`/`master`.

## Interface

- **Endpoint:** `POST /jobs/feature` (`src/server.ts`).
- **Request:** conforms to
  [`packages/schemas/feature-request.schema.json`](../../packages/schemas/feature-request.schema.json)
  — `feature_id`, `repo`, `base_branch`, `title`, `description`,
  `acceptance_criteria[]`, `mode` (`feature|review|dry_run`).
- **Response (success):**
  [`runner-success-response.schema.json`](../../packages/schemas/runner-success-response.schema.json)
  — `success:true`, branch (must match `^ai/`), `commit_sha`, `summary`,
  `files_changed[]`, `tests_run[]`, `test_result`, `risks[]`, `pr_title`, `pr_body`.
- **Response (failure):**
  [`runner-failure-response.schema.json`](../../packages/schemas/runner-failure-response.schema.json)
  — `success:false`, `stage`, `error`, `logs`, `recommendation`.

Full behavior spec: [`docs/runner-contract.md`](../runner-contract.md).

## Dependencies

- **n8n** (caller/orchestrator; owns PR creation + notifications).
- **Runtime prompts** in [`packages/prompts/`](../../packages/prompts/) — the 5-stage
  pipeline `feature-analyzer → repo-context-reader → coding-agent → reviewer →
  pull-request-writer` (JSON-only output; index:
  [`packages/prompts/README.md`](../../packages/prompts/README.md)). Role docs:
  [`docs/agents/`](../agents/).
- **Policy** from [`.ai-agent.yaml`](../../.ai-agent.yaml) /
  [`packages/config/tikfood.ai-agent.yaml`](../../packages/config/tikfood.ai-agent.yaml)
  (protected paths, allowed/blocked commands, anti-goals).
- **Model provider** (OpenAI) — planned; see
  [`docs/openai-integration-plan.md`](../openai-integration-plan.md).
- **Guards** (`src/guards/`): `commandGuard`, `fileGuard`, `secretGuard`.

## Consumers

n8n (`ai-feature-to-pr` workflow). See [`workflows/n8n/`](../../workflows/n8n/).

## Change rules

- Enforce the hard contract: validate input first; branch under `ai/` only; never push
  `main`/`master`; never force-push; never read secrets; never claim success unless the
  branch was created, verified, committed, and pushed
  ([`docs/ai/AI-CONTRACT.md`](../ai/AI-CONTRACT.md), [`docs/security.md`](../security.md)).
- Responses **must** validate against the JSON schemas — changing a schema is a
  **protected-path** change (human approval).
- Mark unfinished stages clearly as `TODO`; make no production claims.
- Runtime prompts must return valid JSON only (see [`docs/standards/prompt-engineering.md`](../standards/prompt-engineering.md)).

## Current state — MVP skeleton (honest)

Present: HTTP entry, guards, and repo tools (`gitWorkspace`, `gitDiff`, `readFile`,
`writeFile`, `searchRepo`, `runCommand`, `repoContext`), and all 5 pipeline prompts in
`packages/prompts/`. **`TODO`:** model calls, guarded edits, verification, review,
commit, and push — i.e. wiring the prompts into the runner. See
[`docs/architecture.md`](../architecture.md) §1.

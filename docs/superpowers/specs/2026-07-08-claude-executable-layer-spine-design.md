# Design — `.claude/` Executable Layer (Spine, Batch 1)

**Date:** 2026-07-08 · **Branch:** `ai/claude-executable-layer` · **Status:** approved for planning

## Problem

The repo has a mature, documentation-based AI Operating System — rules
(`CLAUDE.md`, `docs/ai/AI-CONTRACT.md`), navigation (`REPOSITORY-MAP.md`,
`CONTEXT-LOADING.md`), decision memory (8 ADRs), standards, recipes, workflows,
and agent descriptions. All 8 phases of `docs/ai/ROADMAP.md` are complete.

But every one of those assets is **passive Markdown**: Claude only uses a recipe
if a human remembers to point at it. There is **no `.claude/` directory** — no
native Claude Code skills, agents, or commands that *self-trigger*. The gap is
not content; it is the missing **executable layer** that turns existing docs into
behavior Claude invokes automatically (the way Superpowers skills fire from their
`name`/`description` frontmatter).

This spec covers **Batch 1: the spine** — a minimal vertical slice that proves the
authoring pattern end to end (one skill wrapped, one skill net-new, one subagent,
one command) plus the conventions doc that makes every later batch a cheap,
consistent repeat. It does **not** wrap the full catalog; that is deferred to
future batches.

## Non-goals (this batch)

- Wrapping all 8 recipes / 5 agents / 4 workflows (future batches).
- Authoring every net-new skill the user named (`add-error-code`,
  `backend-service-change`) — only `contract-first-change` is in scope as the
  net-new proof.
- A `code-reviewer` agent — three review tools already exist (`/code-review`,
  `/review`, Superpowers `code-review`); a fourth would collide.
- Changing any existing doc's content, any `docs/` policy, or any repo rule.

## Architecture — how `.claude/` relates to `docs/`

Thin wrappers. The **source of truth stays in `docs/`**; each `.claude/` artifact
is a trigger + short checklist + `@docs/...` pointer, with zero procedure
duplication (honors the ROADMAP principle *"index/link, never duplicate"*).

| `.claude/` artifact | Wraps | Rationale |
|---|---|---|
| `skills/*/SKILL.md` | `docs/recipes/*` | Auto-triggers on a task description; procedural |
| `agents/*.md` | `docs/agents/*` | Runs in isolated context; returns a result, not chatter |
| `commands/*.md` | `docs/workflows/*` | Human-invoked orchestration (e.g. `/feature`) |
| *(stays reference)* | `docs/ai/*`, `docs/thinking`, `docs/verification/*` | Always-loaded context, not triggered |

## The spine — five artifacts

1. **`.claude/skills/add-api-endpoint/SKILL.md`** — wraps
   `docs/recipes/add-api-endpoint.md`. Triggers when adding/extending an
   `apps/api` endpoint. *Proves: wrapping an existing recipe.*

2. **`.claude/skills/contract-first-change/SKILL.md`** — **net-new**. Triggers
   when a change touches the `{data,error}` envelope, `apps/web/lib/api.ts`, or
   `packages/schemas`. Forces "enumerate every consumer before editing the
   contract" (per `thinking/README.md` §5). Ships with a **new source doc**
   `docs/recipes/contract-first-change.md` so the wrap-a-doc invariant holds.
   *Proves: authoring a net-new skill + its backing recipe.*

3. **`.claude/agents/repo-context-reader.md`** — a subagent that runs the
   `CONTEXT-LOADING` protocol and returns a compact **context pack** (relevant
   files, applicable ADRs, rules in scope, gated-area flags), keeping heavy
   file-reading out of the main thread. Wraps `docs/agents/repo-context-reader.md`.
   *Proves: an isolated-context subagent; directly serves Context Engineering.*

4. **`.claude/commands/feature.md`** (`/feature`) — wraps
   `docs/workflows/feature.md`; orchestrates classify → load context → brainstorm
   → (spec/plan if large) → branch → implement → test → verify → docs → PR,
   honoring every 🚦 gate. *Proves: a slash command over a workflow.*

5. **`.claude/CONVENTIONS.md` + `.claude/README.md`** — the meta-spec governing
   all future artifacts: frontmatter rules, trigger-phrasing style, the mandatory
   `make verify` + DoD gate, the protected-path guard, and the doc→artifact
   mapping table above. This is what makes fan-out cheap and consistent.

## Conventions every artifact must follow

- **Frontmatter:** `name` (kebab-case); `description` begins `"Use when …"` so the
  artifact self-triggers (Superpowers style).
- **Guard rails inlined:** each skill/command restates only the AI-CONTRACT gates
  it touches (never push `main`; protected paths → human approval;
  honesty-before-done) and links to the canonical source — it never re-authors
  policy.
- **Verification gate:** every artifact ends by pointing at `make verify` +
  `docs/verification/definition-of-done.md`.
- **No new rules:** artifacts link to `CLAUDE.md` / `AI-CONTRACT.md`; they never
  restate policy as if authoritative. If an artifact would need a rule that does
  not exist, that is a signal to add it to the owning source first.

## Verification for this build

`.claude/` artifacts are not unit-testable, so "done" means:

1. **Frontmatter lint** — valid YAML, required keys present on every artifact.
2. **Link check** — every `@docs/...` reference resolves (the repo already runs a
   link check: 394 links, 0 broken — this must stay 0).
3. **Live trigger test** — for each skill/command, describe a matching task and
   confirm Claude auto-invokes the artifact; for the subagent, confirm it returns
   a context pack without polluting the main thread.
4. **No policy drift** — diff shows no changes to `CLAUDE.md`, `AI-CONTRACT.md`, or
   any `docs/` file except the one net-new recipe (`contract-first-change.md`).

## After the spine (future batches, out of scope here)

Once conventions are proven: fan-out remaining `recipes → skills`,
`agents → agents`, `workflows → commands`; then author the other net-new skills
(`add-error-code`, `backend-service-change`). Each future batch is a mechanical
repeat of the proven pattern, one spec per batch.

## Open questions

None. Spine size (5), spine composition (context-reader in, code-reviewer out),
and the thin-wrapper architecture were confirmed during brainstorming.

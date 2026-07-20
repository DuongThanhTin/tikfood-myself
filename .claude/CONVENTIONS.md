# .claude/ Authoring Conventions

> **Why this file exists:** `.claude/` artifacts are the *executable* layer over the
> repo's documentation-based AI Operating System. Each artifact is a **thin trigger**,
> not a copy. This file is the contract every skill / agent / command must follow so the
> layer stays consistent and never drifts from `docs/`. Design: see
> `docs/superpowers/specs/2026-07-08-claude-executable-layer-spine-design.md`.

## The one rule

Artifacts **wrap docs; they do not replace them.** The authoritative procedure lives in
`docs/`. An artifact contains: frontmatter (to self-trigger) + a short checklist + an
`@docs/...` pointer to the real steps. Never duplicate a procedure. Never re-author a
policy — link to `../CLAUDE.md` / `docs/ai/AI-CONTRACT.md`.

## Doc -> artifact mapping

| Artifact | Wraps | Type |
| --- | --- | --- |
| `.claude/skills/<name>/SKILL.md` | `docs/recipes/*` | auto-triggered procedure |
| `.claude/agents/<name>.md` | `docs/agents/*` | isolated-context subagent |
| `.claude/commands/<name>.md` | `docs/workflows/*` | human-invoked orchestration |
| *(stays reference)* | `docs/ai/*`, `docs/thinking`, `docs/verification/*` | always-loaded, not triggered |

## Frontmatter

- **Skills** (`.claude/skills/<name>/SKILL.md`): `name` (kebab-case) + `description`
  beginning `Use when …` (this is what makes it self-trigger).
- **Agents** (`.claude/agents/<name>.md`): `name` + `description` (`Use when …`);
  optional `tools`, `model`.
- **Commands** (`.claude/commands/<name>.md`): `description`; optional `argument-hint`.
  The slash name is the filename (`feature.md` -> `/feature`).

## @-ref vs bare citation

Use `@path` for the authoritative pointer (the doc the artifact auto-loads or the
reader must follow first) and for any path that `check-artifacts.sh` must verify.
Plain citations (narrative mentions that are not load-bearing) may be bare.
Concretely: `@docs/verification/definition-of-done.md`, `@docs/ai/AI-CONTRACT.md`,
and `@CLAUDE.md` in body prose must carry the `@` prefix.

## Body shape (every skill and command)

1. One-line purpose.
2. `Authoritative steps: @docs/...` — the pointer, read-first.
3. A short quick-checklist (the gist, not the full procedure).
4. Guard rails it touches, each linking its owning source (never push `main`; protected
   paths -> human approval; honesty-before-done).
5. Close with the gate: `make verify` + `@docs/verification/definition-of-done.md`.

Agents are exempt from step 5 — a read-only context subagent has nothing to verify.

## Verification

Run `bash .claude/check-artifacts.sh` (or `make verify-claude`): valid frontmatter +
every referenced repo path resolves. Then a **live trigger test** — describe a matching
task and confirm the artifact self-invokes.

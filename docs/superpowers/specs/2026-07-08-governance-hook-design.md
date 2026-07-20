# Design — Governance PreToolUse Hook

**Date:** 2026-07-08 · **Branch:** `ai/claude-executable-layer` · **Status:** approved for planning

## Problem

The repo's most safety-critical rules — never push `main`/`master`, never force-push,
never read secrets, protected paths require human approval — live only as **prose** in
`CLAUDE.md` / `docs/ai/AI-CONTRACT.md`. In an interactive Claude Code session the model is
*asked* to obey them; nothing *enforces* them. The automation runner already enforces the
same rules deterministically (`apps/ai-code-runner/src/guards/`: `commandGuard.ts`,
`secretGuard.ts`, `fileGuard.ts`, driven by `.ai-agent.yaml`). The interactive side has no
equivalent. This closes V2 vision gap #3 (governance = enforced, not advised) and brings
the interactive side to parity with the runner.

## Non-goals

- **Not a replacement for the prose.** The rules docs stay (they explain *why*); the hook
  makes the critical subset *true*.
- **Not prompt-injection defense.** That is judgment-heavy and out of scope here.
- **Not a full policy engine.** Only the deterministic, pattern-matchable subset of the
  existing `.ai-agent.yaml` lists.
- **No new rules.** The hook enforces exactly the lists already in `.ai-agent.yaml`; it
  invents nothing.

## Architecture

A `PreToolUse` hook registered in project `.claude/settings.json` runs a guard on every
tool call. The guard reads the tool-call JSON from stdin, reads the canonical lists from
**`.ai-agent.yaml` at runtime** (single source of truth — the interactive guard cannot
drift from the runner's rules), and returns allow/deny.

The guard is **python3** (not bash): the hook input is JSON, and `jq` is not guaranteed to
be installed, whereas python3 is present and gives clean JSON parsing + path/command
matching.

### Two-tier enforcement

The existing rules are not equal, so enforcement is two-tier:

| Tier | Source list / patterns | On match | Override |
|------|------------------------|----------|----------|
| **Never** | `.ai-agent.yaml` `blocked_commands` (push `main`/`master`, `git push --force`, `git reset --hard`, `git clean -fdx`, `rm -rf /`, `curl \| sh`, `wget \| sh`, `docker system prune`) **and** secret reads (`.env`, `.env.*`, `secrets/**`, `credentials/**`, `**/*.pem`, `**/*.key`, `**/*id_rsa*`, `**/*private_key*`) | **deny** | none — never legitimate |
| **Gated** | `.ai-agent.yaml` `protected_paths` writes (`packages/schemas/**`, `.github/workflows/**`, `docker-compose.yml`, `apps/ai-code-runner/Dockerfile`, and the secret paths above) | **deny** unless override set | `TIKFOOD_ALLOW_PROTECTED=1` in the environment |

Rationale: `blocked_commands` mean *never*; `protected_paths` mean *require human
approval*. The override exists only for the gated tier — so a human who has genuinely
approved a schema/workflow change can proceed, while catastrophic commands and secret
reads stay blocked regardless.

### What the guard inspects per tool

- **`Bash`**: the `command` string → Never-tier `blocked_commands` patterns; secret-read
  patterns (`cat`/`less`/`head`/`tail` of a secret path); writes/moves into protected or
  secret paths (redirection `>`, `tee`, `mv`, `cp` targets).
- **`Write` / `Edit` / `NotebookEdit`**: the `file_path` → protected-path (gated) and
  secret-path (never) checks.
- **`Read`**: the `file_path` → secret-path (never) check.
- Any other tool → allow (out of scope).

### Block signal

The guard denies via the Claude Code PreToolUse block convention: **exit code 2** with a
human-readable reason on stderr (the reason is surfaced to the model so it knows *why* and
what to do — e.g. "protected path; get human approval then set `TIKFOOD_ALLOW_PROTECTED=1`").
Exit 0 allows. (If the running harness expects the JSON `permissionDecision` form instead,
the guard is adapted during implementation and confirmed by the live check below — the
contract that matters is *the tool call is actually blocked*.)

## Components

- `.claude/hooks/governance-guard.py` — the guard.
- `.claude/hooks/governance-guard.test.sh` — unit tests: mock PreToolUse JSON on stdin,
  assert allow/deny. The RED/GREEN proof (mirrors `commandGuard.test.ts`).
- `.claude/settings.json` — registers the `PreToolUse` hook to run the guard.
- `.claude/hooks/README.md` — enforced rules, two-tier model, the override, honest limits.
- `Makefile` — `make verify-hooks` runs the guard tests; added to the `verify` aggregate
  (deterministic, so unlike the probabilistic `/eval-skills` it belongs in `make verify`).

## Verification

1. **Unit (deterministic):** the guard test suite proves each tier — a `blocked_command`
   denies; a secret read denies; a protected-path write denies without the override and
   **allows with** `TIKFOOD_ALLOW_PROTECTED=1`; a benign call (e.g. `Read apps/api/main.go`,
   `Write docs/foo.md`, `git status`) passes. Every deny is proven RED before GREEN.
2. **Live (controller):** after wiring `settings.json`, the controller attempts a blocked
   op in-session and confirms the hook actually denies it. A hook that silently fails to
   block is governance theatre — the live check is mandatory. If hooks require a session
   reload to take effect, that is a documented manual step.

## Conscious choices (flagged, approved during brainstorming)

- **Placement = project `.claude/settings.json` (checked in),** not `settings.local.json`.
  The hook is an AI-OS artifact and should travel with the repo (H2 portability). Cost: it
  activates for anyone using the repo; a guard bug could block legit work until fixed —
  mitigated by the override and the unit tests.
- **Guard reads `.ai-agent.yaml` at runtime,** not hardcoded patterns — DRY, drift-proof.

## Open questions

None. Enforcement mode (hard-block + override), override signal (`TIKFOOD_ALLOW_PROTECTED`),
placement (project settings.json), and single-source-of-truth (`.ai-agent.yaml`) were
confirmed during brainstorming.

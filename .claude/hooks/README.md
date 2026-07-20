# .claude/hooks — governance guard

A `PreToolUse` hook (`governance-guard.py`, registered in `.claude/settings.json`) that
makes the repo's safety rules *enforced*, not just documented. It mirrors the automation
runner's guards (`apps/ai-code-runner/src/guards/`) on the interactive side, reading the
same lists from `.ai-agent.yaml` at runtime so the two can never drift.

## What it blocks

**Never (no override):**
- `blocked_commands` from `.ai-agent.yaml` — push to `main`/`master`, `git push --force`
  (incl. `-f` / `--force-with-lease`), `git reset --hard`, `git clean -fdx`, `rm -rf /`,
  `curl | sh`, `wget | sh`, `docker system prune`.
- Reading a secret file (`.env`, `.env.*`, `secrets/**`, `credentials/**`, `*.pem`,
  `*.key`, `*id_rsa*`, `*private_key*`) — via the `Read` tool or a `cat`-like shell read.

**Gated (override with `TIKFOOD_ALLOW_PROTECTED=1`):**
- Writing a `protected_paths` target (`packages/schemas/**`, `.github/workflows/**`,
  `docker-compose.yml`, `apps/ai-code-runner/Dockerfile`) — because the rule is *requires
  human approval*, not *never*. When you have approval, run with the env var set.

## Honest limitations
- **Fail-open on the guard's own parse error** — a guard bug never bricks the session; the
  deny paths are covered by `governance-guard.test.sh` (`make verify-hooks`).
- **Bash matching is pragmatic** — it catches the listed commands, `cat`-style secret reads,
  and redirection/`tee` writes into protected paths; it is not a full shell parser.
- **Enforces only the deterministic subset** of the rules. Judgment rules (prompt injection,
  product scope) stay in `CLAUDE.md` / `docs/ai/AI-CONTRACT.md`.
- **Case-sensitive matching** — an uppercase secret name at repo root (e.g. `SERVER.PEM`) may not match; keep secret files lowercase or rely on `secrets/**` placement.
- **No two-step fetch detection** — `curl -o f.sh <url>; sh f.sh` (fetch then run separately) is not caught; only direct `curl … | sh` pipes are.

## Verify / disable
- Test: `make verify-hooks`.
- To do approved protected-path work: `TIKFOOD_ALLOW_PROTECTED=1` in the environment.
- To disable entirely: remove the `PreToolUse` block from `.claude/settings.json`.

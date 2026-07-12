# .claude/hooks — PreToolUse guards

Two `PreToolUse` hooks (registered in `.claude/settings.json`) that make the repo's rules
*enforced*, not just documented:

1. **`governance-guard.py`** — safety: blocked commands, secret reads, protected-path writes.
2. **`verify-evidence-guard.py`** — honesty: no commit / PR without fresh proof that tests pass.

## 1. governance-guard.py

Mirrors the automation runner's guards (`apps/ai-code-runner/src/guards/`) on the
interactive side, reading the same lists from `.ai-agent.yaml` at runtime so the two can
never drift.

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

## 2. verify-evidence-guard.py

Makes "I verified" non-fakeable. Each `make verify-<surface>` records a pass under
`.verify-evidence/<surface>` (via `record-verify.sh`, gitignored). Before a `git commit`
or `gh pr create`, the hook maps every changed file to the surface(s) that must be fresh:

| Changed path | Requires fresh evidence for |
| --- | --- |
| `apps/api/**` | `api` |
| `apps/web/**` | `web` |
| `apps/ai-code-runner/**` | `runner` |
| `packages/**` (shared contracts) | `api` + `web` + `runner` |
| `.claude/hooks/**`, `.claude/*.sh` | `hooks` |
| anything else (docs, Makefile, root) | — (not gated) |

"Fresh" = the surface's evidence file is **newer than the changed source files** — so
editing code after verifying makes the evidence stale and re-blocks you. Denies with the
exact `make verify-<surface>` command to run.

### Honest limitations (this hook)
- **Fail-open** on parse error or when git is unavailable (never bricks the session).
- **Freshness is by mtime** of the working tree — it proves verify ran after your last edit,
  not that the diff is bug-free. Pasted evidence in the PR is still the human-readable proof.
- **Not per-file**: any change under a surface requires that whole surface's `make verify`.

## Verify / disable
- Test: `make verify-hooks` (runs both guards' unit tests).
- To do approved protected-path work: `TIKFOOD_ALLOW_PROTECTED=1` in the environment.
- Emergency commit without fresh evidence: `TIKFOOD_SKIP_VERIFY_GATE=1` (say so in the PR).
- To disable entirely: remove the relevant `PreToolUse` block from `.claude/settings.json`.

# Governance PreToolUse Hook Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Enforce the repo's `.ai-agent.yaml` `blocked_commands` + `protected_paths` in interactive Claude Code via a `PreToolUse` guard, so the safety rules are mechanically true, not just prose.

**Architecture:** A python3 guard reads each tool call's JSON from stdin, reads the canonical lists from `.ai-agent.yaml` at runtime (single source of truth), and denies (exit 2 + stderr reason) on a Never-tier match (blocked_commands, secret reads) or a Gated-tier match (protected-path writes) unless `TIKFOOD_ALLOW_PROTECTED=1`. It is registered as a `PreToolUse` hook in project `.claude/settings.json`.

**Tech Stack:** python3 (JSON parsing + glob/command matching); Bash (unit test harness); GNU Make; Claude Code hooks (settings.json).

## Global Constraints

- **Single source of truth = `.ai-agent.yaml`.** The guard reads `protected_paths` and `blocked_commands` at runtime; it hardcodes no rule list.
- **Two tiers.** Never: `blocked_commands` + secret-file **reads** → deny, no override. Gated: `protected_paths` **writes** → deny unless `TIKFOOD_ALLOW_PROTECTED=1`.
- **Block signal:** deny = write reason to **stderr** and **exit 2**; allow = **exit 0**. On the guard's own parse error → **exit 0 (fail-open)** so a guard bug never bricks the session (conscious tradeoff — the guard is a solo-dev safety net, not a hostile-input boundary; a fail-closed guard that bricks every tool call would just get disabled).
- **Placement:** project `.claude/settings.json` (checked in).
- **No new rules;** enforce exactly the `.ai-agent.yaml` lists (force-push / push-main are matched robustly as the *same* rule, not new ones).
- **Deterministic → in `make verify`.** `make verify-hooks` runs the guard unit tests and joins the `verify` aggregate.
- **Branch:** `ai/claude-executable-layer`. Never push `main`, never force-push, never auto-merge.

---

### Task 1: The guard + unit tests + Make wiring

The guard and its test suite are one deliverable: the guard is only trustworthy with the tests that prove each deny/allow branch. TDD — the test script is written to fail first.

**Files:**
- Create: `.claude/hooks/governance-guard.py`
- Create: `.claude/hooks/governance-guard.test.sh`
- Modify: `Makefile` (add `verify-hooks` to `.PHONY` + the `verify` aggregate + the target)

**Interfaces:**
- Consumes: `.ai-agent.yaml` (`protected_paths`, `blocked_commands`); env `CLAUDE_PROJECT_DIR`, `TIKFOOD_ALLOW_PROTECTED`.
- Produces: `.claude/hooks/governance-guard.py` — a CLI reading PreToolUse JSON on stdin; exit 2 = deny (reason on stderr), exit 0 = allow.

- [ ] **Step 1: Write the guard**

Create `.claude/hooks/governance-guard.py`:

```python
#!/usr/bin/env python3
"""Governance PreToolUse guard for interactive Claude Code.

Enforces the SAME lists the runner uses (.ai-agent.yaml protected_paths +
blocked_commands) so the interactive side matches the automation side.

Tiers:
  NEVER  — blocked_commands (Bash) and secret-file READS: denied, no override.
  GATED  — writes to protected_paths: denied unless TIKFOOD_ALLOW_PROTECTED=1.

Contract: read the PreToolUse JSON on stdin. Deny = reason to stderr + exit 2
(Claude Code blocks the call and shows the reason to the model). Allow = exit 0.
On our own parse error -> exit 0 (fail-open) so a guard bug never bricks the
session; the deny paths are covered by governance-guard.test.sh.
"""
import json, os, re, sys

ROOT = os.environ.get("CLAUDE_PROJECT_DIR") or os.getcwd()
AI_AGENT = os.path.join(ROOT, ".ai-agent.yaml")
SECRET_HINTS = (".env", "secret", "credential", "id_rsa", "private_key", ".pem", ".key")


def _read_list(text, key):
    """Extract a simple YAML list: '  - value' lines under 'key:' until the next
    top-level key. Strips surrounding quotes."""
    out, in_block = [], False
    for line in text.splitlines():
        if re.match(r"^%s\s*:" % re.escape(key), line):
            in_block = True
            continue
        if in_block:
            m = re.match(r"^\s+-\s+(.*\S)\s*$", line)
            if m:
                out.append(m.group(1).strip().strip('"').strip("'"))
            elif re.match(r"^\S", line):  # next top-level key ends the block
                break
    return out


def _load_lists():
    try:
        with open(AI_AGENT, encoding="utf-8") as f:
            text = f.read()
    except OSError:
        return [], []
    return _read_list(text, "protected_paths"), _read_list(text, "blocked_commands")


def _glob_to_re(pattern):
    """Translate a protected-path glob to a regex fragment: ** -> any, * -> non-slash."""
    p = pattern.strip().lstrip("./")
    out, i = "", 0
    while i < len(p):
        if p[i:i + 2] == "**":
            out += ".*"
            i += 2
        elif p[i] == "*":
            out += "[^/]*"
            i += 1
        else:
            out += re.escape(p[i])
            i += 1
    return out


def path_matches(path, pattern):
    path = path.strip().lstrip("./")
    frag = _glob_to_re(pattern)
    base = os.path.basename(path)
    return bool(
        re.fullmatch(frag, path)
        or re.fullmatch(frag + r"(/.*)?", path)
        or re.search(r"(^|/)" + frag + r"(/|$)", path)
        or re.fullmatch(frag, base)  # basename catch (e.g. **/*.pem at repo root)
    )


def is_secret_pattern(pattern):
    pl = pattern.lower()
    return any(h in pl for h in SECRET_HINTS)


def deny(reason):
    sys.stderr.write("BLOCKED by governance-guard: " + reason + "\n")
    sys.exit(2)


def main():
    try:
        data = json.loads(sys.stdin.read())
    except (ValueError, TypeError):
        sys.exit(0)  # fail-open on our own parse bug
    tool = data.get("tool_name") or data.get("tool") or ""
    ti = data.get("tool_input") or data.get("input") or {}
    if not isinstance(ti, dict):
        sys.exit(0)
    protected, blocked = _load_lists()
    override = os.environ.get("TIKFOOD_ALLOW_PROTECTED") == "1"
    secret_patterns = [p for p in protected if is_secret_pattern(p)]

    if tool == "Bash":
        cmd = " ".join((ti.get("command") or "").split())
        low = cmd.lower()
        for b in blocked:                                   # NEVER: blocked_commands
            if " ".join(b.lower().split()) in low:
                deny("command matches blocked_commands: %r" % b)
        if re.search(r"\bgit\s+push\b", low) and re.search(r"(--force\b|(?<!\w)-f\b|--force-with-lease\b)", low):
            deny("force-push is forbidden")
        if re.search(r"\bgit\s+push\b", low) and re.search(r"\b(main|master)\b", low):
            deny("pushing to main/master is forbidden")
        if re.search(r"\b(cat|less|more|head|tail|xxd|od|strings|base64|openssl)\b", low):  # NEVER: secret read
            for tok in cmd.split():
                for sp in secret_patterns:
                    if path_matches(tok, sp):
                        deny("reading a secret path is forbidden: %s" % tok)
        write_targets = re.findall(r">>?\s*(\S+)", cmd) + re.findall(r"\btee\s+(?:-a\s+)?(\S+)", cmd)
        for tgt in write_targets:                           # GATED: protected write via shell
            for pp in protected:
                if path_matches(tgt, pp) and not override:
                    deny("writing a protected path needs approval "
                         "(set TIKFOOD_ALLOW_PROTECTED=1): %s" % tgt)
        sys.exit(0)

    if tool in ("Write", "Edit", "NotebookEdit"):
        path = ti.get("file_path") or ti.get("notebook_path") or ""
        for pp in protected:                                # GATED: protected write
            if path_matches(path, pp):
                if override:
                    sys.exit(0)
                deny("writing a protected path needs approval "
                     "(set TIKFOOD_ALLOW_PROTECTED=1): %s" % path)
        sys.exit(0)

    if tool == "Read":
        path = ti.get("file_path") or ""
        for sp in secret_patterns:                          # NEVER: secret read
            if path_matches(path, sp):
                deny("reading a secret path is forbidden: %s" % path)
        sys.exit(0)

    sys.exit(0)  # any other tool: out of scope


if __name__ == "__main__":
    main()
```

- [ ] **Step 2: Write the failing test harness**

Create `.claude/hooks/governance-guard.test.sh`:

```bash
#!/usr/bin/env bash
# Unit tests for governance-guard.py. Feeds mock PreToolUse JSON on stdin and
# asserts the exit code (2 = deny, 0 = allow).
set -uo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
GUARD="$ROOT/.claude/hooks/governance-guard.py"
export CLAUDE_PROJECT_DIR="$ROOT"

pass=0; fail=0

# assert <expected-exit> <label> <json> [override]
assert() {
  local want="$1" label="$2" json="$3" ov="${4:-}"
  local got
  TIKFOOD_ALLOW_PROTECTED="$ov" bash -c 'python3 "$1" <<<"$2"' _ "$GUARD" "$json" >/dev/null 2>&1
  got=$?
  if [ "$got" = "$want" ]; then
    pass=$((pass+1))
  else
    echo "FAIL: $label — want exit $want, got $got"
    fail=$((fail+1))
  fi
}

# NEVER — blocked_commands
assert 2 "push to main"        '{"tool_name":"Bash","tool_input":{"command":"git push origin main"}}'
assert 2 "force push"          '{"tool_name":"Bash","tool_input":{"command":"git push origin ai/x --force"}}'
assert 2 "force push -f"       '{"tool_name":"Bash","tool_input":{"command":"git push -f origin ai/x"}}'
assert 2 "git reset --hard"    '{"tool_name":"Bash","tool_input":{"command":"git reset --hard HEAD~1"}}'
assert 2 "rm -rf /"            '{"tool_name":"Bash","tool_input":{"command":"rm -rf /"}}'
# NEVER — secret reads (Read tool + Bash cat), even with override
assert 2 "read .env"           '{"tool_name":"Read","tool_input":{"file_path":"apps/web/.env"}}'
assert 2 "read pem"            '{"tool_name":"Read","tool_input":{"file_path":"secrets/tls/server.pem"}}'
assert 2 "cat secret bash"     '{"tool_name":"Bash","tool_input":{"command":"cat .env"}}'
assert 2 "read .env w/override" '{"tool_name":"Read","tool_input":{"file_path":".env"}}' 1
# GATED — protected writes deny WITHOUT override
assert 2 "write schema"        '{"tool_name":"Write","tool_input":{"file_path":"packages/schemas/foo.json"}}'
assert 2 "edit workflow"       '{"tool_name":"Edit","tool_input":{"file_path":".github/workflows/ci.yml"}}'
assert 2 "tee into schema"     '{"tool_name":"Bash","tool_input":{"command":"echo x | tee packages/schemas/x.json"}}'
# GATED — protected writes ALLOW WITH override
assert 0 "write schema +ov"    '{"tool_name":"Write","tool_input":{"file_path":"packages/schemas/foo.json"}}' 1
assert 0 "edit workflow +ov"   '{"tool_name":"Edit","tool_input":{"file_path":".github/workflows/ci.yml"}}' 1
# ALLOW — benign
assert 0 "read normal go"      '{"tool_name":"Read","tool_input":{"file_path":"apps/api/main.go"}}'
assert 0 "write normal doc"    '{"tool_name":"Write","tool_input":{"file_path":"docs/foo.md"}}'
assert 0 "git status"          '{"tool_name":"Bash","tool_input":{"command":"git status --short"}}'
assert 0 "edit app code"       '{"tool_name":"Edit","tool_input":{"file_path":"apps/web/lib/api.ts"}}'
# ALLOW — malformed input fails open
assert 0 "malformed json"      'not json at all'

echo ""
if [ "$fail" -eq 0 ]; then
  echo "PASS: $pass guard assertions"
else
  echo "FAIL: $fail failing / $pass passing"
fi
exit $((fail > 0))
```

- [ ] **Step 3: Run the tests — expect RED**

Run:
```bash
chmod +x .claude/hooks/governance-guard.py .claude/hooks/governance-guard.test.sh
bash .claude/hooks/governance-guard.test.sh; echo "exit: $?"
```
Expected before the guard is correct: failures printed and `exit: 1`. (If Step 1's guard is already correct, you will instead see GREEN — that is acceptable; the point is the test must be able to fail. To confirm the tests are real, temporarily rename `deny` to a no-op and re-run: the NEVER/GATED asserts must then FAIL. Restore.)

- [ ] **Step 4: Run the tests — expect GREEN**

Run:
```bash
bash .claude/hooks/governance-guard.test.sh; echo "exit: $?"
```
Expected: `PASS: 19 guard assertions` and `exit: 0`.

- [ ] **Step 5: Wire `make verify-hooks`**

In `Makefile`: add `verify-hooks` to the `.PHONY` line; append `verify-hooks` to the `verify` aggregate target's dependency list (after `verify-ai-os`); add the target at end of file:
```makefile
## verify-hooks: governance guard unit tests (blocked_commands + protected_paths enforcement)
verify-hooks:
	bash .claude/hooks/governance-guard.test.sh
```

- [ ] **Step 6: Run it**

Run:
```bash
make verify-hooks
```
Expected: `PASS: 19 guard assertions`.

- [ ] **Step 7: Commit**

```bash
git add .claude/hooks/governance-guard.py .claude/hooks/governance-guard.test.sh Makefile
git commit -m "feat(claude): governance guard (enforces .ai-agent.yaml protected_paths + blocked_commands)

python3 PreToolUse guard: NEVER-tier (blocked_commands, secret reads) + GATED-tier
(protected-path writes, override via TIKFOOD_ALLOW_PROTECTED=1). 19 unit assertions;
wired into make verify-hooks + make verify.

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

### Task 2: Register the hook + README + controller live-check

**Files:**
- Create: `.claude/settings.json`
- Create: `.claude/hooks/README.md`

**Interfaces:**
- Consumes: `.claude/hooks/governance-guard.py` (Task 1).
- Produces: an active `PreToolUse` hook.

- [ ] **Step 1: Create `.claude/settings.json`**

Create `.claude/settings.json`:

```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash|Write|Edit|NotebookEdit|Read",
        "hooks": [
          {
            "type": "command",
            "command": "python3 \"$CLAUDE_PROJECT_DIR/.claude/hooks/governance-guard.py\""
          }
        ]
      }
    ]
  }
}
```

- [ ] **Step 2: Create `.claude/hooks/README.md`**

Create `.claude/hooks/README.md`:

```markdown
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

## Verify / disable
- Test: `make verify-hooks`.
- To do approved protected-path work: `TIKFOOD_ALLOW_PROTECTED=1` in the environment.
- To disable entirely: remove the `PreToolUse` block from `.claude/settings.json`.
```

- [ ] **Step 3: Verify the syntactic gate still passes**

Run:
```bash
make verify-claude
```
Expected: `PASS: 5 artifact(s), ...`. The new files are under `.claude/hooks/` and `.claude/settings.json`, none under `skills/agents/commands`, so the artifact count is unchanged (still 5).

- [ ] **Step 4: Controller live-check (NOT delegated)**

This step is run by the controller in the interactive session, because it verifies the hook actually fires. A hook that silently fails to block is governance theatre.
1. Confirm the hook is loaded (if the harness only loads hooks at session start, note that a reload may be required — record honestly whether it took effect).
2. Attempt a **gated** denied op — e.g. ask to `Write` a throwaway file under `packages/schemas/` — and confirm the tool call is **blocked** with the guard's reason.
3. Confirm the **override** path: with `TIKFOOD_ALLOW_PROTECTED=1`, the same write is allowed. (Immediately revert any test file.)
4. Confirm a benign op (e.g. `Read` an app source file) is **not** blocked.
Record the outcomes. If the block signal (exit 2) does not actually deny in this harness, adapt the guard to emit the JSON `permissionDecision:"deny"` form and re-check — the contract that matters is that the call is blocked.

- [ ] **Step 5: Commit**

```bash
git add .claude/settings.json .claude/hooks/README.md
git commit -m "feat(claude): register governance PreToolUse hook + hooks README

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

## Self-Review

**Spec coverage:**
- python3 guard reading `.ai-agent.yaml` at runtime → Task 1 Step 1. ✅
- Two tiers (Never: blocked_commands + secret reads; Gated: protected writes + override) → guard `main()` + tests. ✅
- Block signal exit 2 / allow exit 0 / fail-open on parse error → guard `deny()` + `main()` + a test asserting malformed→0. ✅
- Override = `TIKFOOD_ALLOW_PROTECTED=1` → guard + tests (`+ov` cases). ✅
- Guard inspects Bash / Write·Edit·NotebookEdit / Read → guard branches + tests per tool. ✅
- `make verify-hooks` deterministic + in `verify` aggregate → Task 1 Step 5. ✅
- Placement project `.claude/settings.json` (checked in) → Task 2 Step 1. ✅
- README states rules/override/limitations → Task 2 Step 2. ✅
- Live verification (controller, mandatory) → Task 2 Step 4. ✅
- Single source of truth / no new rules → guard reads yaml; force-push/push-main matched as same rule. ✅

**Placeholder scan:** guard, test harness, settings.json, README all shown in full; commands have expected output; no TBD/TODO. ✅

**Name/count consistency:** `TIKFOOD_ALLOW_PROTECTED`, `governance-guard.py`, `verify-hooks`, exit-2-deny, and `PASS: 19 guard assertions` are consistent across guard, tests, README, Make, and tasks. Assertion count (19) matches the test harness. `make verify-claude` count stays `5` (no skills/agents/commands added). ✅

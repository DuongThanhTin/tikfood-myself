#!/usr/bin/env bash
# Unit tests for verify-evidence-guard.py. Each case builds a throwaway git repo, plants
# changed files + evidence with controlled mtimes (touch -t), feeds mock PreToolUse JSON,
# and asserts the exit code (2 = deny, 0 = allow).
set -uo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
GUARD="$ROOT/.claude/hooks/verify-evidence-guard.py"

TMPS=()
cleanup() { for d in "${TMPS[@]:-}"; do [ -n "$d" ] && rm -rf "$d"; done; }
trap cleanup EXIT

pass=0; fail=0
COMMIT='{"tool_name":"Bash","tool_input":{"command":"git commit -m x"}}'
ADD_COMMIT='{"tool_name":"Bash","tool_input":{"command":"git add -A && git commit -m x"}}'
PR='{"tool_name":"Bash","tool_input":{"command":"gh pr create --base dev"}}'
STATUS='{"tool_name":"Bash","tool_input":{"command":"git status --short"}}'

OLD=202601010000   # touch -t: 2026-01-01 00:00
NEW=202612310000   # touch -t: 2026-12-31 00:00

newrepo() { local d; d="$(mktemp -d)"; TMPS+=("$d"); git init -q "$d"; echo "$d"; }
mkfile()  { mkdir -p "$(dirname "$1")"; echo x > "$1"; }             # mkfile <path>
evid()    { mkdir -p "$(dirname "$1")"; echo "ts	head	PASS" > "$1"; } # evid <path>

# check <want-exit> <label> <repo> <json> [env=val ...]
check() {
  local want="$1" label="$2" repo="$3" json="$4"; shift 4
  local got
  env "$@" CLAUDE_PROJECT_DIR="$repo" python3 "$GUARD" <<<"$json" >/dev/null 2>&1
  got=$?
  if [ "$got" = "$want" ]; then pass=$((pass+1)); else
    echo "FAIL: $label — want exit $want, got $got"; fail=$((fail+1)); fi
}

# T1 — docs-only change: nothing verify covers -> allow
d="$(newrepo)"; mkfile "$d/docs/x.md"
check 0 "docs-only commit" "$d" "$COMMIT"

# T2 — api change, no evidence -> deny
d="$(newrepo)"; mkfile "$d/apps/api/x.go"
check 2 "api change, no evidence" "$d" "$COMMIT"

# T3 — api change, FRESH api evidence -> allow
d="$(newrepo)"; mkfile "$d/apps/api/x.go"; touch -t $OLD "$d/apps/api/x.go"
evid "$d/.verify-evidence/api"; touch -t $NEW "$d/.verify-evidence/api"
check 0 "api change, fresh evidence" "$d" "$COMMIT"

# T4 — api change, STALE api evidence (code edited after verify) -> deny
d="$(newrepo)"; mkfile "$d/apps/api/x.go"; touch -t $NEW "$d/apps/api/x.go"
evid "$d/.verify-evidence/api"; touch -t $OLD "$d/.verify-evidence/api"
check 2 "api change, stale evidence" "$d" "$COMMIT"

# T5 — stale evidence but emergency override -> allow
d="$(newrepo)"; mkfile "$d/apps/api/x.go"; touch -t $NEW "$d/apps/api/x.go"
evid "$d/.verify-evidence/api"; touch -t $OLD "$d/.verify-evidence/api"
check 0 "stale + override" "$d" "$COMMIT" TIKFOOD_SKIP_VERIFY_GATE=1

# T6 — non-commit command (git status) with gated change, no evidence -> allow
d="$(newrepo)"; mkfile "$d/apps/api/x.go"
check 0 "status not gated" "$d" "$STATUS"

# T7 — gh pr create, web change, missing evidence -> deny
d="$(newrepo)"; mkfile "$d/apps/web/y.ts"
check 2 "pr create, web missing" "$d" "$PR"

# T8 — packages change requires api+web+runner; only api fresh -> deny
d="$(newrepo)"; mkfile "$d/packages/schemas/z.json"; touch -t $OLD "$d/packages/schemas/z.json"
evid "$d/.verify-evidence/api"; touch -t $NEW "$d/.verify-evidence/api"
check 2 "packages, partial evidence" "$d" "$COMMIT"

# T9 — packages change, all three surfaces fresh -> allow
d="$(newrepo)"; mkfile "$d/packages/schemas/z.json"; touch -t $OLD "$d/packages/schemas/z.json"
for s in api web runner; do evid "$d/.verify-evidence/$s"; touch -t $NEW "$d/.verify-evidence/$s"; done
check 0 "packages, all fresh" "$d" "$COMMIT"

# T10 — `git add -A && git commit` (nothing staged yet), untracked api file, no evidence -> deny
d="$(newrepo)"; mkfile "$d/apps/api/x.go"
check 2 "add&&commit catches untracked" "$d" "$ADD_COMMIT"

# T11 — malformed JSON -> fail open (allow)
d="$(newrepo)"; mkfile "$d/apps/api/x.go"
check 0 "malformed json" "$d" 'not json at all'

# T12 — not a git repo -> fail open (allow)
d="$(mktemp -d)"; TMPS+=("$d"); mkfile "$d/apps/api/x.go"
check 0 "not a git repo" "$d" "$COMMIT"

# T13 — guard source (.claude/hooks) change maps to 'hooks'; no evidence -> deny
d="$(newrepo)"; mkfile "$d/.claude/hooks/foo.py"
check 2 "hooks change, no evidence" "$d" "$COMMIT"

echo ""
if [ "$fail" -eq 0 ]; then
  echo "PASS: $pass verify-evidence assertions"
else
  echo "FAIL: $fail failing / $pass passing"
fi
exit $((fail > 0))

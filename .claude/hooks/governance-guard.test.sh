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

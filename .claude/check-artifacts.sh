#!/usr/bin/env bash
# Verify every .claude/ artifact: valid frontmatter + all referenced repo paths resolve.
# Usage: bash .claude/check-artifacts.sh
set -uo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

fail=0
count=0

check_frontmatter() {
  local file="$1"
  # First line must be '---' opening a YAML frontmatter block.
  if [ "$(head -n1 "$file")" != "---" ]; then
    echo "FAIL [$file]: missing opening '---' frontmatter"
    fail=1; return
  fi
  # Frontmatter must contain a name: key (skills/agents) — commands use description:.
  local fm
  fm="$(awk 'NR>1{ if($0=="---") exit; print }' "$file")"
  if ! grep -Eq '^(name|description):' <<<"$fm"; then
    echo "FAIL [$file]: frontmatter has no name:/description: key"
    fail=1
  fi
}

while IFS= read -r f; do
  count=$((count+1))
  check_frontmatter "$f"
  refs="$(grep -oE '@?[A-Za-z0-9_.-]+(/[A-Za-z0-9_.-]+)+\.(md|ts|go|yaml|yml|json|sh)' "$f" | sed 's/^@//' | sort -u)"
  for ref in $refs; do
    case "$ref" in http*|*node_modules*) continue ;; esac
    if [ ! -e "$ROOT/$ref" ]; then echo "FAIL [$f]: broken reference -> $ref"; fail=1; fi
  done
done < <(find .claude/skills .claude/agents .claude/commands -name '*.md' 2>/dev/null)

if [ "$fail" -eq 0 ]; then
  echo "PASS: $count artifact(s), frontmatter valid, all references resolve"
else
  echo "FAIL: fix the issues above"
fi
exit $fail

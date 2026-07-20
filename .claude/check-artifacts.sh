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
  local fm
  fm="$(awk 'NR>1{ if($0=="---") exit; print }' "$file")"
  # Branch on artifact type: skills/agents require name:, commands require description:.
  case "$file" in
    .claude/skills/*|.claude/agents/*)
      if ! grep -Eq '^name:' <<<"$fm"; then
        echo "FAIL [$file]: frontmatter missing required 'name:' key (skills/agents must be self-triggerable)"
        fail=1
      fi
      ;;
    .claude/commands/*)
      if ! grep -Eq '^description:' <<<"$fm"; then
        echo "FAIL [$file]: frontmatter missing required 'description:' key (commands must have a description)"
        fail=1
      fi
      ;;
  esac
}

while IFS= read -r f; do
  count=$((count+1))
  check_frontmatter "$f"
  # Extract slash-path refs (with or without leading @) and bare @-prefixed refs (no slash).
  refs="$(
    { grep -oE '@?[A-Za-z0-9_.-]+(/[A-Za-z0-9_.-]+)+\.(md|ts|go|yaml|yml|json|sh)' "$f" | sed 's/^@//';
      grep -oE '@[A-Za-z0-9_.-]+\.(md|ts|go|yaml|yml|json|sh)' "$f" | sed 's/^@//'; } | sort -u
  )"
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

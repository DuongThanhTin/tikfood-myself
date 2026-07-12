#!/usr/bin/env bash
# Verify AI-OS semantic integrity:
#   Check 1 — Anchor integrity: each skill's wrapped recipe still contains ## Steps + ## Verification.
#   Check 2 — Code drift: rooted paths cited in recipes (backtick-quoted apps/, packages/, docs/)
#             still exist on disk.
#
# Limitation (Check 2): only rooted repo paths are tested (apps/, packages/, docs/).
# App-relative paths (e.g. `internal/http`) are deliberately skipped because they are
# ambiguous from repo root and would require per-app prefix resolution.
#
# Usage: bash .claude/check-ai-os.sh
set -uo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$ROOT"

fail=0
skill_count=0
recipe_count=0

# ---------------------------------------------------------------------------
# Check 1 — Anchor integrity (skills -> wrapped recipe still has real procedure)
# ---------------------------------------------------------------------------
while IFS= read -r skill_file; do
  skill_count=$((skill_count + 1))

  # Extract first @docs/recipes/<name>.md pointer (strip leading @).
  recipe_ptr="$(grep -oE '@docs/recipes/[A-Za-z0-9_.-]+\.md' "$skill_file" | head -n1 | sed 's/^@//')"

  if [ -z "$recipe_ptr" ]; then
    echo "FAIL [$skill_file]: no @docs/recipes/*.md authoritative pointer"
    fail=1
    continue
  fi

  recipe_path="$ROOT/$recipe_ptr"
  if [ ! -f "$recipe_path" ]; then
    echo "FAIL [$skill_file]: wrapped recipe missing -> $recipe_ptr"
    fail=1
    continue
  fi

  has_steps=0
  has_verification=0
  if grep -qE '^## Steps[[:space:]]*$' "$recipe_path"; then has_steps=1; fi
  if grep -qE '^## Verification[[:space:]]*$' "$recipe_path"; then has_verification=1; fi

  if [ "$has_steps" -eq 0 ] || [ "$has_verification" -eq 0 ]; then
    echo "FAIL [$skill_file]: wrapped recipe $recipe_ptr is missing required section(s): Steps/Verification"
    fail=1
  fi
done < <(find .claude/skills -name 'SKILL.md' 2>/dev/null)

# ---------------------------------------------------------------------------
# Check 2 — Recipe -> code drift (cited rooted code paths still exist)
# ---------------------------------------------------------------------------
while IFS= read -r recipe_file; do
  recipe_count=$((recipe_count + 1))

  # Extract backtick-quoted tokens starting with apps/, packages/, or docs/.
  # Skips tokens containing a space or a glob (*).
  while IFS= read -r token; do
    # Strip trailing slash for existence check.
    path_to_check="${token%/}"
    if [ -e "$ROOT/$path_to_check" ]; then
      continue
    fi
    echo "FAIL [$recipe_file]: cited path no longer exists -> $path_to_check"
    fail=1
  done < <(
    grep -oE '`(apps|packages|docs)/[^` ]+`' "$recipe_file" \
      | sed "s/^\`//; s/\`$//" \
      | grep -v '[* ]' \
      | sort -u
  )
done < <(find docs/recipes -name '*.md' 2>/dev/null)

# ---------------------------------------------------------------------------
# Summary
# ---------------------------------------------------------------------------
if [ "$fail" -eq 0 ]; then
  echo "PASS: ${skill_count} skill(s), ${recipe_count} recipe(s) — anchors intact, no code drift"
else
  echo "FAIL: fix the issues above"
fi
exit $fail

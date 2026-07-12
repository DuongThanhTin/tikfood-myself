#!/usr/bin/env bash
# record-verify.sh <surface> — write a PASS record for <surface> under .verify-evidence/.
#
# Called by the Makefile at the END of each `verify-*` target, so it only runs when
# every tool in that target exited 0 (make fails fast). One file per surface, so each
# surface carries its own mtime — the verify-evidence-guard.py hook compares that mtime
# against the changed source files to prove verify ran AFTER the last edit.
#
#   .verify-evidence/<surface>   content: <utc-timestamp>\t<short-head>\tPASS
set -euo pipefail

surface="${1:?usage: record-verify.sh <surface>}"
root="${CLAUDE_PROJECT_DIR:-$(git rev-parse --show-toplevel 2>/dev/null || pwd)}"
dir="$root/.verify-evidence"
ts="$(date -u +%Y-%m-%dT%H:%M:%SZ)"
head="$(git -C "$root" rev-parse --short HEAD 2>/dev/null || echo nogit)"

mkdir -p "$dir"
printf '%s\t%s\tPASS\n' "$ts" "$head" > "$dir/$surface"

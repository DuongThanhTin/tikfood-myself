#!/usr/bin/env python3
"""Verify-evidence PreToolUse guard for interactive Claude Code.

Makes "I verified" non-fakeable: before a `git commit` or `gh pr create`, every code
surface touched by the working tree must have a FRESH pass recorded in .verify-evidence/
(written by `make verify-*` via record-verify.sh). "Fresh" = the surface's evidence file
is newer than the changed source files it covers — proving verify ran AFTER the last edit.

Tiers:
  GATED — commit / PR with stale-or-missing evidence for a touched surface: denied unless
          TIKFOOD_SKIP_VERIFY_GATE=1 (the documented escape hatch for emergencies).

Contract: read the PreToolUse JSON on stdin. Deny = reason to stderr + exit 2 (Claude Code
blocks the call and shows the reason to the model). Allow = exit 0. On our own parse error,
or when git is unavailable, exit 0 (fail-open) so a guard bug never bricks the session; the
deny paths are covered by verify-evidence-guard.test.sh.
"""
import os, re, subprocess, sys

ROOT = os.environ.get("CLAUDE_PROJECT_DIR") or os.getcwd()
EVIDENCE_DIR = os.path.join(ROOT, ".verify-evidence")

# Each changed path maps to the surface(s) whose `make verify-<surface>` must be fresh.
# packages/** are shared contracts consumed by all three apps, so they require all three.
SURFACE_RULES = [
    ("apps/api/", ["api"]),
    ("apps/web/", ["web"]),
    ("apps/ai-code-runner/", ["runner"]),
    ("packages/", ["api", "web", "runner"]),
]


def _read_stdin_json():
    import json
    try:
        return json.loads(sys.stdin.read())
    except (ValueError, TypeError):
        return None


def _git(*args):
    """Run a git command at ROOT; return stdout (str) or None on failure."""
    try:
        out = subprocess.run(
            ["git", "-C", ROOT, *args],
            capture_output=True, text=True, timeout=10,
        )
    except (OSError, subprocess.SubprocessError):
        return None
    if out.returncode != 0:
        return None
    return out.stdout


def _changed_paths():
    """Working-tree + staged + untracked paths (porcelain), robust to `git add && commit`
    in one command (nothing is staged yet at PreToolUse time)."""
    # --untracked-files=all so new files show individually (git collapses whole
    # untracked dirs to just "apps/" otherwise, which would dodge the surface map).
    out = _git("status", "--porcelain", "--untracked-files=all")
    if out is None:
        return None  # not a repo / git missing -> caller fails open
    paths = []
    for line in out.splitlines():
        if len(line) < 4:
            continue
        rest = line[3:]
        if " -> " in rest:  # rename/copy: take the destination
            rest = rest.split(" -> ", 1)[1]
        paths.append(rest.strip().strip('"'))
    return paths


def _surfaces_for(path):
    for prefix, surfaces in SURFACE_RULES:
        if path.startswith(prefix):
            return surfaces
    # the guards themselves (.claude/hooks/*.sh|.py, .claude/*.sh) are covered by verify-hooks
    if re.match(r"^\.claude/(hooks/.*|[^/]+\.sh)$", path):
        return ["hooks"]
    return []


def _mtime(path):
    try:
        return os.path.getmtime(path)
    except OSError:
        return None


def deny(reason):
    sys.stderr.write("BLOCKED by verify-evidence-guard: " + reason + "\n")
    sys.exit(2)


def main():
    if os.environ.get("TIKFOOD_SKIP_VERIFY_GATE") == "1":
        sys.exit(0)

    data = _read_stdin_json()
    if not isinstance(data, dict):
        sys.exit(0)  # fail-open on our own parse bug
    tool = data.get("tool_name") or data.get("tool") or ""
    ti = data.get("tool_input") or data.get("input") or {}
    if tool != "Bash" or not isinstance(ti, dict):
        sys.exit(0)

    cmd = " ".join((ti.get("command") or "").split())
    low = cmd.lower()
    is_commit = bool(re.search(r"\bgit\s+commit\b", low))
    is_pr = bool(re.search(r"\bgh\s+pr\s+create\b", low))
    if not (is_commit or is_pr):
        sys.exit(0)

    paths = _changed_paths()
    if paths is None:
        sys.exit(0)  # git unavailable -> fail open

    # required surface -> newest mtime of the changed files that need it
    required = {}
    for p in paths:
        for s in _surfaces_for(p):
            m = _mtime(os.path.join(ROOT, p))
            if m is None:  # deleted file: require evidence exists, no freshness compare
                required.setdefault(s, 0.0)
            else:
                required[s] = max(required.get(s, 0.0), m)

    if not required:
        sys.exit(0)  # nothing verify covers changed -> allow (e.g. docs-only)

    action = "gh pr create" if is_pr else "git commit"
    stale, missing = [], []
    for surface, newest_src in sorted(required.items()):
        ev = os.path.join(EVIDENCE_DIR, surface)
        ev_m = _mtime(ev)
        if ev_m is None:
            missing.append(surface)
        elif ev_m < newest_src:
            stale.append(surface)

    if missing or stale:
        parts = []
        if missing:
            parts.append("no verify evidence for: " + ", ".join(missing))
        if stale:
            parts.append("stale evidence (code changed after verify) for: " + ", ".join(stale))
        surfaces = sorted(set(missing + stale))
        cmds = " && ".join("make verify-%s" % s for s in surfaces)
        deny(
            "%s needs fresh proof that tests pass — %s. "
            "Run: %s   (then commit). Emergency override: TIKFOOD_SKIP_VERIFY_GATE=1."
            % (action, "; ".join(parts), cmds)
        )
    sys.exit(0)


if __name__ == "__main__":
    main()

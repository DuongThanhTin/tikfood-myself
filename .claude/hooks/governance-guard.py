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
        or re.fullmatch(frag, base.lstrip("./"))  # basename catch (e.g. **/*.pem at repo root)
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

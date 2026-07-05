# Opening a new terminal / tab

> 🇻🇳 Tiếng Việt: [`../new-terminal.md`](../new-terminal.md)
>
> Use when: you just opened a new terminal tab and want Claude Code (or yourself) to
> understand the repo immediately and keep going. See [`new-feature.md`](new-feature.md)
> to start new work, [`worktrees.md`](worktrees.md) to run in parallel.

## What loads automatically (quick recap)

A new session in the repo auto-reads `CLAUDE.md` + `apps/*/CLAUDE.md`, the AI-OS spine
(`docs/ai/AI-CONTRACT.md` → `docs/REPOSITORY-MAP.md` → `docs/ai/CONTEXT-LOADING.md`), and
auto-memory. You **don't** need to re-paste these — just point at the task.

## Checklist when opening a new tab

```bash
# 1. Enter the repo directory (CLAUDE.md loads from here)
cd /Users/tinduong245/Documents/Myself/tikfood-ai-automation-starter

# 2. See which branch you're on and whether anything is in progress
git status --short
git branch --show-current

# 3. (if needed) refresh the base
git checkout dev && git pull
```

Then paste the kickoff for your purpose (full templates in
[`templates/session-kickoff.md`](templates/session-kickoff.md)).

## Sample kickoffs

**A. Continue / catch up** (new session, want to re-orient):
> Read CLAUDE.md and docs/REPOSITORY-MAP.md. Based on git + memory, summarize the current
> state and what's in progress, then ask me what to do next. Don't change anything yet.

**B. Start a specific task right away:**
> Read CLAUDE.md, docs/REPOSITORY-MAP.md, docs/ai/CONTEXT-LOADING.md.
> Task: **\<one line>**. If it's UI, follow docs/design/. Brainstorm/plan before coding.

## Worked example

**Context:** Monday morning, new tab, forgot what was in progress.

```text
You paste:
> Read CLAUDE.md + docs/REPOSITORY-MAP.md. Based on git status/log and memory,
> summarize in-progress work and suggest the next step.

Claude will:
1. Read the pointed files + auto-memory.
2. Run git status/log to see the branch & latest commit.
3. Reply: "Branch ai/... has commit X (design system docs). In progress: Y.
   Do you want to (a) open a PR, (b) start a new feature, (c) …?"
```

→ You just pick a direction; no need to retell history.

## Notes

- **Don't start new work directly on an in-progress branch** — see
  [`new-feature.md`](new-feature.md) to create a clean branch.
- If two tabs edit the same file → collisions. For true parallelism use
  [worktrees](worktrees.md).

## Related

- [`new-feature.md`](new-feature.md) · [`worktrees.md`](worktrees.md) · [`templates/session-kickoff.md`](templates/session-kickoff.md)
- [`README.md`](README.md) (decision map) · [`CLAUDE.md`](../../../CLAUDE.md)

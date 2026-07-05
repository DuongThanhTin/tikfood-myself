# Template: Kickoff prompt for a new tab

> 🇻🇳 Tiếng Việt: [`../../templates/session-kickoff.md`](../../templates/session-kickoff.md)
>
> Paste one of the templates below into a new Claude Code session. Replace the `<…>`.
> How to use: [`../new-terminal.md`](../new-terminal.md), [`../new-feature.md`](../new-feature.md).

---

## A. Catch up / continue prior work

```text
Read CLAUDE.md and docs/REPOSITORY-MAP.md first.
Based on git status/log and memory, summarize:
- current branch + latest commit,
- what's in progress,
- 2–3 possible next steps.
Then ask me which direction. Don't change anything yet.
```

## B. Start a new feature

```text
Read CLAUDE.md, docs/REPOSITORY-MAP.md, docs/ai/CONTEXT-LOADING.md.
If it's UI: also read docs/design/ and the matching playbook in docs/design/playbooks/.

Feature: <one-line goal>
Scope: <apps/web or apps/api; which part>
Out-of-scope: <what NOT to do>
Acceptance: <what "done" means>
Constraints: discovery-only, no anti-goals; reuse first; use tokens; data through lib/api.ts.

Brainstorm/plan before coding. For multi-step work, write the plan to a file first.
```

## C. Start from an existing brief/spec

```text
Read CLAUDE.md, then read <path to brief or docs/superpowers/specs/…>.
Confirm you understand it, list risks/ambiguities, then propose a plan.
Don't code until I approve the plan.
```

## D. Review a UI change

```text
Read docs/design/playbooks/review-ui.md.
Review the current diff against that checklist (token consistency, reuse, a11y, responsive,
anti-goals). Return must-fix vs nice-to-have, with file:line.
```

---

### Tips

- Always say "**plan before coding**" for anything over one step.
- For parallel work across tabs → create a separate branch/worktree first (see
  [`../worktrees.md`](../worktrees.md)).

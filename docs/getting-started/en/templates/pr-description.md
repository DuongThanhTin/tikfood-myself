# Template: PR Description

> 🇻🇳 Tiếng Việt: [`../../templates/pr-description.md`](../../templates/pr-description.md)
>
> Use when opening a PR (branch `ai/<feature>` → `main`/`dev`). No auto-merge; wait for
> human review. How to use: [`../new-feature.md`](../new-feature.md).

---

```md
## Goal
<1–2 lines: what this feature does / why>

## Changes
- <main change 1>
- <main change 2>

## Scope & Out-of-scope
- Belongs to: apps/web | apps/api | …
- Does NOT include: <…>

## Acceptance / how to test
- [ ] <test step 1 / expected result>
- [ ] <step 2>
- [ ] Ran relevant tests; checked against docs/verification/definition-of-done.md

## Constraints honored
- [ ] Discovery-only, no anti-goals
- [ ] Reuse/tokens/`lib/api.ts`; no breaking API/response changes
- [ ] Vietnamese copy; a11y (labels on icon-only controls, image alt)
- [ ] Needs human approval? <no / yes — reason>

## Screenshots / notes (if UI)
<screenshot or description>

🤖 Generated with [Claude Code](https://claude.com/claude-code)
```

---

### Important reminders (repo rules)

- PR branches under `ai/`; **never** push to `main`/`master`, **no** force-push, **no** auto-merge.
- Commit message ends with: `Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>`
- PR body ends with: `🤖 Generated with [Claude Code](https://claude.com/claude-code)`

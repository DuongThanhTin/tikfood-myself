# Getting Started — Working in this repo (with Claude Code)

> 🇻🇳 Tiếng Việt: [`../README.md`](../README.md)
>
> **This is the front door.** Just opened a new terminal, about to start a feature, or
> forgot "where do I begin"? Read the one file below and follow it. **Detailed worked
> examples** live in this file; **copy-paste templates** live in [`templates/`](templates/).

## I want to… → read this file

| Situation | Guide | Detailed example |
|-----------|-------|------------------|
| Just **opened a new terminal/tab** | [`new-terminal.md`](new-terminal.md) | [Example 1](#example-1--new-terminal-catch-up--continue) |
| **Start a new feature** (branch → brief → PR) | [`new-feature.md`](new-feature.md) | [Example 2](#example-2--build-a-feature-end-to-end) |
| **Run several tabs in parallel** without collisions | [`worktrees.md`](worktrees.md) | [Example 3](#example-3--two-features-in-parallel-with-worktrees) |
| Need a **kickoff / brief / PR** template to paste | [`templates/`](templates/) | shown in each example |

## What loads AUTOMATICALLY each session (so you don't lose context)

Every new tab is an independent session, **but it doesn't start from zero** — it always loads:

1. **`CLAUDE.md`** (root) + **`apps/*/CLAUDE.md`** — the rules, anti-goals, structure.
2. **AI-OS spine**: [`docs/ai/AI-CONTRACT.md`](../../ai/AI-CONTRACT.md) → [`docs/REPOSITORY-MAP.md`](../../REPOSITORY-MAP.md) → [`docs/ai/CONTEXT-LOADING.md`](../../ai/CONTEXT-LOADING.md).
3. **Auto-memory** (`MEMORY.md`) — long-lived decisions/progress across sessions.

→ A new session sees **results & decisions**, just not the previous **chat transcript**.
So: write important things to **files** (brief/spec/memory), don't keep them in the chat.

## 6 golden rules (details in `CLAUDE.md` — don't restate, just remember)

1. **Reuse first** — search code/`docs/` before adding anything.
2. **Use tokens, not hardcoded values** — follow [`docs/design/`](../../design/00-overview.md) for UI.
3. **Vietnamese copy**, English identifiers.
4. **Data through `lib/api.ts`** — no scattered `fetch` in components.
5. **Anti-goals** — no cart/order/checkout/payment/booking/chat/follow/monetization/livestream.
6. **Finish against** [`docs/verification/definition-of-done.md`](../../verification/definition-of-done.md).

## Layout of this folder

```
docs/getting-started/
├─ README.md            Vietnamese hub + map + 3 detailed examples
├─ new-terminal.md      what to do in a fresh tab
├─ new-feature.md       start a feature (branch → brief → PR)
├─ worktrees.md         run in parallel with git worktree
├─ templates/           kickoff / feature-brief / pr-description
└─ en/                  ← this English mirror
```

---

# ⭐ Detailed worked examples (end-to-end)

> Each example follows the same shape: **Context → Steps (commands + prompt to paste) →
> Expected result**. The "Claude replies" blocks are **illustrative**, not a real transcript.

## Example 1 — New terminal (catch up & continue)

**Context:** You open a new tab in the morning, forgot what was in progress, want to continue safely.

**Step 1 — Enter the repo & check state**
```bash
cd /Users/tinduong245/Documents/Myself/tikfood-ai-automation-starter
git status --short
git branch --show-current
```

**Step 2 — Paste the kickoff** (variant A in [`templates/session-kickoff.md`](templates/session-kickoff.md)):
```text
Read CLAUDE.md and docs/REPOSITORY-MAP.md first.
Based on git status/log and memory, summarize:
- current branch + latest commit,
- what's in progress,
- 2–3 possible next steps.
Then ask me which direction. Don't change anything yet.
```

**Expected result (illustrative):**
```text
Claude:
• Branch: ai/ai-operating-system, latest commit 9b038a9 "docs(design): …".
• In progress: 6 docs files (architecture, services/*, standards/*) modified earlier, NOT committed.
• Next: (a) commit/PR the design-system docs, (b) start a new feature,
  (c) handle the 6 in-progress files. Which one?
```
→ You just type e.g. `b`, then jump to [Example 2](#example-2--build-a-feature-end-to-end).

**Note:** don't start new work on an in-progress branch — create a clean branch (see Example 2).

---

## Example 2 — Build a feature end-to-end

**Context:** Add a **"Clear filters"** button to the venue rail's empty-state (when there
are no results). Small feature, `apps/web` only.

**Step 1 — Clean branch**
```bash
git checkout dev && git pull
git checkout -b ai/rail-empty-state-reset
```

**Step 2 — Write the brief** (see [`templates/feature-brief.md`](templates/feature-brief.md)):
```md
# Feature: Clear-filters button in the venue rail empty-state
- Goal: Empty rail → show message + a "Clear filters" button to get results back.
- Belongs to: apps/web
- Scope: add the button to the empty-state; click → run clearFilters + re-search.
- Out-of-scope: API changes; pagination; illustration.
- Acceptance:
  - [ ] Empty → shows .emptyText + button.
  - [ ] Click → clearFilters() runs, list reloads.
  - [ ] Works in light/dark and both breakpoints (1180/640).
- Constraints: discovery-only; follow docs/design; no new component if reusable.
- Track: docs/design/patterns/empty-state.md + playbooks/update-component.md
- Branch: ai/rail-empty-state-reset (base: dev)
```

**Step 3 — Kickoff Claude** (variant B in session-kickoff, embedding the brief):
```text
Read CLAUDE.md, apps/web/CLAUDE.md, docs/design/patterns/empty-state.md,
docs/design/playbooks/update-component.md.
Here is the brief: <paste the brief above>.
Plan before coding (list files touched, how to reuse clearFilters). Don't code until I approve.
```

**Step 4 — Track:** since this **extends an existing pattern** (not a new component) →
follow `update-component.md`; the button reuses the existing `secondaryButton`; behavior
reuses `clearFilters` (already in `DiscoveryExperience.tsx`).

**Step 5 — Build (Claude does this after you approve the plan):** edit the rail empty-state
block to render `.emptyText` + `<button className="secondaryButton" onClick={clearFilters}>`;
no new tokens/colors; Vietnamese copy "Xoá bộ lọc".

**Step 6 — Finish & PR**
```bash
# check against docs/verification/definition-of-done.md, then:
git add apps/web/components/DiscoveryExperience.tsx
git commit -m "feat(web): clear-filters button in rail empty-state

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
git push -u origin ai/rail-empty-state-reset
gh pr create --base dev --fill   # or paste body from templates/pr-description.md
```
PR body (short form of [`templates/pr-description.md`](templates/pr-description.md)):
```md
## Goal
Empty rail → a "Clear filters" button to get results back.
## Changes
- Rail empty-state: add secondaryButton calling clearFilters.
## Acceptance / how to test
- [ ] Filter down to empty → see message + button; click → list reloads.
- [ ] OK in light/dark + both breakpoints.
🤖 Generated with [Claude Code](https://claude.com/claude-code)
```

**Expected result:** a small, in-scope PR that reuses existing code, no breaking API changes.

---

## Example 3 — Two features in parallel with worktrees

**Context:** Simultaneously build (A) a new endpoint in `apps/api` and (B) a filter-chip
tweak in `apps/web`. You want two independent tabs, no stashing back and forth.
Details: [`worktrees.md`](worktrees.md).

**Tab 1 — backend (`apps/api`)**
```bash
git worktree add -b ai/api-nearby-endpoint ../tikfood-api-nearby dev
cd ../tikfood-api-nearby
# open Claude in this tab, paste:
#   Read CLAUDE.md + apps/api/CLAUDE.md. Feature: "nearby venues" endpoint.
#   Plan first; needs human approval if it touches external network/migration.
```

**Tab 2 — web (`apps/web`)** — open another terminal:
```bash
git worktree add -b ai/web-filter-chip ../tikfood-web-chip dev
cd ../tikfood-web-chip/apps/web && npm install    # a new worktree needs its own deps
# open Claude in this tab, paste:
#   Read CLAUDE.md + apps/web/CLAUDE.md + docs/design/.
#   Feature: "Open late" filter chip. Reuse the existing chip/tag (docs/design/components/input.md).
```

**Check & clean up**
```bash
git worktree list                       # see open worktrees
# when each is done: commit → PR → then:
git worktree remove ../tikfood-api-nearby
git worktree remove ../tikfood-web-chip
git worktree prune
```

**Expected result:** two directories, two `ai/…` branches, two PRs — no file collisions,
no `git stash`.

**Common traps:** forgetting `npm install` in the web worktree (build fails, missing
module); the same branch can't be checked out in two worktrees.

---

## Related

- Rules & context: [`CLAUDE.md`](../../../CLAUDE.md) · [`docs/ai/`](../../ai)
- Recipes by task type: [`docs/recipes/`](../../recipes)
- UI playbooks: [`docs/design/playbooks/`](../../design/playbooks/README.md)
- Spec-driven trail: [`docs/superpowers/specs/`](../../superpowers/specs)
- Full onboarding narrative: [`docs/handbook/`](../../handbook)

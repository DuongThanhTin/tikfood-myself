# Running in parallel with git worktree

> 🇻🇳 Tiếng Việt: [`../worktrees.md`](../worktrees.md)
>
> Use when: you want to open **several terminal tabs doing different work at once** without
> stepping on each other (no constant `git stash`/branch-switching in one directory).

## The problem & the fix

One repo directory is on **one branch** at a time. Two tabs in the same directory doing two
different things will collide. **`git worktree`** lets one repo have **multiple working
directories**, each on its own branch → one worktree per tab, fully independent.

## Core commands

```bash
# Create a new worktree + new branch for a feature (base = dev)
git worktree add -b ai/<feature> ../tikfood-<feature> dev

# List existing worktrees
git worktree list

# Remove when done (after merge / no longer needed)
git worktree remove ../tikfood-<feature>
git worktree prune          # clean up leftover metadata
```

> Note: all worktrees **share one .git/history**. A given branch can be checked out in
> only **one** worktree at a time.

## Conventions for this repo

- Put worktrees **next to** the repo, not nested inside: `../tikfood-<feature>`.
- Branches still use the **`ai/<feature>`** prefix (per `CLAUDE.md`); base off `dev`.
- Each worktree needs its own dependencies (no shared `node_modules`, `.next`, or Go build):
  ```bash
  cd ../tikfood-<feature>/apps/web && npm install
  ```
- **Split work so tabs don't overlap files** (e.g. tab A = `apps/api`, tab B = `apps/web`).

## Worked example — two features in parallel

**Goal:** at the same time build (1) a new discovery endpoint in `apps/api`, and (2) a UI
tweak in `apps/web`.

```bash
# Tab 1 — backend
git worktree add -b ai/api-nearby-endpoint ../tikfood-api-nearby dev
cd ../tikfood-api-nearby
# open Claude here: "Read CLAUDE.md + apps/api/CLAUDE.md. Feature: <...>. Plan first."

# Tab 2 — web (open another terminal)
git worktree add -b ai/web-filter-chip ../tikfood-web-chip dev
cd ../tikfood-web-chip/apps/web && npm install
# open Claude here: "Read CLAUDE.md + apps/web/CLAUDE.md + docs/design/. Feature: <...>."
```

→ Two tabs, two directories, two branches, two PRs — no stash, no collisions.
When each is done: commit → PR → `git worktree remove ../tikfood-...`.

## Inside a Claude session

You can also ask Claude to work in an isolated worktree (it has a built-in worktree
mechanism for collision-prone tasks). But for **opening several tabs by hand**, the
`git worktree` CLI above is the clearest option.

## Common errors

- *"fatal: '<branch>' is already checked out"* → that branch is open in another worktree;
  use a different branch or `git worktree list` to find it.
- Forgetting `npm install` in a new worktree → the web build fails on missing modules.
- Deleting a worktree directory by hand without `git worktree prune` → stale metadata.

## Related

- [`new-terminal.md`](new-terminal.md) · [`new-feature.md`](new-feature.md) · [`README.md`](README.md)
- Branch/commit rules: [`CLAUDE.md`](../../../CLAUDE.md) (Security & Git)

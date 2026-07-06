# Context-Loading Protocol

> **Why this file exists:** `docs/handbook/01` §1.4 gives a default *order* for loading
> context, and handbook Part 6 (Context Engineering) explains the principle but defers
> the concrete protocol to this file. Without a single operational protocol, every
> session re-improvises what to read — loading too little
> (missing a rule) or too much (burning context on irrelevant files). This is the
> operational protocol: **what to read, in what order, per task type, and when to
> drop context or delegate to a sub-agent.** It is the runtime companion to
> [`AI-CONTRACT.md`](AI-CONTRACT.md) (the rules) and [`../REPOSITORY-MAP.md`](../REPOSITORY-MAP.md)
> (where things live).

**Principle:** load the *minimum* context that lets you act correctly. Prefer reading
a small owning doc over guessing. When a search would sweep many files, delegate it to
a sub-agent and keep only the conclusion (see §4).

## 1. Always-on core (every session)

Read these first, in order — they are small and set the boundaries:

1. [`/CLAUDE.md`](../../CLAUDE.md) — workspace rules & anti-goals.
2. [`docs/ai/AI-CONTRACT.md`](AI-CONTRACT.md) — the rules-of-engagement index.
3. [`docs/REPOSITORY-MAP.md`](../REPOSITORY-MAP.md) — to locate the task's area.

Stop there until you have classified the task (§2). Do **not** pre-load app code or
standards you may not need.

## 2. Task-type loading matrix

After classifying the task (use [`docs/thinking/README.md`](../thinking/README.md) §1),
load the app-specific and topic docs it needs. Read the owning `CLAUDE.md` **before**
editing any app. For a common task type, the matching [`docs/recipes/`](../recipes/)
recipe already bundles the skill + docs + steps — prefer it over reconstructing from
this matrix. For the whole **request → PR** path (classify → context → spec/plan →
branch → implement → verify → PR, with approval gates), run the matching
[`docs/workflows/`](../workflows/) workflow — it orchestrates the recipe, skills, and
gates end-to-end.

| Task type | Read before acting (in addition to the core) |
| --- | --- |
| **Backend change** (`apps/api`) | [`apps/api/CLAUDE.md`](../../apps/api/CLAUDE.md) → [`docs/standards/backend/README.md`](../standards/backend/) → the specific backend standard (architecture / patterns / errors / request-response / testing) for the file you touch |
| **Frontend change** (`apps/web`) | [`apps/web/CLAUDE.md`](../../apps/web/CLAUDE.md) → [`docs/standards/frontend-architecture.md`](../standards/frontend-architecture.md); keep `lib/api.ts` types in sync with Go JSON tags |
| **API contract / response shape** | [`docs/contracts/api.md`](../contracts/api.md) + [`docs/standards/backend/request-response.md`](../standards/backend/request-response.md); update `apps/web/lib/api.ts` too |
| **Runner / automation** (`apps/ai-code-runner`) | [`docs/contracts/runner.md`](../contracts/runner.md) + [`docs/agents/`](../agents/) + [`packages/config/*.yaml`](../../packages/config/) |
| **Runtime prompt** (`packages/prompts`) | the paired [`docs/agents/*`](../agents/) role doc + [`docs/standards/prompt-engineering.md`](../standards/prompt-engineering.md); output must be JSON-only |
| **Schema / contract** (`packages/schemas`) | **Protected path** — get human approval; then [`docs/contracts/runner.md`](../contracts/runner.md) |
| **Migration** (`apps/api/migrations`) | **Requires human approval**; then [`docs/standards/backend/database.md`](../standards/backend/database.md) |
| **Product scope / feature framing** | [`docs/tikfood/`](../tikfood/) (vision / mvp-scope / anti-goals) |
| **Architecture question** | [`docs/architecture.md`](../architecture.md) + [`docs/services/`](../services/) + [`docs/adr/`](../adr/) |
| **Security-sensitive change** | [`docs/security.md`](../security.md) (runner) + [`docs/standards/application-security.md`](../standards/application-security.md) (product) |

If a task spans areas, load each area's owner doc — but only the sections you need.

## 3. Process layer (which skill first)

Load the **process skill before implementation**, per handbook §1.5:

- "Build / add X" → `brainstorming` → `writing-plans` → `executing-plans`.
- "Fix bug Y" → `systematic-debugging`.
- Before claiming done → `verification-before-completion`.

Hard skills (TDD, systematic-debugging) are followed exactly; pattern skills adapt.
Recipes in [`docs/recipes/`](../recipes/) name the exact skill + docs per task type.

## 4. When to drop context or delegate

- **Drop context** once a file has served its purpose (e.g. you confirmed a signature)
  — do not keep large files resident "just in case."
- **Delegate to a sub-agent** when answering a question means sweeping many files or
  directories and you only need the conclusion (naming conventions, "where is X used",
  broad audits). Keep the summary, not the file dumps.
- **Spawn parallel sub-agents** for independent read-only investigations; run them
  concurrently and synthesize.
- **Use a git worktree** for isolated feature work per the `using-git-worktrees` skill.
- Prefer the dedicated read/search tools over shell `cat`/`grep` for single files you
  actually need in context.
- **Read ranges, not whole files** — for a large file (>~400 lines) where you need one
  function/section, read the relevant range instead of the entire file.

## 5. Order of authority when sources conflict

Code > `CLAUDE.md` (app-specific > root, within the app) > `.ai-agent.yaml` /
`packages/config` > `docs/security.md` > `docs/standards/**` > handbook/other docs.
When docs disagree with code, **code wins** — flag the mismatch rather than following
stale docs. See [`AI-CONTRACT.md`](AI-CONTRACT.md) → *Owning sources*.

## 6. Session-start checklist

- [ ] Read the always-on core (§1).
- [ ] Classify the task and load only its row from the matrix (§2).
- [ ] Load the correct process skill first (§3).
- [ ] Confirm the change is not in a protected/anti-goal area (else stop → human approval).
- [ ] Do not write code before Discovery is complete (per the Work Session Playbook).

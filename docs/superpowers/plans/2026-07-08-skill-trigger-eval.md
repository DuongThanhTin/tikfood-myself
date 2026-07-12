# Skill-Trigger Eval Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add `/eval-skills` — a Claude-run behavioral eval that checks whether the `.claude/` skills self-trigger correctly and don't collide, using a judge-on-descriptions proxy.

**Architecture:** Three markdown files. `.claude/eval/scenarios.md` is a declarative task→expected-skill table. `.claude/commands/eval-skills.md` is the runner Claude executes: it reads every skill's `description` live, dispatches one fresh **blind** judge subagent per scenario (task + descriptions only), compares each verdict to `expected`, and prints a scorecard. `.claude/eval/README.md` documents purpose and honest limitations. No shell/CI gate — the eval is opt-in and human-reviewed because it is probabilistic.

**Tech Stack:** Markdown + YAML frontmatter (Claude Code command); subagent dispatch as the judge mechanism (session model, no API key, no network).

## Global Constraints

- **No `make` target; not in `make verify`.** Opt-in, human-reviewed (probabilistic proxy).
- **Proxy, not the real trigger engine.** Tests description-discriminability, not live auto-invocation.
- **Skills only (v1).** Do not eval command/agent triggering.
- **No new skills/recipes.** This adds an eval only.
- **Judges are fresh and blind** — each judge sees the task + the `{skill: description}` map, never the expected answer, never conversation history.
- **Descriptions read live** from `.claude/skills/*/SKILL.md` at run time — never hardcoded into the command.
- **`expected` is a set;** `none` means no skill should fire (catches over-triggering). `compose-field-web` expects `contract-first-change`.
- **`.claude/eval/*` files carry no frontmatter** (they are not under skills/agents/commands, so `check-artifacts.sh` does not scan them). The command file **is** scanned: it needs `description:` frontmatter and every `@`-ref must resolve.
- **Branch:** `ai/claude-executable-layer`. Never push `main`, never force-push, never auto-merge.

---

### Task 1: Author the eval (scenarios + command + README)

The whole eval is one cohesive deliverable. Create all three files, then verify the command passes the syntactic gate.

**Files:**
- Create: `.claude/eval/scenarios.md`
- Create: `.claude/commands/eval-skills.md`
- Create: `.claude/eval/README.md`

**Interfaces:**
- Consumes: the `description` frontmatter of `.claude/skills/*/SKILL.md` (read at run time); `.claude/check-artifacts.sh` (Task gate).
- Produces: the `/eval-skills` command and its scenario data.

- [ ] **Step 1: Create the scenario table**

Create `.claude/eval/scenarios.md`:

```markdown
# Skill-trigger eval — scenarios

> Data for `/eval-skills`. Each row: a task a user might describe, and the skill that
> SHOULD fire (`expected`). `expected` is a set — comma-separate if more than one answer is
> acceptable. `none` means no skill in `.claude/skills/` should fire (tests over-triggering).
> Add a row to cover a new skill or a new collision you want guarded. Keep tasks realistic
> and one sentence.

| id | task | expected |
|----|------|----------|
| pos-add-endpoint | Add a query filter to the venues endpoint in `apps/api` | add-api-endpoint |
| compose-field-web | Add a field to the venue API response that the web app will display | contract-first-change |
| pos-envelope | Change the `{data,error}` response shape used across api and web | contract-first-change |
| boundary-schema | Add a field to a `packages/schemas` contract | contract-first-change |
| neg-refactor | Rename a Go handler function for readability, no behavior change | none |
| neg-frontend-copy | Fix a typo in the login page heading | none |
```

- [ ] **Step 2: Create the `/eval-skills` command**

Create `.claude/commands/eval-skills.md`:

```markdown
---
description: Run the skill-trigger eval — check whether .claude/ skills self-trigger correctly and do not collide, using a judge-on-descriptions proxy. Opt-in, human-reviewed, probabilistic; not a CI gate.
---

# /eval-skills — behavioral skill-trigger eval

Purpose, mechanism, and honest limitations: **@.claude/eval/README.md** — read it first.
Scenario data: **@.claude/eval/scenarios.md**.

You are the eval runner. Execute these steps and print a scorecard. Do not modify any file.

## Steps
1. **Collect skills.** Read the `description` frontmatter from every
   `.claude/skills/*/SKILL.md`. Build the map `{skill-name: description}`. If there are zero
   skills, print `No skills to eval` and stop.
2. **Load scenarios** from `@.claude/eval/scenarios.md` — parse the table rows into
   `{id, task, expected}` (split `expected` on commas into a set; `none` = empty set meaning
   "no skill should fire").
3. **Judge each scenario independently.** For every scenario, dispatch ONE fresh
   general-purpose subagent (model: haiku; run scenarios concurrently) with EXACTLY this
   brief and nothing else — the judge must stay blind to the expected answer:
   > You are a routing judge. Given a developer task and a set of skills (name +
   > description), decide which SINGLE skill should automatically handle the task, or
   > `none` if no skill fits. Consider only the descriptions provided. Reply with two
   > lines: `pick: <skill-name|none>` and `why: <one sentence>`. Do not invent skill names.
   >
   > Task: "<scenario.task>"
   >
   > Skills:
   > <for each: "- <name>: <description>">
4. **Score.** For each scenario, parse the judge's `pick`. PASS if: (a) `expected` is
   non-empty and `pick` is in `expected`; or (b) `expected` is empty (`none`) and `pick` is
   `none`. Otherwise FAIL.
5. **Report a scorecard** — one line per scenario:
   `PASS/FAIL  <id>  expected=<...>  pick=<...>  — <judge why>`
   then a summary line `Score: <passed>/<total>`. For each FAIL, add one line naming the
   likely cause (e.g. "two descriptions both match — tighten one to disambiguate", or
   "description too broad — it matched an unrelated task").

## Reading the result
A FAIL is a **finding about the skill descriptions**, not a bug in the eval: it means the
descriptions are not discriminative enough for reliable auto-triggering. This is opt-in and
probabilistic — re-run if a single scenario looks borderline; treat a persistent FAIL as a
description to fix.
```

- [ ] **Step 3: Create the README**

Create `.claude/eval/README.md`:

```markdown
# .claude/eval — skill-trigger eval

Behavioral eval for the `.claude/` executable layer, run via `/eval-skills`. It answers a
question the syntactic gate (`check-artifacts.sh`) and the integrity gate
(`check-ai-os.sh`) cannot: **do the skills actually self-trigger for the right tasks, and
do overlapping skills collide?**

## How it works
For each scenario in `scenarios.md`, a fresh **blind** judge subagent is given the task plus
every skill's `description` (read live) and asked which single skill should fire (or
`none`). Its pick is compared to the scenario's `expected`. The runner prints a scorecard.

## Honest limitations
- **Proxy, not the real trigger engine.** It tests whether the `description`s are
  discriminative enough for a model to route correctly — not Claude Code's actual
  auto-invocation in a live session. Confirm the real thing with a manual live test when it
  matters.
- **Probabilistic.** The judge is an LLM; a borderline scenario may flip between runs.
  Re-run before treating a single FAIL as real.
- **Opt-in, human-reviewed. Deliberately NOT in `make verify`.** A flaky hard gate is worse
  than none (see the V2 vision, "pristine gates"). Run it when you add or change a skill.
- **Skills only (v1).** Command/agent triggering is not evaluated yet.

## Adding a scenario
Add a row to `scenarios.md`: a realistic one-sentence task and the `expected` skill (or
`none`). Cover each new skill with at least one positive, and add a `none` row if the new
skill's description risks grabbing unrelated tasks.
```

- [ ] **Step 4: Verify the syntactic gate passes**

Run:
```bash
make verify-claude
```
Expected: `PASS: 3 artifact(s), frontmatter valid, all references resolve`. The count rises from 2 to 3 because `.claude/commands/eval-skills.md` is now scanned (it has `description:` frontmatter and its `@`-refs — `@.claude/eval/README.md`, `@.claude/eval/scenarios.md` — must resolve). The two `.claude/eval/*.md` files are NOT scanned (not under skills/agents/commands) and correctly have no frontmatter.

- [ ] **Step 5: Verify the integrity gate is unaffected**

Run:
```bash
make verify-ai-os
```
Expected: `PASS: 2 skill(s), 9 recipe(s) — anchors intact, no code drift` (unchanged — this task adds no skill or recipe).

- [ ] **Step 6: Commit**

```bash
git add .claude/eval/scenarios.md .claude/commands/eval-skills.md .claude/eval/README.md
git commit -m "feat(claude): /eval-skills behavioral skill-trigger eval

- .claude/eval/scenarios.md: task -> expected-skill scenarios (positives, compose, negatives, boundary)
- .claude/commands/eval-skills.md: runner — blind judge subagent per scenario, live descriptions, scorecard
- .claude/eval/README.md: proxy/probabilistic/opt-in limitations

Co-Authored-By: Claude Opus 4.8 (1M context) <noreply@anthropic.com>"
```

---

### Task 2: Behavioral verification (controller-run)

Running the eval dispatches judge subagents, which the main session does best. This is the controller's verification step, analogous to the spine's live-trigger tests — not delegated to an implementer.

**Files:** none (verification only; do not persist the scorecard — it is a point-in-time report, and a committed scorecard would go stale).

- [ ] **Step 1: Run the eval**

Invoke `/eval-skills` in the session. It reads the two current skills' descriptions, runs the six scenarios through blind haiku judges, and prints a scorecard.

- [ ] **Step 2: Interpret the scorecard**

Expected healthy result: `Score: 6/6`, with the two `neg-*` scenarios returning `pick=none`, `compose-field-web` and `boundary-schema` returning `contract-first-change`, and `pos-add-endpoint` returning `add-api-endpoint`.

- [ ] **Step 3: Triage any FAIL**

A FAIL is a finding about the **skill descriptions**, not the eval. Record it: which scenario, expected vs pick, and the likely description overlap. Do **not** silently "fix" the eval to make it green. Two legitimate responses: (a) tighten the losing skill's `description` (a real improvement — then re-run), or (b) if the collision is intended-and-acceptable, note it in the scenario row. Report the outcome to the human; do not auto-edit skill descriptions without surfacing the finding first.

- [ ] **Step 4: Report**

Report the scorecard and any findings to the human. The eval's value is precisely a FAIL that surfaces a real description collision before it bites in a live session.

---

## Self-Review

**Spec coverage:**
- `.claude/eval/scenarios.md` (6 seeded scenarios incl. compose + negatives + boundary) → Task 1 Step 1. ✅
- `.claude/commands/eval-skills.md` (live descriptions, blind judge per scenario, scorecard) → Task 1 Step 2. ✅
- `.claude/eval/README.md` (proxy/probabilistic/opt-in/skills-only limitations) → Task 1 Step 3. ✅
- No make target / not in `make verify` → no Makefile change in any task; README states it. ✅
- Judges fresh + blind, descriptions read live → command Steps 1 & 3. ✅
- `expected` is a set, `none` = over-trigger check → command Step 4 scoring. ✅
- `compose-field-web` → `contract-first-change` → scenarios table. ✅
- Behavioral run is controller work → Task 2. ✅
- Syntactic gate still passes with the new command → Task 1 Step 4. ✅

**Placeholder scan:** All three files' full content is in the plan; commands have expected output; no TBD/TODO. ✅

**Name consistency:** File paths, the `/eval-skills` name, `expected`/`pick`/`none` vocabulary, and the `make verify-claude` → `PASS: 3 artifact(s)` expectation are consistent across tasks and the command body. ✅

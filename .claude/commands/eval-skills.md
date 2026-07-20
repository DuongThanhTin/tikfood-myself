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

# Design — Skill-Trigger Eval (`/eval-skills`)

**Date:** 2026-07-08 · **Branch:** `ai/claude-executable-layer` · **Status:** approved for planning

## Problem

The `.claude/` layer's whole premise is that skills **self-trigger** from their
`description` frontmatter. Nothing tests whether they actually do. The two shipped gates
check *syntax* (`check-artifacts.sh`: frontmatter + links) and *semantic integrity*
(`check-ai-os.sh`: recipe anchors + code drift) — neither checks **behavior**: given a
task, does the right skill fire, and do overlapping skills mis-fire?

The final spine review flagged the concrete risk: `add-api-endpoint` and
`contract-first-change` overlap ("add a field to the venue response the web reads" could
match both). Today the only detector is a human hitting the collision live. The V2 vision
promised a behavioral eval (§2.3) but the integrity harness delivered only the
deterministic part. This spec closes that gap: gaps #1 (reduce human-verification burden)
and #2 (deliver the behavioral eval, not just the easy 80%).

## Non-goals

- **No `make` target / not in `make verify`.** A behavioral eval needs a model in the loop
  (bash can't judge; an external API needs a secret key + network, both off-limits). It is
  probabilistic, so a hard CI gate would be flaky — worse than none. It is **opt-in and
  human-reviewed**.
- **Not the real Claude Code trigger engine.** This is a *proxy*: it tests whether the
  `description`s are discriminative enough for a model to pick the right skill. It does not
  spawn live sessions to observe actual auto-invocation.
- **Skills only (v1).** Command/agent triggering is out of scope for this batch.
- **No new skills/recipes.** This adds an eval, not more artifacts to trigger.

## Architecture — a Claude-run procedure, three files

The "runner" is Claude executing a command, not a shell script. The model in the loop is
the session/subagent model — no API key, no network.

1. **`.claude/eval/scenarios.md`** — declarative scenario table (human-editable):

   | id | task | expected |
   |----|------|----------|
   | `pos-add-endpoint` | Add a query filter to the venues endpoint in `apps/api` | `add-api-endpoint` |
   | `compose-field-web` | Add a field to the venue API response the web will display | `contract-first-change` |
   | `pos-envelope` | Change the `{data,error}` shape used across api and web | `contract-first-change` |
   | `neg-refactor` | Rename a Go handler for readability, no behavior change | `none` |
   | `neg-frontend-copy` | Fix a typo in the login page heading | `none` |
   | `boundary-schema` | Add a field to `packages/schemas` | `contract-first-change` |

   `expected` is a **set** of acceptable answers (usually one; `none` means "no skill
   should fire" — this catches *over*-triggering). Rationale for `compose-field-web` →
   `contract-first-change`: a web-read field is a cross-service contract; `add-api-endpoint`
   step 2 explicitly routes such changes to `contract-first-change`.

2. **`.claude/commands/eval-skills.md`** (`/eval-skills`) — the runner. On invocation Claude:
   1. Reads the `description` frontmatter from **every** `.claude/skills/*/SKILL.md` at
      runtime (so the eval never goes stale as skills change).
   2. For each scenario, dispatches a **fresh judge subagent** given *only* the task text +
      the set of `{skill-name: description}` pairs — **blind to the expected answer** — asking
      "which single skill should fire, or `none`?" with a one-line reason.
   3. Compares each verdict to the scenario's `expected` set; prints a **scorecard**:
      per-scenario PASS/FAIL + the judge's pick + reason, then a summary (`X/Y passed`).
   4. On any FAIL, names the likely cause (e.g. "descriptions of A and B both match task").

3. **`.claude/eval/README.md`** — purpose + honest limitations: proxy (not the real trigger
   engine); probabilistic (LLM judge); opt-in, human-reviewed, deliberately not in
   `make verify`; skills-only v1; how to add a scenario.

## Data flow

`/eval-skills` → read live skill descriptions → N fresh **blind** judge subagents (isolated
context, run concurrently) → compare each verdict vs `scenarios.md` `expected` → scorecard
to the human.

## Design rationale (isolation & honesty)

- **Fresh blind judges** = independence: no bias from the expected answer or from
  conversation history; each scenario judged on the descriptions alone.
- **Descriptions read live**, not hardcoded = the eval measures the *current* skills and
  cannot silently drift out of sync.
- **Report, not gate** = honest about a probabilistic proxy. The human reads the scorecard
  and decides; the eval never blocks a commit on an LLM's mood.

## Verification of this build

Run `/eval-skills` against the current two skills. The six seeded scenarios should score as
`expected`. The real signals are `compose-field-web`, `boundary-schema`, and the two `neg-*`
cases: if the two descriptions are discriminative, positives map correctly and negatives
return `none`. A FAIL here is a *finding about the skill descriptions*, not a bug in the
eval — which is exactly the value.

## Open questions

None. Mechanism (judge-on-descriptions proxy), runner form (command, not `make`),
`compose-field-web` expectation (`contract-first-change`), and skills-only scope were
confirmed during brainstorming.

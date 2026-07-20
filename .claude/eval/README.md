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

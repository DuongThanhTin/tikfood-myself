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

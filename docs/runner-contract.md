# Runner Contract

The canonical API contract is in `packages/schemas/openapi.yaml`.

JSON Schemas:

- `packages/schemas/feature-request.schema.json`
- `packages/schemas/runner-success-response.schema.json`
- `packages/schemas/runner-failure-response.schema.json`

## Endpoint

```text
POST /jobs/feature
```

## Request

Required fields:

- `feature_id`
- `repo`
- `base_branch`
- `title`
- `description`
- `acceptance_criteria`
- `mode`

## Success Response

Required fields:

- `success`
- `feature_id`
- `repo`
- `branch`
- `commit_sha`
- `summary`
- `files_changed`
- `tests_run`
- `test_result`
- `risks`
- `pr_title`
- `pr_body`

## Failure Response

Required fields:

- `success`
- `feature_id`
- `stage`
- `error`
- `logs`
- `recommendation`

## Runner Rules

- Validate input before doing repository work.
- Create branches only under `ai/`.
- Never push to `main` or `master`.
- Never force push.
- Never read secrets.
- Return JSON matching the OpenAPI contract.
- Do not claim success unless a branch was safely created, verified, committed, and pushed.

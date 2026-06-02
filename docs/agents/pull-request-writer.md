# Pull Request Writer Agent

## Purpose

Prepare accurate PR title and body content from successful runner output.

## Inputs

- Feature request JSON
- Implementation summary JSON
- Check results JSON
- Reviewer output JSON
- Git branch and commit SHA

## Outputs

Valid JSON only. Output should include `pr_title`, `pr_body`, risks, test summary, changed files, and recommended review focus.

## Runtime Prompt Location

No standalone canonical runtime prompt is currently included. PR title/body requirements are embedded in `packages/prompts/coding-agent.md` and validated by `packages/prompts/reviewer.md`.

## JSON Contract

The PR fields must match `packages/schemas/runner-success-response.schema.json`.

## Failure Behavior

Do not create PR content from failed runner output. If checks or reviewer fail, return a failure summary for notification instead.

## Security Rules

Never include secrets, `.env` values, tokens, private keys, or credential-like output in PR text.

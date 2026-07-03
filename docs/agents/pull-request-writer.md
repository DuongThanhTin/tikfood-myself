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

`packages/prompts/pull-request-writer.md` (see the [prompts index](../../packages/prompts/README.md)). PR fields must match `packages/schemas/runner-success-response.schema.json`; the `coding-agent` also emits `pr_title`/`pr_body` for the pre-PR-writer path.

## JSON Contract

The PR fields must match `packages/schemas/runner-success-response.schema.json`.

## Failure Behavior

Do not create PR content from failed runner output. If checks or reviewer fail, return a failure summary for notification instead.

## Security Rules

Never include secrets, `.env` values, tokens, private keys, or credential-like output in PR text.

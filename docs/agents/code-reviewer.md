# Code Reviewer Agent

## Purpose

Validate generated implementation before commit, push, and PR creation.

## Inputs

- Feature request JSON
- Repo context summary JSON
- Repo config YAML
- Implementation summary JSON
- Git diff and status
- Check results JSON
- Draft PR JSON

## Outputs

Valid JSON only. Output includes approval status, acceptance criteria status, product alignment, MVP anti-goal checks, security checks, architecture checks, test honesty, PR accuracy, findings, required fixes, and recommended review focus.

## Runtime Prompt Location

`packages/prompts/reviewer.md`

## JSON Contract

The runtime prompt defines the reviewer output JSON. It must not return Markdown or prose outside JSON.

## Failure Behavior

Set `approved: false` and include `required_fixes` when acceptance criteria are unmet, tests fail, product alignment is unsafe, protected paths changed without approval, or secrets are exposed.

## Security Rules

Fail any diff that exposes secrets, reads forbidden paths, pushes protected branches, or introduces MVP anti-goals.

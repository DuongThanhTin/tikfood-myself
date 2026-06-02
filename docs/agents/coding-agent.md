# Coding Agent

## Purpose

Implement scoped TikFood feature requests after the context reader has summarized the repository.

## Inputs

- Feature request JSON
- Repo context summary JSON
- Repo config YAML
- Available tools
- Git status

## Outputs

Valid JSON only. Output includes product alignment, docs read, implementation plan, files changed, tests, commands run, verification results, risks, assumptions, approval needs, commit message, PR title, and PR body.

## Runtime Prompt Location

`packages/prompts/coding-agent.md`

## JSON Contract

The runtime prompt defines the implementation output JSON. It must not return Markdown or prose outside JSON.

## Failure Behavior

Return `success: false`, list blocked reasons, leave PR fields empty, and set `requires_human_approval: true` when protected paths, anti-goals, or unsafe changes are involved.

## Security Rules

Never read secrets. Never push `main` or `master`. Never merge. Modify the smallest number of files possible.

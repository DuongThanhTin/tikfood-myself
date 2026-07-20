# .claude/ — Executable AI Operating System layer

Native Claude Code artifacts that make the repo's documentation-based AI Operating
System *self-triggering*. Each artifact is a thin wrapper over a `docs/` source — see
`.claude/CONVENTIONS.md` for the authoring contract and the doc->artifact mapping.

## Contents (Batch 1 — spine)

- `skills/add-api-endpoint/` — wraps `docs/recipes/add-api-endpoint.md`.
- `skills/contract-first-change/` — net-new; guards the `{data,error}` contract.
- `agents/repo-context-reader.md` — runs the context-loading protocol, returns a context pack.
- `commands/feature.md` — `/feature`, wraps `docs/workflows/feature.md`.

## Verify

`bash .claude/check-artifacts.sh` — or `make verify-claude`.

## Roadmap

Batch 1 proves the pattern. Later batches fan out the remaining recipes/agents/workflows
and author more net-new skills (`add-error-code`, `backend-service-change`). One spec per
batch under `docs/superpowers/specs/`.

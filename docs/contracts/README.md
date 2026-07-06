# Contracts — one place for TikFood's wire contracts

> **Why this directory exists:** the API response contract and the automation-runner
> job contract used to live apart (`docs/standards/api-contracts.md` and
> `docs/runner-contract.md`), so "what shape goes over the wire" had no single index.
> This directory is that index. **It links, it does not copy** — each contract keeps
> one owning document, and the machine-readable schemas stay the source of truth.

## Contracts

| Contract | Doc | Scope |
| --- | --- | --- |
| **Product API** | [`api.md`](api.md) | `{ data, error }` envelope, status codes, versioning, frontend type sync for `/api/v1/...` |
| **Automation runner** | [`runner.md`](runner.md) | `POST /jobs/feature` request / success / failure shapes + runner rules |

## Source of truth

- **Machine-readable schemas win over prose.** The runner contract is defined by
  [`packages/schemas/`](../../packages/schemas/) (`openapi.yaml` +
  `feature-request` / `runner-success-response` / `runner-failure-response`
  JSON Schemas). `packages/schemas/**` is a **protected path** — changes need human approval.
- The API response **envelope** decision is recorded in
  [ADR-0003 — data/error response envelope](../adr/0003-data-error-response-envelope.md).
- If a doc here disagrees with the code or a schema, the **code / schema wins** — flag
  the mismatch (per [`/CLAUDE.md`](../../CLAUDE.md) → *Source of Truth*).

# Service Contracts

> **Why this directory exists:** the repo documents API shapes
> ([`docs/contracts/api.md`](../contracts/api.md)) and the runner
> job ([`docs/contracts/runner.md`](../contracts/runner.md)), but there was no per-service
> page answering, at a glance: *what does this service expose, what does it depend on,
> who consumes it, and what are the rules for changing it?* These contracts give each
> runnable service one such page so cross-service changes are safe and intentional.
>
> They **link** to the owning standards/`CLAUDE.md`/ADRs rather than restating rules.
> Code wins on conflict — flag the mismatch.

**Related:** [`docs/architecture.md`](../architecture.md) ·
[`docs/REPOSITORY-MAP.md`](../REPOSITORY-MAP.md) · [`docs/adr/`](../adr/).

| Service | Contract | Kind |
| --- | --- | --- |
| `apps/api` | [`api.md`](api.md) | HTTP discovery API |
| `apps/web` | [`web.md`](web.md) | Next.js frontend (API consumer) |
| `apps/ai-code-runner` | [`ai-code-runner.md`](ai-code-runner.md) | Automation runner (MVP skeleton) |

Each contract has the same sections: **Purpose · Interface · Dependencies · Consumers ·
Change rules · Current state**.

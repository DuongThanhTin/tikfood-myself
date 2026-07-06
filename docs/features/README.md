# Feature Design Packages

> **Why this directory exists:** before a non-trivial feature is implemented, its design
> is worked out as a set of short, single-topic documents — so the approach is reviewable
> up front and a later session (human or agent) can build it without re-deriving
> decisions. Each feature gets **one subdirectory** holding that package. This is the
> design layer between Discovery and the implementation plan.

**Related:** [`docs/workflows/feature.md`](../workflows/feature.md) (the request → PR flow) ·
[`docs/thinking/README.md`](../thinking/README.md) (classify / reason) ·
[`docs/superpowers/`](../superpowers/) (specs + plans) ·
[`docs/verification/definition-of-done.md`](../verification/definition-of-done.md) (the "done" gate).

## Convention

- One subdirectory per feature: `docs/features/<feature-name>/`.
- **Design only — no code.** Each document distinguishes **Existing / Proposed / Open
  questions** and links to the owning repo docs instead of restating them.
- Language: English, per [ADR-0008](../adr/0008-documentation-language-policy.md).
- A design package precedes the implementation plan; on approval the plan/tasks drive
  the build (see the feature's `tasks.md` and [`docs/superpowers/plans/`](../superpowers/plans/)).

## Standard document set (use the subset a feature needs)

`overview` · `specification` · `architecture` · `api` · `database` · `security` ·
`frontend` · `ui` · `component-reuse` · `testing` · `release` · `tasks`

## Index

| Feature | Status | Entry |
| --- | --- | --- |
| Authentication | Design — not yet implemented | [`authentication/overview.md`](authentication/overview.md) |

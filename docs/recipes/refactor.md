# Recipe — Refactor (no behavior change)

**Goal:** improve structure/readability/size while keeping behavior identical.

**When to use:** a file/function carries multiple responsibilities, duplication has
accumulated, or an approved split is due (e.g. `DiscoveryExperience.tsx`).

## Context to load
- Owner doc for the area (`apps/*/CLAUDE.md`) and its "Future Direction" section — some
  refactors are pre-approved there.
- [`docs/thinking/README.md`](../thinking/README.md) §3 (is the refactor justified yet?).

## Process skill
`writing-plans` for anything non-trivial; keep the plan's steps small and independently
verifiable.

## Steps
1. **Justify** — is the complexity real *now*? Don't split/abstract speculatively
   ([`thinking/README.md`](../thinking/README.md)).
2. **Establish a safety net** — ensure tests cover the current behavior; add
   characterization tests first if coverage is thin.
3. **Refactor in small steps** — one responsibility moved at a time; keep each step
   green. Preserve public contracts (envelope, `lib/api.ts` types, exported signatures).
4. **No behavior change** — if you find a bug mid-refactor, note it and fix it separately
   ([fix-bug](fix-bug.md)).
5. **Match surrounding style** — naming, idioms, comment density.

## Verification
Full test suite passes unchanged before and after; `make verify-api` or
`make verify-web` clean. Diff shows movement, not behavior change.

## Exit
Behavior identical, tests green throughout, contracts intact, structure genuinely
improved.

## Common mistakes
- Mixing a refactor with a feature or bug fix (hard to review, hides regressions).
- Refactoring without a test safety net.
- Introducing an abstraction before it's justified (e.g. `AppError`, a state library).
- Breaking a public contract "while in there".

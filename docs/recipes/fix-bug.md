# Recipe — Fix a bug

**Goal:** find the real root cause and fix it with a regression test — not a symptom
patch.

**When to use:** something is broken, fails, or behaves incorrectly.

## Context to load
- The failing area's owner doc (`apps/*/CLAUDE.md`, [`docs/services/`](../services/)).
- For search/geo/dish bugs: ADR [0004](../adr/0004-in-memory-fallback-repository.md) —
  the bug may be a **known fallback-vs-Postgres divergence**, not a defect.

## Process skill
**`systematic-debugging` first** — do not propose a fix before the root cause is found.

## Steps
1. **Reproduce** — get a failing test or exact repro steps. If you can't reproduce, keep
   investigating; don't guess-fix.
2. **Localize** — narrow to the layer (handler/service/repository/frontend) and the
   specific input. Check whether it's the Postgres or fallback path.
3. **Root cause** — explain *why* it happens, not just where. Confirm against the code.
4. **Write a failing test** that captures the bug (table-driven case for the exact
   input).
5. **Fix** at the root, smallest change; keep layering intact.
6. **Confirm** the new test passes and no others regress.
7. **Check siblings** — the same class of bug elsewhere (e.g. alias logic exists in
   several places — see `IMPROVEMENTS.md`).

## Verification
`go test ./...` (or `web:typecheck`+`web:build`); the new regression test fails before /
passes after. See [Definition of Done](../verification/definition-of-done.md).

## Exit
Root cause stated, regression test added, fix minimal, no new failures, related
occurrences noted or fixed.

## Common mistakes
- Patching the symptom without understanding the cause.
- "Fixing" a documented fallback divergence that's actually expected behavior.
- No regression test (the bug can silently return).
- Scope creep — refactoring while fixing (do that separately, see
  [refactor](refactor.md)).

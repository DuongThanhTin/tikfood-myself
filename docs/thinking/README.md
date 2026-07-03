# Thinking Framework

> **Why this doc exists:** many standards in this repo deliberately defer to judgment —
> "add pagination only when an endpoint can realistically return large data", "don't
> introduce the error abstraction until enough error paths justify it", "add a library
> only when a real need appears". A senior engineer resolves these calls consistently; a
> new contributor or an AI agent needs the *heuristics* to do the same. This doc makes
> that reasoning explicit. It is the reasoning companion to the handbook's 3-layer model
> and skill-orchestration principle (`docs/handbook/01`).

**Related:** [`docs/handbook/01-ai-operating-system.md`](../handbook/01-ai-operating-system.md)
(3-layer model) · [`docs/ai/CONTEXT-LOADING.md`](../ai/CONTEXT-LOADING.md) ·
[`docs/recipes/`](../recipes/) (step-by-step task guides, added in a later phase) ·
[`docs/verification/definition-of-done.md`](../verification/definition-of-done.md).

## 1. Classify the task first

Before loading deep context or writing code, name the task type — it selects the process
skill, the docs to read, and the recipe:

| Signal | Task type | Start with |
| --- | --- | --- |
| "Add / build / support X" | Feature | `brainstorming` → `writing-plans` |
| "It's broken / wrong / fails" | Bug | `systematic-debugging` |
| "Make it cleaner / faster / smaller" | Refactor / perf | plan + [`performance.md`](../standards/performance.md) if perf |
| "Is it safe / is X exposed" | Security review | [`application-security.md`](../standards/application-security.md) |
| "How does X work / where is Y" | Research | read-only exploration (delegate if broad) |
| Touches auth/migration/infra/external/anti-goal | **Gated** | **stop → human approval** |

## 2. Process before implementation

Per handbook §1.5: **load the process skill first**, then implementation. Build →
`brainstorming`; bug → `systematic-debugging`; before "done" →
`verification-before-completion`. Hard skills (TDD, systematic-debugging) are followed
exactly; pattern skills adapt.

## 3. Reading "when justified" standards

Several standards say "not until justified." Resolve them with these triggers — build
the thing when a trigger fires, not before:

| Deferred capability | Build it when… |
| --- | --- |
| Pagination `meta` | an endpoint can realistically return large/unbounded data |
| Error abstraction (`AppError`) | there are enough distinct error paths to remove real duplication |
| React Query / Zod / form lib | real caching, mutations, or non-trivial forms actually appear |
| CSS framework | plain CSS demonstrably can't sustain the UI's needs ([ADR-0005](../adr/0005-plain-css-frontend.md)) |
| Cache / new index | a hot, stable query is measured (not assumed) |
| Splitting a big file | it carries multiple real responsibilities and change is friction |

Default when no trigger has fired: **do the smallest thing that satisfies the current
requirement**, and note the future option rather than building it.

## 4. Spike vs. build

- **Spike (throwaway exploration)** when the approach is unknown or risk is high — learn,
  then discard and implement cleanly with tests.
- **Build directly** when the pattern already exists in the codebase (reuse it) and the
  change is well understood.
- Prefer **reuse over new**: search for an existing helper/service/component before
  adding one (per `CLAUDE.md` → *Before Finishing*).

## 5. Scope & risk heuristics

- **Smallest correct change.** Match surrounding code's idioms, naming, and comment
  density. No drive-by rewrites.
- **One-way vs. two-way doors.** Reversible change → proceed with tests. Hard-to-reverse
  (public contract, schema, data shape, external call) → slow down; consider an ADR and
  human approval.
- **Blast radius.** Before changing a shared contract (the `{data,error}` envelope,
  `lib/api.ts` types, `packages/schemas/*`), enumerate consumers (see
  [`docs/services/`](../services/)) and update them together.
- **Two divergent paths.** For search/geo/dish changes, reason about **both** the
  Postgres and in-memory fallback paths ([ADR-0004](../adr/0004-in-memory-fallback-repository.md)).

## 6. When to stop and ask

Stop and request human approval when the task enters a gated area (auth, migrations,
infra, external network calls, protected paths, any anti-goal), when a decision is
hard to reverse and you're unsure, or when the requirement itself is ambiguous. Asking
early is cheaper than an unwanted one-way door.

## 7. Decide if it's ADR-worthy

If the task made a significant, constraining, or surprising choice, write an
[ADR](../adr/). If not, no ADR. This is the retrospective question at the end of the
core workflow.

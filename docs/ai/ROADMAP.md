# AI Operating System — Implementation Roadmap

> **Why this file exists:** it is the tracking artifact for turning
> `tikfood-ai-automation-starter` into an AI Operating System — a repository where
> Claude Code (and the `ai-code-runner`) can work like a senior engineer because the
> rules, context, decisions, standards, and task recipes are all explicit and linked.
> It records the audit that motivated the work, the phased plan, and progress.

Status: **Complete** (2026-07-03) — all 8 phases generated. Owner: Repository Architect role.

> **Scope:** the build plan for the AI Operating System **documentation**. Sibling
> planning docs: [`docs/roadmap.md`](../roadmap.md) = automation build phases ·
> [`IMPROVEMENTS.md`](../../IMPROVEMENTS.md) = prioritized code/infra findings.

## Guiding principles

1. **Index/link, never duplicate.** Layer-1 rules stay canonical in `CLAUDE.md`,
   `.ai-agent.yaml`, `docs/standards/**`, `docs/security.md`. New docs *point to*
   them. When docs conflict with code, **code wins** — flag the mismatch.
2. **Reusable reference docs over large prompts.** Prefer a small linkable doc a
   prompt can cite over embedding knowledge in a mega-prompt.
3. **Every document explains why it exists** and what it links to.
4. **Honest about current-vs-future state.** Skeletons and roadmap items are marked;
   no fake "production-ready" claims.
5. **English for reference/standard/contract docs** (matches `docs/standards/`);
   handbook additions stay bilingual (VN prose + EN technical) to match Parts 1–2.

## Audit summary

The repo is already AI-first with a strong Layer-1 (rules) and deep backend
standards. The gaps are cross-cutting AI-OS capabilities that the handbook itself
already forecasts (`docs/handbook/02` §2.5 names `recipes/`, `playbooks/`,
`verification/`, `thinking/`, `adr/` as future dirs; Parts 3–8 are skeleton/
placeholder). This roadmap fulfills that plan.

Gap register (priority: P0 foundation · P1 high leverage · P2 later):

| # | Capability | State | Pri |
|---|---|---|---|
| 1 | Repository Map | scattered across 3+ docs | P0 |
| 2 | AI Contract (canonical rules index) | distributed, no single entry | P0 |
| 3 | Context-Loading Protocol | default order only; Part 6 placeholder | P0 |
| 4 | ADR process + decision memory | absent (decisions implicit) | P1 |
| 5 | Architecture refresh + service contracts | `architecture.md` partly aspirational | P1 |
| 6 | Testing Standard (general/frontend) | backend-only; web has 0 tests | P1 |
| 7 | Application Security Standard | product security scattered | P1 |
| 8 | Performance Standard | only `backend/scaling.md` | P1 |
| 9 | Prompt Standard | implicit convention, undocumented | P1 |
| 10 | Thinking Framework | principle only; heuristics missing | P1 |
| 11 | Verification + Review guides | skill/prompt exist; no DoD/rubric doc | P1 |
| 12 | Recipes / Domain Playbooks | handbook Parts 3–4 skeleton | P1 |
| 13 | Prompt library completion | 2 role docs lack runtime prompts | P2 |
| 14 | Product domain model / glossary | vision exists; entity model missing | P2 |
| 15 | Handbook Parts 3,5,6,7,8 | skeleton/placeholder | P2 |

**Out of scope (routed to `IMPROVEMENTS.md`, needs human/code action):** CI
(`.github/workflows/` — protected path), frontend test infra + seed tests,
location-alias consolidation, fabricated UI data, worker infrastructure.

## Phases

Each phase lists purpose, documents, dependencies, complexity, impact.

### Phase 1 — Spine & Navigation  ·  P0  ·  Low  ·  Impact: High
Purpose: make every session's context-loading deterministic; give the rest of the
OS a single spine to link back to.
- `docs/REPOSITORY-MAP.md` — every directory → purpose, owner CLAUDE.md, read-before-editing.
- `docs/ai/AI-CONTRACT.md` — canonical index of rules-of-engagement (links, no copies).
- `docs/ai/CONTEXT-LOADING.md` — what to read per task type; when to drop context / spawn sub-agent.
Depends on: none.

### Phase 2 — Decision & Architecture Memory  ·  P1  ·  Medium  ·  Impact: High
Purpose: stop re-deriving or contradicting settled design.
- `docs/adr/README.md` (process) + `docs/adr/TEMPLATE.md` (MADR-lite).
- Seed ADRs: Gin+pgx (not GORM), plain CSS, `{data,error}` envelope, handler→service→repo
  layering, in-memory fallback repository, bilingual handbook.
- `docs/architecture.md` refresh — mark current-vs-roadmap; reconcile mismatches.
Depends on: Phase 1.

### Phase 3 — Service Contracts  ·  P1  ·  Medium  ·  Impact: Med-High
Purpose: make cross-service change safe.
- `docs/services/api.md`, `docs/services/web.md`, `docs/services/ai-code-runner.md`
  (inputs/outputs/dependencies/current endpoints/change rules), linking
  `docs/standards/api-contracts.md` and `docs/runner-contract.md`.
Depends on: Phase 2.

### Phase 4 — Standards Completion  ·  P1  ·  Medium  ·  Impact: High
- `docs/standards/testing.md` (general + frontend + coverage philosophy).
- `docs/standards/application-security.md` (authz, validation, rate-limit, CORS, PII,
  product prompt-injection) — complements `docs/security.md` (runner security).
- `docs/standards/performance.md` (latency/query/index/frontend budgets).
- `docs/standards/prompt-engineering.md` (the Prompt Standard/template).
Depends on: Phase 1.

### Phase 5 — Thinking / Verification / Review  ·  P1  ·  Medium  ·  Impact: High
- `docs/thinking/README.md` — task classification, when to brainstorm/research/spike,
  how to read "when justified" standards, risk heuristics.
- `docs/verification/definition-of-done.md` — DoD per surface (api/web/prompt).
- `docs/verification/review-guide.md` — review checklist + Critical/Major/Minor/Nit rubric.
Depends on: Phases 1, 4.

### Phase 6 — Recipes (Domain Playbooks)  ·  P1  ·  Med-High  ·  Impact: High
- `docs/recipes/README.md` + recipes: add-api-endpoint, fix-bug, refactor,
  add-runtime-prompt, frontend-component, add-migration (approval-gated),
  performance-investigation, security-review. Each: goal/inputs/steps/context/verification/exit.
Depends on: Phases 3, 4, 5.

### Phase 7 — Prompt Library & Product Knowledge  ·  P2  ·  Medium  ·  Impact: Medium
- `packages/prompts/README.md` (index) + `packages/prompts/feature-analyzer.md`
  + `packages/prompts/pull-request-writer.md`.
- `docs/tikfood/domain-model.md` (Venue, Dish, TrendScore, AISummary, Location/alias) + glossary.
Depends on: Phase 4.

### Phase 8 — Handbook Integration & Closeout  ·  P2  ·  Low-Med  ·  Impact: High
- Fill handbook Parts 3,5,6,7,8 as thin explainers that LINK to the new artifacts.
- Update handbook `README.md` status table; update `.ai-agent.yaml` `docs_required`;
  update root `README.md` / `CLAUDE.md` pointers.
Depends on: all prior phases.

## Progress log

- 2026-07-03 — Audit complete; roadmap approved (dedicated dirs + handbook links;
  English reference docs; docs + runtime prompts in scope; CI/tests flagged).
- 2026-07-03 — Phase 1 **complete**: `docs/REPOSITORY-MAP.md`, `docs/ai/AI-CONTRACT.md`,
  `docs/ai/CONTEXT-LOADING.md`.
- 2026-07-03 — Phase 2 **complete**: `docs/adr/README.md` + `TEMPLATE.md` + seed ADRs
  0001–0006; `docs/architecture.md` refreshed (current-vs-roadmap, ADR links).
- 2026-07-03 — Phase 3 **complete**: `docs/services/README.md` + `api.md` + `web.md`
  + `ai-code-runner.md` (verified routes/params against code).
- 2026-07-03 — Phase 4 **complete**: `docs/standards/testing.md`,
  `application-security.md`, `performance.md`, `prompt-engineering.md`; wired into
  `docs/standards/README.md`.
- 2026-07-03 — Phase 5 **complete**: `docs/thinking/README.md`;
  `docs/verification/definition-of-done.md` + `review-guide.md`.
- 2026-07-03 — Phase 6 **complete**: `docs/recipes/README.md` + 8 recipes
  (add-api-endpoint, frontend-component, fix-bug, refactor, add-runtime-prompt,
  add-migration, performance-investigation, security-review).
- 2026-07-03 — Phase 7 **complete**: `packages/prompts/README.md` +
  `feature-analyzer.md` + `pull-request-writer.md`; `docs/tikfood/domain-model.md`.
- 2026-07-03 — Phase 8 **complete**: handbook Parts 3–8 filled (explainer+links);
  handbook README status/index updated; `.ai-agent.yaml` `docs_required` extended; root
  `CLAUDE.md` "Start Here" + `README.md` AI-OS section added. Link check: 394 links, 0 broken.
- 2026-07-03 — **All phases complete.** Remaining handoffs (CI, frontend test harness,
  location-alias consolidation, fabricated-UI-data, worker infra) tracked in `IMPROVEMENTS.md`.

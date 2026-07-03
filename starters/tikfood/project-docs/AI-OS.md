# AI Operating System — Setup Guide (starter)

> **Why this file exists:** when you copy this starter into a real TikFood product repo,
> you want AI agents (Claude Code, `ai-code-runner`) to work like senior engineers from
> day one — with explicit rules, context, decisions, and reusable task recipes rather
> than ad-hoc mega-prompts. This is the checklist to stand up that "AI Operating System"
> in the new repo. The **reference implementation** lives in the automation workspace
> under `docs/` — copy and adapt those files; don't reinvent them.

## The model (3 layers)

An AI OS separates:

1. **Truth (what's allowed)** — `CLAUDE.md` (root + per-app), `.ai-agent.yaml`,
   `docs/standards/**`, product docs, security policy.
2. **Process (how to work)** — Superpowers skills + a handbook + task recipes.
3. **Runtime (who runs)** — Claude Code (interactive) and `ai-code-runner` (automated),
   driven by `packages/prompts/**`.

Keep a rule in exactly one layer: a rule change edits Truth; a workflow change edits
Process; don't bake rules into runtime prompts.

## Bring-up checklist

Copy from the reference implementation (paths are in the automation workspace) and adapt
names to the new repo:

- [ ] **Spine** — `docs/ai/AI-CONTRACT.md` (rules index), `docs/REPOSITORY-MAP.md`
      (dir → purpose), `docs/ai/CONTEXT-LOADING.md` (what to read per task).
- [ ] **Root `CLAUDE.md`** with a "Start Here" pointer to the three spine docs, plus
      MVP anti-goals, security/git rules, naming. Add per-app `CLAUDE.md` describing each
      app "as built today".
- [ ] **`.ai-agent.yaml`** — protected paths, allowed/blocked commands,
      `docs_required.before_coding` (include the spine docs). Seed provided in
      `starters/tikfood/.ai-agent.yaml`.
- [ ] **Decision memory** — `docs/adr/README.md` + `docs/adr/TEMPLATE.md`; write an ADR
      per significant choice (framework, DB, API envelope, layering).
- [ ] **Standards** — start from `docs/standards/` (backend/frontend/api/testing/
      application-security/performance/prompt-engineering).
- [ ] **Service contracts** — `docs/services/<service>.md` per runnable service.
- [ ] **Judgment & gates** — `docs/thinking/README.md`,
      `docs/verification/definition-of-done.md`, `docs/verification/review-guide.md`.
- [ ] **Recipes** — `docs/recipes/` step-by-step playbooks for your common task types.
- [ ] **Runtime prompts** — `packages/prompts/` following a Prompt Standard
      (JSON-only, product context + anti-goals, one role per prompt); index in
      `packages/prompts/README.md`.
- [ ] **Handbook** — onboarding narrative that links (not duplicates) all of the above.

## Principles to preserve

- **Index/link, never duplicate** — reference docs point to the owning source; code wins
  on conflict, flag mismatches.
- **Reusable docs over large prompts.**
- **Honest current-state** — mark skeletons/roadmap; no fake "production-ready" claims.
- **Preserve TikFood anti-goals** in prompts, config, and product docs (discovery only;
  no delivery/cart/checkout/booking/chat/etc.).

See the automation workspace `docs/ai/ROADMAP.md` for the phased build order that
produced the reference implementation.

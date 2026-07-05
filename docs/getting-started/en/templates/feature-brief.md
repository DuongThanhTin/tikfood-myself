# Template: Feature Brief

> 🇻🇳 Tiếng Việt: [`../../templates/feature-brief.md`](../../templates/feature-brief.md)
>
> Fill in this short form before starting. Small task: keep it at the top of the PR. Large
> task: save as `docs/superpowers/specs/<date>-<feature>-design.md`. How to use:
> [`../new-feature.md`](../new-feature.md).

---

```md
# Feature: <short name>

- **Goal (1 line):** <what the user can do / what problem is solved>
- **Belongs to:** apps/web | apps/api | both | packages/…
- **Scope (what to do):**
  - <…>
- **Out-of-scope (what NOT to do):**
  - <…>
- **Acceptance (done when):**
  - [ ] <observable condition 1>
  - [ ] <condition 2>
- **Constraints / risks:**
  - Discovery-only; no anti-goals (cart/order/checkout/payment/booking/chat/follow/monetization/livestream).
  - Needs human approval? (auth / migration / infra / external network call) → <yes/no>
  - Reuse first; tokens not hardcoded values; data through lib/api.ts; Vietnamese copy.
- **Track:** UI → docs/design/ + playbook <…>; backend/domain → docs/recipes/<…>.
- **Branch:** ai/<feature>  (base: dev)
```

---

### Filled example

```md
# Feature: Clear-filters button in the venue rail empty-state

- **Goal:** When the rail is empty, show a message + a "Clear filters" button to get results back.
- **Belongs to:** apps/web
- **Scope:** add a reset button to the rail empty-state; call clearFilters + re-search.
- **Out-of-scope:** API changes; pagination; illustration.
- **Acceptance:**
  - [ ] Empty → shows `.emptyText` + button.
  - [ ] Click → clearFilters() runs, list reloads.
  - [ ] Works in light/dark and both breakpoints.
- **Constraints:** discovery-only; follow docs/design; no new component if reusable.
- **Track:** docs/design/patterns/empty-state.md + playbooks/update-component.md
- **Branch:** ai/rail-empty-state-reset (base: dev)
```

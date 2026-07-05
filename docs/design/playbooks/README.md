# Design Playbooks

Step-by-step guides for building and changing UI in `apps/web` **the TikFood way** —
reusing the design system rather than inventing. Every playbook assumes you've read
[`../00-overview.md`](../00-overview.md) and the relevant foundations (`03–09`).

## Workflow playbooks (actionable today)

- [create-page.md](create-page.md) — add a route/page
- [create-component.md](create-component.md) — add a reusable component
- [update-component.md](update-component.md) — change an existing component safely
- [review-ui.md](review-ui.md) — review a UI change before merge

## "Create X" playbooks

Some targets exist today (map, card); others are **Proposed** and each playbook links to
its spec:

- [create-card.md](create-card.md) *(exists)* · [create-map.md](create-map.md) *(exists)*
- [create-form.md](create-form.md) *(from existing inputs)*
- [create-modal.md](create-modal.md) *(Proposed — [spec](../components/modal.md))*
- [create-dialog.md](create-dialog.md) *(Proposed — confirm/prompt variant)*
- [create-table.md](create-table.md) *(Proposed — [spec](../components/table.md))*
- [create-dashboard.md](create-dashboard.md) *(Proposed — no dashboard exists; anti-goal caution)*

## Ground rules (all playbooks)

1. **Reuse first** — search `globals.css` and `components/` before adding anything (`apps/web/CLAUDE.md`).
2. **Tokens, not literals** — reference foundations (`03/04/06`); don't hardcode colors/sizes.
3. **Bilingual copy** — Vietnamese UI strings, English identifiers.
4. **Data through `lib/api.ts`** — no raw `fetch` in components.
5. **Accessibility** — labels on icon-only controls, native controls, `alt` text (`09`).
6. **Anti-goals** — no cart/order/checkout/payment/booking/chat/follow/monetization/livestream.
7. **Finish against** [`../../verification/definition-of-done.md`](../../verification/definition-of-done.md).

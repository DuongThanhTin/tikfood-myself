# Recipe — Add / change a frontend component (`apps/web`)

**Goal:** build or modify discovery UI consistently with the frontend conventions and the
API contract.

**When to use:** new component, UI change, or wiring new API data into the web app.

## Context to load
- [`apps/web/CLAUDE.md`](../../apps/web/CLAUDE.md) · [`docs/services/web.md`](../services/web.md)
- [`docs/standards/frontend-architecture.md`](../standards/frontend-architecture.md)
- ADR [0005](../adr/0005-plain-css-frontend.md) (plain CSS)
- For distinctive new UI: the `frontend-design` skill.

## Process skill
`brainstorming` for new UX → implement. Use `frontend-design` when the visual design
matters.

## Steps
1. **Scope check** — discovery-only; no cart/checkout/booking/chat/etc.
2. **Data via `lib/api.ts`** — never add raw `fetch` in a component. If new data is
   needed, add/extend a typed function there; keep types in sync with Go JSON tags
   (snake_case).
3. **Component** — descriptive name (`VenueMap`, `VenueDetail`, `VenueRailCard`); place in
   `components/`. Prefer server component for initial data, client component for
   interactivity (mirror `page.tsx` → `DiscoveryExperience` pattern).
4. **Styling** — reuse existing `globals.css` classes; inline `style` only for computed
   values (e.g. marker positions). No new CSS framework without justification.
5. **Accessibility** — labels on inputs, `aria-label` on icon-only buttons; Vietnamese UI
   copy is the norm.
6. **Heavy deps lazy** — keep map/optional libs lazy-imported; don't add React
   Query/Zod/form libs until a real need exists.
7. **Docs** — update [`docs/services/web.md`](../services/web.md) if the data contract or
   structure changes.

## Verification
`npm run web:typecheck` and `npm run web:build`. Unit tests if the harness exists (else
say none ran). See [Definition of Done → Frontend](../verification/definition-of-done.md).

## Exit
Typecheck + build pass, data flows through `lib/api.ts`, accessible, no fabricated data
shown as real, docs synced.

## Common mistakes
- Scattering `fetch` calls; letting web types drift from the API.
- Presenting placeholder ratings/reviews/photos as real (see `IMPROVEMENTS.md`).
- Adding a CSS framework or data library prematurely.
- Growing `DiscoveryExperience.tsx` further instead of extracting a focused component.

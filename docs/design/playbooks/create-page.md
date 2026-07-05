# Playbook: Create a Page

> Add a new route/page to `apps/web` (Next.js App Router). Today the app has one real
> route ([`../pages/home.md`](../pages/home.md)); most "pages" here will be **Proposed**
> promotions of existing panels (see [restaurant](../pages/restaurant.md),
> [search](../pages/search.md)).

## When to use

You need a deep-linkable/shareable route (not just an in-panel view).

## Prerequisites

- Read [`../00-overview.md`](../00-overview.md), [`../05-grid-system.md`](../05-grid-system.md),
  and the target's page doc under [`../pages/`](../pages).
- Confirm the route isn't an anti-goal surface (`CLAUDE.md`).

## Steps

1. **Decide server vs client.** Mirror `app/page.tsx`: an async **server component**
   fetches initial data via `lib/api.ts`; a **client** component owns interaction.
2. **Create the route** at `app/<segment>/page.tsx` (e.g. `app/venue/[slug]/page.tsx`).
3. **Fetch data** with an existing `lib/api.ts` function (`getDiscoveryVenues`,
   `fetchVenueDetail`, …). Respect the `{ data, error }` envelope and the
   `fallbackVenues` path. **Never** add raw `fetch` in the component.
4. **Compose from existing pieces** — the app shell/grid ([`05`](../05-grid-system.md)),
   [cards](../components/card.md), [inputs](../components/input.md), [map](../components/map.md),
   and the [patterns](../patterns) (search/filter/loading/empty/error). Don't rebuild them.
5. **URL state** (if it's a search-like page): mirror `VenueSearchParams` in the query
   string; reuse `toSearchParams`.
6. **Copy** in Vietnamese; keep metadata (`title`/`description`) in `layout`/`page`.
7. **Accessibility**: labelled landmark regions, focusable controls, `alt` text.

## Checklist

- [ ] Server fetch through `lib/api.ts`; fallback path works with no API.
- [ ] Reused shell/components/patterns; no duplicated styles or tokens.
- [ ] Vietnamese copy; English identifiers.
- [ ] Landmarks labelled; keyboard-navigable; images have `alt`.
- [ ] Not an anti-goal surface (or human approval obtained).
- [ ] Meets [`definition-of-done`](../../verification/definition-of-done.md).

## Related

- [pages/home.md](../pages/home.md) · [pages/restaurant.md](../pages/restaurant.md) · [pages/search.md](../pages/search.md) · [create-component.md](create-component.md)

## References

- `apps/web/CLAUDE.md` · `apps/web/app/page.tsx` · `apps/web/lib/api.ts`.

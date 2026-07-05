# Page: Home

> **Status: As-built.** This is the **only** real page/route in the app — the entire
> discovery experience. (`restaurant`, `dish`, `profile`, `search` are **Proposed** —
> they don't exist as routes today.)

## Purpose

The single discovery screen: search + filter venues in the left rail, browse them on the
map, and open a venue's detail in the right panel.

## Usage

- **Route:** `app/page.tsx` (App Router `/`). It is an **async server component** that
  fetches initial venues, then renders the client experience:

```tsx
// app/page.tsx
export default async function Home() {
  const venues = await getDiscoveryVenues();
  return <DiscoveryExperience initialVenues={venues} />;
}
```

- **Data:** `getDiscoveryVenues()` (`lib/api.ts:171-191`) fetches
  `/api/v1/discovery/venues` with `revalidate:30`, falling back to `fallbackVenues` when
  `NEXT_PUBLIC_API_URL` is unset or the request fails.
- **Composition:** the client `DiscoveryExperience` renders the three-column app shell
  ([`../05-grid-system.md`](../05-grid-system.md)):
  - **Left panel** — brand, nav tabs, hero search, chips, filters, venue rail, auth CTA.
  - **Map stage** — the MapLibre map ([`../components/map.md`](../components/map.md)).
  - **Right panel** — venue detail (mounts on selection; documented under the Proposed
    [`restaurant.md`](restaurant.md), since it's the venue view).

## Examples

```jsx
{/* The shell composition (DiscoveryExperience.tsx:398-658) */}
<main data-theme={theme} className="appShell …">
  <aside className="leftPanel" aria-label="Discovery controls">…</aside>
  <section className="mapStage" aria-label="Restaurant map">…</section>
  {selectedVenue ? <aside className="rightPanel" aria-label="Venue detail">…</aside> : null}
</main>
```

## Rules

- Initial venues are fetched **server-side** and passed as `initialVenues`; subsequent
  searches are client-side via `lib/api.ts`.
- The first venue is auto-selected on load (`selectedVenue = initialVenues[0] ?? null`, `:224`).
- The page holds no routing beyond this single screen (no nested routes).

## Do / Don't

- **Do** keep the server/client split: server fetches initial data, client owns interaction.
- **Do** rely on the `fallbackVenues` path so the page renders without an API.
- **Don't** add data fetching in the page beyond `getDiscoveryVenues`.
- **Don't** introduce client-side routing for detail; it's a panel, not a route (today).

## Accessibility

- The three regions are labelled landmarks (see [`../09-accessibility.md`](../09-accessibility.md)).
- All page-level gaps (focus, aria-live, inert buttons) are inherited from the components; none are page-specific.

## Related

- [Grid system](../05-grid-system.md) · [components/map.md](../components/map.md) · [components/card.md](../components/card.md) · [patterns/search.md](../patterns/search.md) · [patterns/filter.md](../patterns/filter.md)
- Proposed venue view: [pages/restaurant.md](restaurant.md).

## Files

- `apps/web/app/page.tsx` — the route.
- `apps/web/app/layout.tsx` — html/body, metadata, font & MapLibre CSS imports.
- `apps/web/components/DiscoveryExperience.tsx` — the experience `:222-660`.
- `apps/web/lib/api.ts` — `getDiscoveryVenues` `:171-191`.

## References

- `apps/web/CLAUDE.md` · [Design principles](../01-design-principles.md).

## Future improvements

- If detail/search/profile become shareable, promote them to real routes (see the Proposed page docs).

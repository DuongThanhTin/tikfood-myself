# Page: Restaurant (Venue)

> **Status: As-built as a panel; Proposed as a route.** The venue view **exists today**
> as the right-hand detail *panel* (`VenueDetail`), not a standalone page/route. This doc
> documents the real panel and proposes promoting it to a shareable route.

## Purpose

Show a selected venue in full: hero + badges, title + rating + quick actions, trend
score, AI/about copy, categories, social videos, dishes, and opening hours.

## Usage (as-built — the detail panel)

Rendered by `VenueDetail` (`DiscoveryExperience.tsx:1032-1178`) inside `.rightPanel`,
mounted only when a venue is selected (`:632-657`). Building blocks:

- **Close** `.closeButton` (`globals.css:1163-1177`).
- **Hero** `.detailHero` (320px → 240px ≤640px) + `.detailBadges` overlay ("HOT"/"TRENDING").
- **Title row** `.detailTitleRow` — `h2` (Plus Jakarta 30px), `.detailMeta` (address · district), `.compactDetailActions` (route/share/save via [`glassAction`](../components/button.md)), `.ratingBlock`.
- **Trend score** — the [Trend score card](../components/card.md) (`ai_summary` + clamped `trend_score`).
- **About** `.aboutText`; **categories** `.category` pills.
- **Videos** — the [Video card](../components/card.md) grid.
- **Dishes** — `.dishDetailItem` rows when `venue.dishes` exists, else `.dish` chips from `trending_dishes`.
- **Opening hours** — `.openingRow` rows (`formatDay` labels), when present.
- **Menu CTA** — `primaryButton large` (currently inert).

Detail data loads lazily on select if `dishes`/`opening_hours` are missing (`:255-285`).

## Examples

```jsx
{/* Mount + close (:632-641) */}
<aside className="rightPanel" aria-label="Venue detail">
  <button className="closeButton" type="button" aria-label="Close venue detail" onClick={() => setSelectedVenue(null)}>
    <Icon name="close" />
  </button>
  <VenueDetail venue={selectedVenue} /* … */ />
</aside>
```

## Rules

- **Real vs mocked:** `name`, `address`, `district`, `about`, `ai_summary`, `categories`,
  `dishes`, `opening_hours` are real API fields; `rating`, `reviews`, badges, and imagery
  come from the hardcoded `mediaByVenue` (mocked).
- Mount conditionally; don't CSS-hide.
- share/save/menu are inert today — wire or `disabled` before shipping.

## Do / Don't

- **Do** use the lazy-detail pattern and show `.detailStatus` while loading.
- **Do** prefer detailed dish rows when dish data exists.
- **Don't** present mocked rating/reviews/images as real.
- **Don't** rely on `font-weight:850` for compact actions (not loaded).

## Accessibility

- `aside[aria-label="Venue detail"]`; close button labelled; `<h3>` section headings structure the body.
- Gaps: inert enabled buttons; no `aria-live` for lazy load — see [`../09-accessibility.md`](../09-accessibility.md).

## Proposed: promote to a route

If venues need to be **shareable/deep-linkable**, add `app/venue/[slug]/page.tsx` that
server-fetches `fetchVenueDetail(slug)` and renders the same `VenueDetail` full-width,
with the map optional. This is Proposed — it does not exist. Keep the panel for in-context
browsing; add the route for sharing.

## Related

- [components/card.md](../components/card.md) · [components/button.md](../components/button.md) · [components/map.md](../components/map.md) · [pages/home.md](home.md) · [Proposed dish page](dish.md)

## Files

- `apps/web/app/globals.css` — detail styles `:1163-1471`.
- `apps/web/components/DiscoveryExperience.tsx` — `VenueDetail` `:1032-1178`; mount `:632-657`; lazy load `:255-285`.
- `apps/web/lib/api.ts` — `Venue`/`VenueDish`/`OpeningHour` `:1-61`; `fetchVenueDetail` `:213-227`.

## References

- [Design principles](../01-design-principles.md) · `apps/web/CLAUDE.md`.

## Future improvements

- Source rating/reviews/images from the API; wire share/save/menu; add the shareable route above.

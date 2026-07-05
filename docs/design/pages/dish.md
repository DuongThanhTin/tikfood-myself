# Page: Dish

> **Status: Proposed — not implemented.** There is no dish page or route today. Dishes
> appear only *inside* a venue: as `trending_dishes` chips and `.dishDetailItem` rows in
> the detail panel. This spec is a forward proposal grounded in the existing `VenueDish`
> data and current components.

## Purpose

A dish-first view — "where is this dish trending?" — matching the product's dish-first
principle ([`../01-design-principles.md`](../01-design-principles.md#2-dish--and-venue-first-content)).
Would show a dish, its social/trend signals, and the venues serving it.

## Data (already exists)

`VenueDish` (`lib/api.ts:39-54`) already carries everything a dish view needs:
`name`, `slug`, `short_description`, `about`, `category`, `cuisine`,
`price_min_vnd`/`price_max_vnd`, `mention_count`, `video_count`, `view_count`,
`trend_score`. **No dish endpoint/route consumes it standalone today.**

## Usage (proposed)

- Route `app/dish/[slug]/page.tsx` (Proposed), server-fetching a dish + its venues.
- Reuse existing components: [Trend score card](../components/card.md) for the dish trend,
  [Video card](../components/card.md) grid for social clips, [Venue rail card](../components/card.md)
  for venues serving it, [Map](../components/map.md) to plot them.
- Reuse `.category`/`.dish` pills and the `formatPrice`/`formatCompactNumber` helpers.

## Rules (proposed)

- Discovery-only: no ordering/cart for the dish (anti-goals).
- Reuse the venue-card and trend-card language; don't invent a new card.
- Prices via `formatPrice` (`k`/`tr`), views via `formatCompactNumber` (`K`/`M`).

## Do / Don't

- **Do** lean on `VenueDish` fields already defined.
- **Do** cross-link each venue to its [restaurant view](restaurant.md).
- **Don't** add commerce actions.
- **Don't** duplicate token/spacing definitions — reference the foundations.

## Accessibility (proposed)

- Same landmark/labeling and image-alt conventions as [home](home.md)/[restaurant](restaurant.md).

## Related

- [pages/restaurant.md](restaurant.md) · [components/card.md](../components/card.md) · [components/map.md](../components/map.md) · [patterns/search.md](../patterns/search.md)

## Files

- **None yet** (Proposed). Data analogue: `apps/web/lib/api.ts` `VenueDish` `:39-54`.

## References

- [Design principles](../01-design-principles.md) · `docs/tikfood/product-vision.md`.

## Future improvements

- N/A until implemented; needs a dish endpoint + route.

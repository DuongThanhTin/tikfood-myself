# Page: Search

> **Status: Proposed — not implemented as a page.** Search **behavior** exists today
> (inline in the left rail — see [`../patterns/search.md`](../patterns/search.md)), but
> there is no dedicated search *route* or full-page results view. This is a forward
> proposal.

## Purpose

A dedicated, deep-linkable search results view (e.g. `?q=…&district=…`) — useful for
sharing a search or landing from an external link. Today search is an in-panel action on
the single home page.

## Usage (proposed)

- Route `app/search/page.tsx` (Proposed) reading query/filter params from the URL and
  server-fetching results, then rendering the same rail + map as [home](home.md).
- Reuse [Input](../components/input.md) (search field, filters), the [search](../patterns/search.md)
  and [filter](../patterns/filter.md) patterns, [Venue rail card](../components/card.md),
  [Map](../components/map.md), and [empty-state](../patterns/empty-state.md).
- The key addition over home is **URL-driven state**: params ↔ query string are the source
  of truth (today they live only in React state).

## Rules (proposed)

- URL params mirror `VenueSearchParams` (`lib/api.ts:72-87`); reuse `toSearchParams`.
- Reuse the existing search/filter patterns exactly — this page is a routed presentation
  of them, not new behavior.
- Keep deliberate-search UX consistent with home (or intentionally decide on as-you-type).

## Do / Don't

- **Do** make filter state URL-addressable (the main reason to have this page).
- **Do** reuse home's rail/map composition.
- **Don't** fork the search logic — share it with the home experience.
- **Don't** duplicate token/pattern docs; reference them.

## Accessibility (proposed)

- Same landmark/labeling as home; announce result counts via `aria-live` (a gap noted for
  the current inline search too).

## Related

- [patterns/search.md](../patterns/search.md) · [patterns/filter.md](../patterns/filter.md) · [pages/home.md](home.md) · [components/input.md](../components/input.md) · [Proposed pagination](../patterns/pagination.md)

## Files

- **None yet** (Proposed). Behavior analogue: `DiscoveryExperience.tsx` `runSearch`/`params` `:242-301`; `lib/api.ts` `toSearchParams` `:229-241`.

## References

- `apps/web/CLAUDE.md`.

## Future improvements

- N/A until implemented; primarily an App Router + URL-state exercise reusing existing pieces.

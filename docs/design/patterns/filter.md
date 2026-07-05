# Pattern: Filter

> **Status: As-built.** How venue filtering is composed and applied.

## Purpose

Narrow discovery results by district, max price, tags, open-now, sort order, and
proximity. Filters combine into one request.

## Usage

- **Controls:** district & max-price selects, the tag bar (toggle tags), the open-now
  toggle, and the sort select ([`../components/input.md`](../components/input.md)); plus
  quick-prompt **chips** that set query/price/sort.
- **State → params:** each control updates its state; the `params` memo assembles
  `VenueSearchParams` (`DiscoveryExperience.tsx:242-253`).
- **Proximity override:** when the user shares location (`nearUser`), sort is forced to
  `distance`, the sort select is disabled, and `radius_m:3000` is applied (`:242-253, 520-524`).
- **Apply:** filters take effect on `runSearch`; **Reset** (`clearFilters`, `:303-318`)
  clears every filter, route, and error, then re-searches `sort:"trending"`.

## Examples

```jsx
{/* Tag toggle (:498-509) */}
<button className={activeTag === tag.value ? "tag active" : "tag"} type="button"
  onClick={() => setActiveTag(activeTag === tag.value ? "" : tag.value)}>{tag.label}</button>

{/* Sort — disabled under proximity mode (:520-524) */}
<select className="sortSelect" value={nearUser ? "distance" : sort} disabled={Boolean(nearUser)} aria-label="Sort venues">…</select>
```

## Rules

- One active tag at a time (`activeTag` string); re-tapping clears it.
- Proximity mode wins: it forces `distance` sort and disables the sort control.
- `clearFilters` resets to `sort:"trending", limit:20` and clears route/error state.
- Empty filter values serialize away (see [search](search.md)).
- Sort options map to API values `trending | videos | price | distance`.

## Do / Don't

- **Do** keep filter state as the single source; derive the request from the `params` memo.
- **Do** disable conflicting controls (sort during proximity) rather than silently overriding without feedback.
- **Don't** allow multiple active tags (the state is a single value).
- **Don't** convey a tag's selected state by color alone for a11y (add `aria-pressed`).

## Accessibility

- Selects/toggle are native and labelled; the sort select has `aria-label`.
- Tag "active" state is color-only (no `aria-pressed`) — a gap (see [`../09-accessibility.md`](../09-accessibility.md) A3).

## Related

- [components/input.md](../components/input.md) · [patterns/search.md](search.md) · [patterns/empty-state.md](empty-state.md)

## Files

- `apps/web/components/DiscoveryExperience.tsx` — controls `:474-532`; chips `:459-471`; `params` `:242-253`; `clearFilters` `:303-318`; `useCurrentLocation` `:320-352`.
- `apps/web/lib/api.ts` — `VenueSearchParams` `:72-87`.

## References

- `apps/web/CLAUDE.md`.

## Future improvements

- Support multi-select tags if product needs it (state + API change).
- Add `aria-pressed` to tags; surface active-filter chips/count for clarity.

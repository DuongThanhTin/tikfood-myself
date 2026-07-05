# Pattern: Pagination

> **Status: Proposed — not implemented.** There is no pagination today. The venue rail
> fetches a fixed `limit: 20` (`DiscoveryExperience.tsx:253`) with no way to load more.
> This spec is a forward proposal.

## Purpose

Let users browse beyond the first page of results. Today the list is capped at 20 and
silently truncated — a discovery ceiling.

## Current state (the gap)

- `params.limit = 20` is hardcoded; the API supports `limit` (`VenueSearchParams`,
  `lib/api.ts:87`) but no `offset`/cursor is sent and no "load more" UI exists.
- Result count is shown in the map legend, but there's no indication that results were
  capped.

## Usage (proposed)

- Prefer **"load more" / infinite scroll** over numbered pages — it fits the map-first,
  rail-scroll UX better than page numbers.
- Reuse [`button`](../components/button.md) `secondaryButton` for an explicit "Xem thêm",
  or observe the rail scroll position for infinite append.
- Requires an API contract for continuation (offset or cursor) — coordinate with
  `apps/api` before implementing.

## Examples (proposed shape)

```jsx
{hasMore ? (
  <button className="secondaryButton" type="button" onClick={loadMore} disabled={isLoadingMore}>
    {isLoadingMore ? "Đang tải" : "Xem thêm"}
  </button>
) : null}
```

## Rules (proposed)

- Append to the existing list; don't replace it (unlike a fresh search).
- Show a loading affordance while fetching the next page ([loading](loading.md)).
- Make truncation visible (e.g. the "Xem thêm" button) rather than silently capping.
- Keep the map markers in sync with the appended venues.

## Do / Don't

- **Do** align with the API's continuation contract before building.
- **Do** reuse loading/empty/error patterns.
- **Don't** add numbered pagination unless product needs deep-linkable pages.
- **Don't** silently cap results with no affordance (the current gap).

## Accessibility (proposed)

- "Load more" is a labelled button; announce newly loaded count via `aria-live`.
- Infinite scroll needs a keyboard-reachable trigger and a "results loaded" announcement.

## Related

- [patterns/loading.md](loading.md) · [patterns/search.md](search.md) · [components/card.md](../components/card.md)

## Files

- **None yet** (Proposed). Touch points: `DiscoveryExperience.tsx:253` (`limit`), `lib/api.ts:87`.

## References

- `apps/web/lib/api.ts` · `apps/api` discovery endpoint contract.

## Future improvements

- N/A until implemented; requires an API `offset`/cursor.

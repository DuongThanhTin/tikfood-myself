# Pattern: Search

> **Status: As-built.** How free-text discovery search works today.

## Purpose

Let users find venues by typing a natural-language-ish query (e.g. "date dưới 500k ở
Quận 1"). Search is deliberate (runs on an explicit button), not as-you-type.

## Usage

- **Input:** the hero search field ([`../components/input.md`](../components/input.md#usage)).
  Typing only updates `query` state (`DiscoveryExperience.tsx:451-455`).
- **Trigger:** the "Tìm kiếm" primary button calls `runSearch()` (`:535`); quick-prompt
  chips can pre-fill `query` via `applyChip` (`:354-364`).
- **Params:** `query` feeds `q` inside the `params` memo alongside the active filters
  (`:242-253`).
- **Fetch:** `runSearch` → `fetchDiscoveryVenues(params)` (`:287-301`), which hits
  `GET /api/v1/discovery/venues?…` when `NEXT_PUBLIC_API_URL` is set, else filters
  `fallbackVenues` client-side (`lib/api.ts:197-282`).
- **Result:** venue list replaces state; the first result is auto-selected (unless
  `selectFirst:false`).

## Examples

```jsx
{/* Deliberate search (:535-537) */}
<button className="primaryButton" type="button" onClick={() => void runSearch()} disabled={isLoading}>
  {isLoading ? "Đang tìm" : "Tìm kiếm"}
</button>
```
```ts
// Query serialization drops empty/false values (lib/api.ts:229-241)
const params = { q: query, district, tags: activeTag, max_price_vnd: Number(maxPrice) || undefined, sort, limit: 20 };
```

## Rules

- Search is explicit (button/chip), not keystroke-triggered.
- Always route data through `lib/api.ts` (respect the `{ data, error }` envelope); don't `fetch` in components.
- Empty/false params are omitted from the query string (`toSearchParams`).
- The client fallback search matches against name, description, about, district, categories, and trending dishes (`filterFallbackVenues`).

## Do / Don't

- **Do** keep the placeholder a realistic Vietnamese example query.
- **Do** reflect loading via the button label + `disabled`.
- **Don't** fire a request on every keystroke (current design is deliberate search).
- **Don't** bypass `lib/api.ts`.

## Accessibility

- The field is labelled via its `<label>` wrapper; the button has visible text.
- Loading is signalled by the label swap + `disabled` only (no `aria-busy`); results update with no `aria-live` — see [`../09-accessibility.md`](../09-accessibility.md) (A4/A5).

## Related

- [components/input.md](../components/input.md) · [patterns/filter.md](filter.md) · [patterns/loading.md](loading.md) · [patterns/empty-state.md](empty-state.md)

## Files

- `apps/web/components/DiscoveryExperience.tsx` — `query` `:228`; `params` `:242-253`; `runSearch` `:287-301`; `applyChip` `:354-364`.
- `apps/web/lib/api.ts` — `fetchDiscoveryVenues` `:197-211`; `toSearchParams` `:229-241`; `filterFallbackVenues` `:243-282`.

## References

- `apps/web/CLAUDE.md` (API access rules).

## Future improvements

- Consider debounced as-you-type search (would need loading/`aria-live` affordances).
- Parse the natural-language query (price/district) rather than sending it raw as `q`.

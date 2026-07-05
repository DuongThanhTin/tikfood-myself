# Form Field & Select

> As-built. Tokens live in [`../01-foundations.md`](../01-foundations.md).

## Purpose

The labelled dropdown pattern used for discovery filters (district, max price, sort).
A `field` is a `<label>` wrapper pairing a caption `<span>` with a native `<select>`.

## Usage

- `.field` (`globals.css:455-461`): `display:grid; gap:7px`, caption in `--on-muted`
  12px/800.
- `.field select, .sortSelect` (`:463-473`): `width:100%`, `min-height:42px`,
  `border-radius:10px`, `1px var(--outline)`, `padding:0 12px`, text `--on-surface`,
  `background:var(--surface-low)`, `outline:0`.
- `.filterGrid` (`:449-453`): two-column grid, `gap:10px`; collapses to one column
  under 640px (`:1566-1570`).
- `.sortSelect` (`:501-503`): `max-width:150px`, used standalone in the filter actions
  row rather than inside a `.field`.

Controls are **native `<select>`** elements — no custom dropdown.

## Examples

```jsx
{/* District field (:476-483) */}
<label className="field">
  <span>Quận</span>
  <select value={district} onChange={(e) => setDistrict(e.target.value)}>
    <option value="">Tất cả</option>
    <option value="District 1">Quận 1</option>
    <option value="District 3">Quận 3</option>
  </select>
</label>

{/* Sort — disabled while sorting by distance, with an explicit aria-label (:520-531) */}
<select
  className="sortSelect"
  value={nearUser ? "distance" : sort}
  onChange={(e) => setSort(e.target.value as VenueSearchParams["sort"])}
  disabled={Boolean(nearUser)}
  aria-label="Sort venues"
>
  <option value="trending">Trending</option>
  <option value="videos">Video nhiều nhất</option>
  <option value="price">Giá tốt</option>
  <option value="distance">Gần nhất</option>
</select>
```

## Rules

- Wrap a captioned select in `<label className="field">` so the caption labels the control implicitly.
- A standalone select (no visible caption, e.g. sort) must carry an explicit `aria-label`.
- Option `value`s map to API params (`district`, `max_price_vnd`, `sort`) — keep them in the shapes `VenueSearchParams` expects (`lib/api.ts:72-87`).
- The empty option ("Tất cả"/"Bất kỳ") represents "no filter" and serializes away in `toSearchParams` (`lib/api.ts:229-241`).
- When location sort is forced (`nearUser` set), the sort select is `disabled` and shows `distance`.

## Do

- Use `.filterGrid` to pair two related selects; it already handles the responsive collapse.
- Give the first option an empty value to mean "unset".
- Keep captions short and in Vietnamese.

## Don't

- Don't build a custom dropdown; the system uses native `<select>` for accessibility and simplicity.
- Don't omit `aria-label` on a caption-less select.
- Don't hardcode option values that the API doesn't accept.

## Accessibility

- Captioned selects rely on the wrapping `<label>` for their accessible name.
- The sort select, having no visible caption, sets `aria-label="Sort venues"`.
- Native `<select>` gives keyboard operability and OS-native option lists for free.
- Disabled state (sort during distance mode) uses the native `disabled` attribute, which assistive tech announces.

## Related components

- [Toggle](toggle.md) — the checkbox filter that shares the filter actions row.
- [Chip & Tag](chip-tag.md) — the tag filters above the grid.
- [Button](button.md) — the search/reset actions below the grid.

## Files

- `apps/web/app/globals.css` — `.field` `:455-461`; select `:463-473`; `.filterGrid` `:449-453`; `.sortSelect` `:501-503`; responsive `:1566-1570`.
- `apps/web/components/DiscoveryExperience.tsx` — `:474-532`.
- `apps/web/lib/api.ts` — `VenueSearchParams` `:72-87`; `toSearchParams` `:229-241`.

## References

- [Component catalog — Form controls](../03-components.md#form-controls)
- [Patterns — Data-flow](../04-patterns.md#data-flow-pattern)

## Future improvements

- No spacing/radius tokens back these controls (see [R6](../05-inconsistencies.md#r6--no-spacing-or-radius-scale-tokens)); a token pass would unify the 42px height / 10px radius with other inputs.

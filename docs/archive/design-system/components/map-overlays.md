# Map Overlays

> As-built. Tokens live in [`../01-foundations.md`](../01-foundations.md).

## Purpose

The two floating glass info cards layered over the map:
- **Count legend** (`mapCountLegend`) — how many venues are in the current area.
- **Route summary** (`routeSummary`) — distance/duration of an active route, with a
  clear button.

(The floating "Tìm quanh đây" bottom action and the topbar controls belong to
[Venue Map](venue-map.md).)

## Usage

- `.mapCountLegend` (`globals.css:1043-1073`): absolute, `z-index:5`, `right:24px`,
  `bottom:100px`; grid, `max-width:210px`; `1px var(--outline)`, radius 16,
  `padding:12px 14px`, `--glass-strong`, `--shadow`, `blur(18px)`. `strong` count in
  `--primary` (Plus Jakarta 24px); `span` muted 11px/800.
- `.routeSummary` (`:1075-1115`): absolute, `z-index:6`, `right:24px`, `bottom:176px`;
  grid, `min-width:210px; max-width:280px`; `1px rgba(47,128,255,0.35)` (blue) border,
  radius 18, `--glass-strong`, `blur(18px)`. `strong` distance in **blue `#2f80ff`**
  (Plus Jakarta 24px); `span` muted; `button` blue 12px/900.

## Examples

```jsx
{/* Count legend — rendered in both the fallback branch (:1010) and the live cluster branch (:1015-1020) */}
<div className="mapCountLegend mapCountLegendLive">
  <strong>{venues.length}</strong>
  <span>cửa hàng trong khu vực hiện tại</span>
</div>

{/* Route summary (:1021-1027) */}
{activeRoute ? (
  <div className="routeSummary">
    <strong>{formatDistance(activeRoute.distanceMeters)}</strong>
    <span>{formatDuration(activeRoute.durationSeconds)} tới {selectedVenue?.name ?? "địa điểm"}</span>
    <button type="button" onClick={onClearRoute}>Xóa chỉ đường</button>
  </div>
) : null}
```

## Rules

- Overlays are absolutely positioned within `.mapCanvasWrap`; their `z-index` (5/6) sits above markers (3–4) but below the detail close button (7).
- The count legend appears in two code paths (pre-map fallback and live cluster mode) with identical copy — keep them in sync (see [R19](../05-inconsistencies.md#secondary-findings)).
- The `mapCountLegendLive` modifier is applied but **has no CSS rule** (see [R2](../05-inconsistencies.md#r2--mapcountlegendlive-used-but-undefined)); the live legend inherits base styling only.
- The route summary is blue-themed to match the route line; that blue is a hardcoded literal, not a token (R5).
- Distance/duration are formatted by `formatDistance`/`formatDuration` (`:1295-1310`).

## Do

- Keep overlays inside the map wrap so they position against the map, not the page.
- Match the route summary color to the route line (both blue) for a coherent routing state.
- Show the route summary only when `activeRoute` exists.

## Don't

- Don't add a third duplicate of the count legend; consolidating the two renders would be better (R19).
- Don't rely on the `mapCountLegendLive` class for styling — it's a no-op (R2).
- Don't reintroduce the blue as a literal in new code; it wants a token (R5).

## Accessibility

- Overlays are text content, read inline by screen readers.
- The route summary's "Xóa chỉ đường" is a real `<button>` with visible text (accessible name present).
- The overlays are visual context; there's no live-region announcement when counts/route update (they change silently for AT users).

## Related components

- [Map Marker](map-marker.md) — what the count/route describe.
- [Venue Map](venue-map.md) — the host stage and the bottom action.
- [Venue Detail Panel](venue-detail-panel.md) — also triggers routing (mirrors the clear action).

## Files

- `apps/web/app/globals.css` — legend `:1043-1073`; route summary `:1075-1115`.
- `apps/web/components/DiscoveryExperience.tsx` — legend `:1010, 1015-1020`; route summary `:1021-1027`; formatters `:1295-1310`.

## References

- [Patterns — Map interaction / UI states](../04-patterns.md#map-interaction-patterns)
- [Component catalog — Cards (map overlays)](../03-components.md#cards)
- [Inconsistency register — R2, R5, R19](../05-inconsistencies.md)

## Future improvements

- Define `.mapCountLegendLive` or drop the modifier (R2).
- Merge the duplicated legend renders (R19).
- Announce route/count changes via an `aria-live` region.

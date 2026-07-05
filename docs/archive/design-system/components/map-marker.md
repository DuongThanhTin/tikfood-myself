# Map Marker

> As-built. Tokens live in [`../01-foundations.md`](../01-foundations.md).

## Purpose

The pins that represent venues (and the user's location) on the map. There are two
rendering systems that share styling: real MapLibre DOM markers and a CSS fallback
layer used before the map is ready (see [R8](../05-inconsistencies.md#r8--two-parallel-marker-systems)).

## Usage

Base `.mapMarker, .fallbackMarker` (`globals.css:874-882`): `2px solid var(--primary)`,
`background:rgba(255,183,127,0.16)`, fire-tinted shadow, 12px/900.

| Element | Class | Role |
|---------|-------|------|
| Cluster pin | `mapMarker cluster` (`:892-895`) | 34×34 pill with a count `strong` (zoom < 12.2) |
| Named marker | `mapMarker named` (`:897-906`) | borderless bubble + pin (zoom ≥ 12.2) |
| Bubble | `markerBubble` (`:908-948`) | pill: icon `em` (22×22) + ellipsised name (max 116px) |
| Pin | `markerPin` (`:950-967`) | 17×17 `--trend-fire` dot, 3px white border, white center |
| User location | `userLocationMarker` / `fallbackUserLocation` (`:995-1013`) | 24×24 **blue `#2f80ff`** dot |
| Fallback | `fallbackMarker` (`:975-989`) | absolute, 36×36, `translate(-50%,-50%)`, positioned by `toMapX/toMapY` |
| Selected | `.selected` (`:1022-1036, 1117-1134`) | expands to a bubble with a 36×36 circular thumb + short name |

Marker mode (`cluster` vs `named`/`venue`) is chosen by `getMapMarkerMode`
(`DiscoveryExperience.tsx:218-220`) at zoom 12.2.

## Examples

```jsx
{/* MapLibre markers are built as raw DOM strings (:828-851) */}
const element = document.createElement("button");
element.className = ["mapMarker", isSelected && !isClusterMode ? "selected" : "", isClusterMode ? "cluster" : "named"].filter(Boolean).join(" ");
element.setAttribute("aria-label", `Select ${venue.name}`);
element.innerHTML = isClusterMode
  ? `<strong>${getVenueClusterCount(venue, venues)}</strong>`
  : `<span class="markerBubble"><em>${getVenueMapIcon(venue)}</em><span>${venue.name}</span></span><span class="markerPin" aria-hidden="true"></span>`;

{/* Fallback marker is JSX, positioned by the linear projection (:989-1008) */}
<button
  className={selectedVenue?.id === venue.id ? "fallbackMarker selected" : "fallbackMarker"}
  type="button"
  style={{ left: `${toMapX(venue.longitude)}%`, top: `${toMapY(venue.latitude)}%` }}
  onClick={() => onSelectVenue(venue)}
  aria-label={`Select ${venue.name}`}
/>
```

## Rules

- Both systems keep the same class contract so styling is shared; changes must be mirrored in the DOM-string branch and the JSX branch.
- Position is the only inline style (`left`/`top` %, from `toMapX/toMapY`) — a computed value, allowed inline.
- The category glyph inside a named marker is an **emoji** from `getVenueMapIcon` (`☕`/`🍽`), not a `uiIcon`.
- Cluster count comes from `getVenueClusterCount`, which currently always returns 1 — clustering is a stub (see [R13](../05-inconsistencies.md#secondary-findings)).
- The user-location marker's blue (`#2f80ff`) is a hardcoded literal, not a token (R5).

## Do

- Give every marker button an `aria-label` naming the venue (both branches do).
- Keep the pin/bubble split so the pin anchors to the coordinate (`anchor:"bottom"`) while the bubble floats above.
- Use `.selected` to promote the active venue's marker.

## Don't

- Don't diverge the two systems' class names or the fallback will render differently from the live map (R8).
- Don't build markers with untrusted `innerHTML`; venue names are interpolated into a DOM string today — keep inputs controlled.
- Don't treat the cluster count as real until clustering is implemented (R13).

## Accessibility

- Each marker is a `<button>` with `aria-label="Select <name>"`.
- The pin and category emoji are decorative (`aria-hidden` on the pin; the emoji sits inside the labelled button).
- Selected markers add a thumbnail `img` with `alt=""` (decorative), since the button label already names the venue.

## Related components

- [Venue Map](venue-map.md) — hosts and manages markers.
- [Map Overlays](map-overlays.md) — the count legend/route summary that accompany markers.
- [Icon](icon.md) — note markers use emoji, not the `Icon` component.

## Files

- `apps/web/app/globals.css` — markers `:874-1042`; selected thumb/name `:1117-1134`.
- `apps/web/components/DiscoveryExperience.tsx` — build `:810-856`; fallback `:975-1014`; `getMapMarkerMode` `:218-220`; `getVenueClusterCount` `:1253-1257`; `getVenueMapIcon` `:1259-1265`; `toMapX/toMapY` `:1336-1342`.

## References

- [Patterns — Map interaction](../04-patterns.md#map-interaction-patterns)
- [Component catalog — Map markers](../03-components.md#map-markers-two-systems)
- [Inconsistency register — R5, R8, R13](../05-inconsistencies.md)

## Future improvements

- Implement real clustering so `getVenueClusterCount` aggregates (R13).
- Tokenize the user-location blue (R5).
- Replace `innerHTML` marker construction with DOM building to avoid interpolating names into markup.

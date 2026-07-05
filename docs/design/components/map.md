# Map

> **Status: As-built.** Tokens: [`../03-color-system.md`](../03-color-system.md).
> Behavior overlaps [`../05-grid-system.md`](../05-grid-system.md) (the map column) and
> [`../08-motion.md`](../08-motion.md) (camera).

## Purpose

The center map stage: a lazily-loaded MapLibre GL map with a styled canvas, a floating
topbar, markers (venue + user + cluster), map overlays (count legend, route summary), a
fallback layer before the map loads, and a "find near me" bottom action.

## Usage

Rendered by the `VenueMap` sub-component (`DiscoveryExperience.tsx:696-1030`).

- **Stage** `.mapStage` (`globals.css:735-740`): relative, `min-height:100vh`, `overflow:hidden`.
- **Topbar** `.mapTopbar` (`:742-759`): absolute, `pointer-events:none`; children re-enable `auto` so the map stays draggable. Hosts glass buttons + `.mapStats`.
- **Layers** (`:814-872`): `.mapFallbackGrid` (gradient stand-in) → `.mapCanvas` (MapLibre, with a saturation/contrast CSS filter) → `.mapShade` (vignette) → marker layer.
- **Markers** — base `.mapMarker/.fallbackMarker` (`:874-882`); `cluster`/`named` modes, `.markerBubble` + `.markerPin`, `userLocationMarker` (blue), and a JSX `.fallbackMarker` layer positioned by `toMapX/toMapY`. Two rendering systems share the class contract.
- **Overlays** — `.mapCountLegend` (`:1043-1073`) and `.routeSummary` (blue, `:1075-1115`).
- **Bottom action** `.mapBottomAction` (`:1136-1161`): centered primary pill.

## Examples

```jsx
{/* Lazy init (:730-753) */}
void import("maplibre-gl").then((maplibregl) => {
  const map = new maplibregl.Map({
    container: mapContainerRef.current,
    center: [106.683, 10.778], zoom: 12.8, attributionControl: false,
    style: buildMapStyle(initialThemeRef.current),   // CARTO dark-matter / positron
  });
  map.addControl(new maplibregl.NavigationControl({ visualizePitch: true }), "bottom-right");
});

{/* Canvas wrap + fallback + overlays (:970-1028) */}
<div className="mapCanvasWrap">
  <div className="mapFallbackGrid" aria-hidden="true" />
  <div ref={mapContainerRef} className="mapCanvas" />
  <div className="mapShade" aria-hidden="true" />
</div>
```

## Rules

- MapLibre is **lazy-imported**; its CSS is imported in `app/layout.tsx`.
- Marker mode is chosen by `getMapMarkerMode` at zoom 12.2 (`:218-220`); `getVenueClusterCount` is a stub that returns 1 (clustering not implemented).
- Theme changes call `setStyle` then re-apply `applyMapLayerPreferences` on `idle` (`:768-791`).
- Keep the topbar `pointer-events:none` with `auto` children or it blocks dragging.
- Camera durations are fixed (`flyTo` 500/700ms, `fitBounds` 700ms).
- Routing (`fetchRoute`) calls the public OSRM demo server directly — an external, un-configurable dependency, and a raw `fetch` outside `lib/api.ts`.
- Both marker systems must keep identical class names or the fallback diverges.
- Position (`left/top` %) is the only inline style on markers.

## Do / Don't

- **Do** keep the import lazy and clean up markers/route/map on unmount (`:755-765`).
- **Do** re-apply layer preferences after any `setStyle`.
- **Do** keep decorative layers `aria-hidden`.
- **Don't** eager-import MapLibre.
- **Don't** fetch map/route data via raw `fetch` in the component (route fetching already does — a known deviation).
- **Don't** depend on the OSRM demo server for anything production-critical.

## Accessibility

- The stage is `section[aria-label="Restaurant map"]`; markers are `<button>`s with `aria-label="Select <name>"`; decorative layers are `aria-hidden`.
- The map canvas isn't keyboard-navigable beyond MapLibre's controls; the left-panel rail is the non-map path to select venues.
- No `aria-live` on count/route changes; no `prefers-reduced-motion` for camera flights — see [`../09-accessibility.md`](../09-accessibility.md).

## Related components

- [Button](button.md) (topbar + bottom action) · [Card](card.md) (overlays are card-like) · [Icon](../07-icons.md)
- [Grid system](../05-grid-system.md) (map column) · [Motion](../08-motion.md) (camera).

## Files

- `apps/web/app/globals.css` — stage/topbar `:735-812`; canvas layers `:814-872`; markers `:874-1134`; overlays `:1043-1115`; bottom action `:1136-1161`.
- `apps/web/components/DiscoveryExperience.tsx` — `VenueMap` `:696-1030`; style/prefs `:136-216`; markers `:810-856, 975-1014`; `fetchRoute` `:1213-1251`; route layers `:909-968`.
- `apps/web/app/layout.tsx` — `import "maplibre-gl/dist/maplibre-gl.css"`.

## References

- [Grid system](../05-grid-system.md) · [Motion](../08-motion.md) · [Color](../03-color-system.md)

## Future improvements

- Move routing behind `lib/api.ts` with a configurable provider (not the OSRM demo).
- Implement real clustering; tokenize the user-location blue and the basemap palette.
- Add `aria-live` for count/route updates and `prefers-reduced-motion` for camera moves.

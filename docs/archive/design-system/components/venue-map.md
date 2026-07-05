# Venue Map

> As-built. Tokens live in [`../01-foundations.md`](../01-foundations.md); map
> interaction behavior is summarized in [`../04-patterns.md`](../04-patterns.md#map-interaction-patterns).

## Purpose

The center map stage: a lazily-loaded MapLibre GL map with a styled canvas, a floating
topbar, a fallback grid + markers before the map is ready, and a bottom "find near me"
action. Hosts markers and overlays.

## Usage

Rendered by the `VenueMap` sub-component (`DiscoveryExperience.tsx:696-1030`).

- `.mapStage` (`globals.css:735-740`): relative, `min-height:100vh`, `overflow:hidden`.
- `.mapTopbar` (`:742-759`): absolute top bar, `pointer-events:none`; children re-enable
  `pointer-events:auto` so the map stays draggable between controls. Holds glass
  buttons ([Button](button.md)) and `.mapStats` (`:808-812`).
- Stacked fill layers `.mapCanvasWrap`/`.mapCanvas`/`.mapFallbackGrid`/`.mapShade`/
  `.fallbackMarkerLayer` (`:814-872`): the MapLibre canvas gets a CSS
  saturation/contrast filter (`:837-843`); `.mapShade` adds a vignette; the fallback
  grid is a radial/linear gradient stand-in.
- `.mapBottomAction` (`:1136-1161`): centered floating primary pill ("Tìm quanh đây"),
  hover `scale(1.03)`.

## Examples

```jsx
{/* Lazy map init (:723-753) */}
void import("maplibre-gl").then((maplibregl) => {
  const map = new maplibregl.Map({
    container: mapContainerRef.current,
    center: [106.683, 10.778],   // HCMC
    zoom: 12.8,
    attributionControl: false,
    style: buildMapStyle(initialThemeRef.current),   // CARTO dark-matter / positron
  });
  map.addControl(new maplibregl.NavigationControl({ visualizePitch: true }), "bottom-right");
});

{/* Canvas wrap with fallback + overlays (:970-1028) */}
<div className="mapCanvasWrap">
  <div className="mapFallbackGrid" aria-hidden="true" />
  <div ref={mapContainerRef} className="mapCanvas" />
  <div className="mapShade" aria-hidden="true" />
  {/* fallback markers + count legend when !mapReady; route summary when activeRoute */}
</div>
```

## Rules

- MapLibre is **lazy-imported** (`import("maplibre-gl")`) so it stays off the initial bundle; the CSS is imported in `app/layout.tsx`.
- Theme changes call `map.setStyle` then re-apply `applyMapLayerPreferences` on `idle` (`:768-791`); marker mode re-syncs on `zoomend` (`:793-808`).
- Basemap tuning (`applyMapLayerPreferences`, `:145-216`) repaints roads/water/parks with a bespoke hardcoded palette per theme — outside the token system (R5).
- The topbar must keep `pointer-events:none` with `auto` children, or it will block map dragging.
- Camera moves are fixed durations: `flyTo` 500ms (select) / 700ms (user), `fitBounds` 700ms (route).
- Routing calls the public OSRM demo server directly (`fetchRoute`) — see [R14/R15](../05-inconsistencies.md#secondary-findings); an external, un-configurable dependency.
- Cleanup on unmount removes markers, user marker, route layers, and the map instance (`:755-765`).

## Do

- Keep the map import lazy and clean up in the effect's teardown.
- Re-apply layer preferences after any `setStyle` (theme switch).
- Keep decorative layers (`mapFallbackGrid`, `mapShade`) `aria-hidden`.

## Don't

- Don't eager-import MapLibre; it's large.
- Don't add map data fetches via raw `fetch` in the component — route fetching already violates the `lib/api.ts` rule (R15).
- Don't depend on the OSRM demo server for anything production-critical (R14).

## Accessibility

- The stage is a labelled region (`<section aria-label="Restaurant map">`).
- Decorative layers are `aria-hidden`; interactive markers/buttons carry their own labels ([Map Marker](map-marker.md), [Button](button.md)).
- The map canvas itself is not keyboard-navigable beyond MapLibre's built-in controls; venue selection is fully available via the left-panel rail as a non-map path.

## Related components

- [Map Marker](map-marker.md) — pins rendered onto the map.
- [Map Overlays](map-overlays.md) — count legend + route summary.
- [Button](button.md) — topbar glass buttons + bottom action.
- [App Shell](app-shell.md) — the middle column that contains the stage.

## Files

- `apps/web/app/globals.css` — stage/topbar `:735-812`; canvas layers `:814-872`; bottom action `:1136-1161`.
- `apps/web/components/DiscoveryExperience.tsx` — `VenueMap` `:696-1030`; map style/prefs `:136-216`; `fetchRoute` `:1213-1251`; route layers `:909-968`.
- `apps/web/app/layout.tsx` — `import "maplibre-gl/dist/maplibre-gl.css"`.

## References

- [Patterns — Map interaction](../04-patterns.md#map-interaction-patterns)
- [Layout — Map stage](../02-layout.md#map-stage--mapstage-735-740)
- [Inconsistency register — R5, R8, R13, R14, R15](../05-inconsistencies.md)

## Future improvements

- Move routing behind `lib/api.ts` and make the routing provider configurable rather than the public OSRM demo (R14/R15).
- Tokenize / centralize the basemap color palette (R5).

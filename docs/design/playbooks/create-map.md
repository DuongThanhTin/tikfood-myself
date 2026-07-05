# Playbook: Create / Extend the Map

> The map **exists** — see [`../components/map.md`](../components/map.md). This guide is for
> extending it (new marker type, overlay, or map-driven view). Don't add a second map stack.

## When to use

Adding markers, overlays, or map interactions to the existing MapLibre stage.

## Steps

1. **Work inside `VenueMap`** (`DiscoveryExperience.tsx:696-1030`); reuse its lazy-import,
   theme handling, and cleanup. Never eager-import MapLibre.
2. **New markers**: extend the shared marker classes (`.mapMarker`/`.markerBubble`/`.markerPin`)
   so the MapLibre DOM markers and the fallback layer stay visually identical. Position via
   `left/top` % (the only allowed inline style).
3. **New overlays**: model on `.mapCountLegend`/`.routeSummary` — glass-strong, blur18,
   absolute with a sensible `z-index` (markers 3–4, overlays 5–6, detail close 7).
4. **Theme**: if you touch basemap paint, do it in `applyMapLayerPreferences` and re-apply
   on `setStyle` `idle`. Prefer tokens; the current basemap palette is hardcoded (known gap).
5. **Data/routing**: fetch through `lib/api.ts` — do **not** repeat the raw-`fetch`-to-OSRM
   deviation; make external providers configurable and human-approved.
6. **a11y**: markers are `<button>`s with `aria-label`; decorative layers `aria-hidden`;
   consider `aria-live` for changing counts and `prefers-reduced-motion` for camera moves.

## Checklist

- [ ] Reused `VenueMap` + shared marker/overlay classes; both marker systems consistent.
- [ ] Lazy import + cleanup preserved.
- [ ] Overlay `z-index` correct; tokens used.
- [ ] No raw `fetch` for new data; external calls approved.
- [ ] Marker labels + hidden decoration; meets [DoD](../../verification/definition-of-done.md).

## Related

- [../components/map.md](../components/map.md) · [../08-motion.md](../08-motion.md) · [../05-grid-system.md](../05-grid-system.md)

## References

- `apps/web/components/DiscoveryExperience.tsx` (`:696-1030, 136-216`).

# Motion

> **Status: As-built.** Source: `globals.css` (CSS transitions) + `DiscoveryExperience.tsx`
> (map camera).

## Principles in code

Motion is subtle and functional: hover feedback, panel collapse, image zoom, and map
camera moves. **There are no `@keyframes`** — no looping or entrance animations.

## Transition durations & easing

Easing is always `ease`. Durations present:

| Duration | Where |
|----------|-------|
| `0.18s` (×5) | marker bubble / fallback marker (`:919, 984`) |
| `0.2s` (×8) | chips, tabs, buttons, thumbnail color, bottom action |
| `0.25s` (×1) | hero-search glow opacity (`:340`) |
| `0.3s` (×2) | left-panel width collapse, sidebar toggle transform (`:135, 161`) |
| `0.45s` (×2) | image zoom on hover (`.venueThumb img`, `.videoCard img`, `:621, 1436`) |

## Hover transforms

- `translateY(-1px)` — chips, marker bubble.
- `scale(1.03–1.08)` — bottom action, thumbnails, fallback marker.

## Map camera motion (JS)

- `flyTo` **500ms** on venue select (`:867`); **700ms** on user location (`:901`).
- `fitBounds` **700ms** on route with asymmetric padding (`:956-964`).

## Rules

- Use `ease`; pick an existing duration (`0.2s` for control feedback, `0.3s` for layout, `0.45s` for image zoom).
- Keep motion functional (feedback/navigation), not decorative loops.
- Map camera durations are fixed constants in the map effects.

## Do / Don't

- **Do** reuse the established durations for consistency.
- **Do** animate `transform`/`opacity` (cheap) as the existing rules do.
- **Don't** add looping/entrance keyframes; none exist and the aesthetic is calm.
- **Don't** animate layout-affecting properties beyond the intentional panel-width transition.

## Accessibility

- There is **no `prefers-reduced-motion` handling** today. The zoom/scale hovers and
  700ms camera flights would ideally be reduced for users who request it — a known gap.

## Related

- [Grid system](05-grid-system.md) (panel collapse) · [components/map.md](components/map.md) (camera)

## Files

- `apps/web/app/globals.css` — transitions throughout; no `@keyframes`.
- `apps/web/components/DiscoveryExperience.tsx` — `flyTo`/`fitBounds` `:863-968`.

## References

- `apps/web/CLAUDE.md`.

## Future improvements

- Add a `@media (prefers-reduced-motion: reduce)` block to soften hovers and camera moves.
- Tokenize durations (`--motion-fast/base/slow`) rather than repeating literals.

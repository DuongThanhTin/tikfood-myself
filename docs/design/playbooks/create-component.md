# Playbook: Create a Component

> Add a reusable UI component to `apps/web`. Components today live as local
> sub-components inside `DiscoveryExperience.tsx` with styles in `globals.css`.

## When to use

A visual unit is reused (or clearly will be) and doesn't already exist. **Reuse beats
create** — check first.

## Prerequisites

- Read [`../00-overview.md`](../00-overview.md) and foundations `03–09`.
- Search `globals.css` and `components/` for an existing class/component you can extend.

## Steps

1. **Confirm it's new.** Grep classes and component names; if a variant of an existing
   piece fits ([button](../components/button.md), [card](../components/card.md),
   [input](../components/input.md)), extend that instead.
2. **Name it descriptively** (`VenueRailCard`, not `Card2`) — `apps/web/CLAUDE.md`.
3. **Style with tokens.** Use color/spacing/type/radius from foundations; reuse the base
   class groups (e.g. the button base, the glass surface). No hardcoded hex/px where a
   token exists.
4. **Structure semantically.** Real `<button>`/`<label>`/`<table>`; inline `style` only
   for computed values (e.g. positions, a progress width).
5. **Wire data via props**; if it needs API data, the data comes from a parent that used
   `lib/api.ts`.
6. **Accessibility**: `aria-label` on icon-only controls, `aria-hidden` on decoration,
   `alt` on images, native controls, and a `:focus-visible` treatment (the system-wide gap
   — add one, per [`../09-accessibility.md`](../09-accessibility.md)).
7. **Document it.** Add a doc under [`../components/`](../components) using the standard
   template (Purpose · Usage · Examples · Rules · Do · Don't · Accessibility · Related ·
   Files · References · Future Improvements).

## Checklist

- [ ] Verified nothing existing covers it.
- [ ] Descriptive name; tokens not literals; reused base classes.
- [ ] Semantic markup; inline style only for computed values.
- [ ] a11y: labels, decoration hidden, focus visible.
- [ ] Component doc added; cross-linked.
- [ ] Meets [`definition-of-done`](../../verification/definition-of-done.md).

## Related

- [update-component.md](update-component.md) · [review-ui.md](review-ui.md) · [../03-color-system.md](../03-color-system.md) · [../06-typography.md](../06-typography.md)

## References

- `apps/web/CLAUDE.md` · `apps/web/app/globals.css` · `apps/web/components/DiscoveryExperience.tsx`.

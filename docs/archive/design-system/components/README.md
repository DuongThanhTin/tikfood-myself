# Component Documentation

Per-component reference for the `apps/web` design system, one document each, following a
fixed template: **Purpose · Usage · Examples · Rules · Do · Don't · Accessibility ·
Related Components · Files · References · Future Improvements**.

These are the *detailed* docs. For the at-a-glance catalog see
[`../03-components.md`](../03-components.md); for tokens see
[`../01-foundations.md`](../01-foundations.md). All examples are real code with source
line references; nothing is invented, and inconsistencies are flagged to the register
[`../05-inconsistencies.md`](../05-inconsistencies.md) rather than fixed.

## Controls

- [Button](button.md) — primary / secondary / glass / icon-glass / glass-action
- [Chip & Tag](chip-tag.md) — quick-prompt chips and toggle filter tags
- [Form Field & Select](form-field-select.md) — labelled native dropdowns
- [Hero Search](hero-search.md) — the primary search input with glow
- [Toggle](toggle.md) — the open-now checkbox
- [Icon](icon.md) — the unicode-glyph icon primitive
- [Navigation Tabs](navigation-tabs.md) — side tabs and collapsed mini tabs

## Content & cards

- [Badge](badge.md) — status/category pills
- [Venue Rail Card](venue-rail-card.md) — the selectable venue row
- [Trend Score Card](trend-score-card.md) — trend percentage + gradient track
- [Video Card](video-card.md) — 9:16 social-video thumbnail

## Map

- [Map Marker](map-marker.md) — venue/user pins (two rendering systems)
- [Map Overlays](map-overlays.md) — count legend + route summary
- [Venue Map](venue-map.md) — the MapLibre map stage

## Structure

- [Venue Detail Panel](venue-detail-panel.md) — the right-hand detail composite
- [App Shell](app-shell.md) — the three-column grid container + theme

## Where components live

Every component is a local sub-component inside
`apps/web/components/DiscoveryExperience.tsx` (`VenueRailCard`, `VenueMap`,
`VenueDetail`, `VideoCard`, `Icon`); their styles are in `apps/web/app/globals.css`.
`components/VenueList.tsx` is **orphaned** (see [R1](../05-inconsistencies.md#r1--venuelisttsx-is-orphaned-its-classes-dont-exist)) and is not documented as a live component.

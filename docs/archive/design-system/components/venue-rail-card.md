# Venue Rail Card

> As-built. Tokens live in [`../01-foundations.md`](../01-foundations.md).

## Purpose

The selectable venue row in the left-panel recommendation rail ("Gợi ý cho bạn"). A
horizontal card: square thumbnail + summary (name, cuisine · district, trend badge,
star rating). Selecting it drives the map and opens the detail panel.

## Usage

Rendered by the `VenueRailCard` sub-component (`DiscoveryExperience.tsx:662-694`) inside
`.venueList` (`globals.css:586-589`, `display:grid; gap:18px`), under a `.sectionHeader`
(`:562-584`).

- `.venueRailCard` (`:591-602`): 2-col grid `96px minmax(0,1fr)`, `gap:16px`,
  `border:0`, `border-radius:14px`, transparent bg, left-aligned. It is a `<button>`.
- `.venueThumb` (`:609-626`): 96×96, radius 12, `--surface-high`; `img` covers and
  zooms `scale(1.08)` on card hover.
- `.venueSummary` (`:628-648`): `strong` name (15px/800, ellipsis, → `--primary` on
  hover/selected), `small` "cuisine · district".
- `.venueMetaRow` (`:650-677`): trend badge `em` + star rating span (`--secondary`).
- Responsive: thumbnail shrinks to 82px under 640px (`:1572-1579`).

## Examples

```jsx
{/* Rail (:547-565) */}
<section className="venueRail" aria-label="Recommended venues">
  <div className="sectionHeader">
    <h2>Gợi ý cho bạn</h2>
    <button type="button" onClick={() => void runSearch({ sort: "trending", limit: 20 })}>Tất cả</button>
  </div>
  {venues.length === 0 ? <p className="emptyText">Không có địa điểm phù hợp bộ lọc.</p> : null}
  <div className="venueList">
    {venues.map((venue) => (
      <VenueRailCard key={venue.id} venue={venue} selected={selectedVenue?.id === venue.id} onSelect={() => selectVenue(venue)} />
    ))}
  </div>
</section>

{/* Card internals (:672-693) */}
<button className={selected ? "venueRailCard selected" : "venueRailCard"} type="button" onClick={onSelect}>
  <span className="venueThumb"><img alt={`${venue.name} preview`} src={media.image} /></span>
  <span className="venueSummary">
    <strong>{venue.name}</strong>
    <small>{media.cuisine} · {venue.district}</small>
    <span className="venueMetaRow">
      <em>{media.badge}</em>
      <span><Icon name="star" />{media.rating}</span>
    </span>
  </span>
</button>
```

## Rules

- The whole card is one `<button>`; the thumbnail and summary are non-interactive `<span>`s inside it.
- Selection is controlled by the parent via the `selected` prop → `.selected` class; hover and selected share the same title-highlight treatment (`:604-607`).
- Name truncates with ellipsis; keep it on one line.
- `media.image` / `media.cuisine` / `media.badge` / `media.rating` come from `getVenueMedia` (mocked — see [R7](../05-inconsistencies.md#r7--presentation-data-is-mocked-not-api-driven)); only `venue.name` / `venue.district` are real API data.

## Do

- Keep the card a single button for one clear tap target.
- Use `venue.name`/`venue.district` from the API; treat thumbnail/rating/badge as placeholder media until wired.
- Let `.venueList` own inter-card spacing.

## Don't

- Don't nest additional buttons/links inside the card (it's already a button).
- Don't present mocked rating/badge as real venue data (R7).
- Don't hardcode the thumb size in two places; the responsive rule already handles 640px.

## Accessibility

- The card is a `<button>` with an accessible name from its text content (venue name, cuisine, district, rating).
- The star `<Icon>` is decorative (`aria-hidden`); the numeric rating text carries the meaning.
- The thumbnail `img` has descriptive `alt` (`"<name> preview"`).
- Selected state is visual only (title color) — no `aria-pressed`/`aria-current`.

## Related components

- [Badge](badge.md) — the `venueMetaRow em` trend badge.
- [Icon](icon.md) — the star glyph.
- [Venue Detail Panel](venue-detail-panel.md) — opened when a card is selected.
- [Venue Map](venue-map.md) — flies to the selected venue.

## Files

- `apps/web/app/globals.css` — rail/header `:558-589`; card `:591-677`; responsive `:1572-1579`.
- `apps/web/components/DiscoveryExperience.tsx` — rail `:547-565`; `VenueRailCard` `:662-694`; `getVenueMedia` `:1209-1211`.

## References

- [Component catalog — Cards](../03-components.md#cards)
- [Patterns — UI states (empty)](../04-patterns.md#ui-states)
- [Inconsistency register — R7](../05-inconsistencies.md)

## Future improvements

- Replace mocked `mediaByVenue` imagery/rating/badge with API fields so each card shows its own real media (R7).
- Add `aria-current`/`aria-pressed` to reflect the selected card.

# Venue Detail Panel

> As-built. Tokens live in [`../01-foundations.md`](../01-foundations.md). This is a
> composite; sub-parts with their own docs are linked, not re-described.

## Purpose

The right-hand panel showing a selected venue in full: hero image with badges, title +
rating + quick actions, trend score, AI/about copy, categories, social videos, dishes,
opening hours, and a menu CTA. Mounts only when a venue is selected.

## Usage

Rendered by the `VenueDetail` sub-component (`DiscoveryExperience.tsx:1032-1178`) inside
`.rightPanel` (see [App Shell](app-shell.md)); the panel also hosts the `.closeButton`.

- `.closeButton` (`globals.css:1163-1177`): absolute, `z-index:7`, top-right, 42×42
  circle, glass, `blur(18px)`, white glyph.
- `.detailHero` (`:1183-1201`): 320px (→ 240px under 640px), `object-fit:cover`, dark
  gradient `::after`; hosts `.detailBadges` → see [Badge](badge.md).
- `.detailBody` (`:1224-1228`): grid, `gap:26px`, `padding:28px 28px 32px`.
- `.detailTitleRow` (`:1230-1249`): grid `minmax(0,1fr) auto` (→ 1 col under 640px);
  `h2` (Plus Jakarta 30px/1.14), `.detailMeta` (address · district).
- `.compactDetailActions` (`:1251-1274`): 3-col grid of `glassAction` buttons
  (route/share/save), min-height 36, 11px/**850** (weight not loaded — R4).
- `.ratingBlock` (`:1276-1295`): right-aligned; `--secondary` number (Plus Jakarta
  22px) + star + review count.
- `.detailStatus` (`:1356-1361`): loading line, muted, negative `margin:-8px 0` (R17).
- `.aboutText` (`:1351-1354`): the about copy.
- `.category, .dish` (`:1363-1377`): pill chips — category = primary on rgba-primary;
  dish = on-surface on glass.
- `.dishDetailItem, .openingRow` (`:1386-1415`): row items (name/day + price/hours).
- `.detailActions` / `.detailMenuAction` (`:1463-1471`): the full-width menu CTA row.

## Examples

```jsx
{/* Panel mount + close (:632-657) */}
{selectedVenue ? (
  <aside className="rightPanel" aria-label="Venue detail">
    <button className="closeButton" type="button" aria-label="Close venue detail" onClick={() => setSelectedVenue(null)}>
      <Icon name="close" />
    </button>
    <VenueDetail venue={selectedVenue} /* …handlers/state… */ />
  </aside>
) : null}

{/* Title row + quick actions (:1068-1094) */}
<div className="detailTitleRow">
  <div>
    <h2>{venue.name}</h2>
    <p className="detailMeta">{venue.address}, {venue.district}</p>
    <div className="compactDetailActions">
      <button className="glassAction" type="button" onClick={isActiveRoute ? onClearRoute : onRequestRoute} disabled={isRouting}>
        <Icon name="route" />{isRouting ? "Đang tải" : isActiveRoute ? "Xóa route" : "Chỉ đường"}
      </button>
      <button className="glassAction" type="button"><Icon name="share" />Chia sẻ</button>
      <button className="glassAction" type="button"><Icon name="bookmark" />Lưu</button>
    </div>
  </div>
  <div className="ratingBlock"><span><Icon name="star" />{media.rating}</span><small>{media.reviews}</small></div>
</div>
```

## Rules

- The panel mounts conditionally (`selectedVenue` truthy); it isn't hidden with CSS.
- Detail data loads lazily: on select, `fetchVenueDetail` runs only if `dishes`/`opening_hours` are absent (`:255-285`), showing `.detailStatus` meanwhile.
- Dishes render as detailed rows when `venue.dishes` exists, else fall back to plain `.dish` chips from `trending_dishes` (`:1138-1156`).
- Opening hours render only when present; days use `formatDay` abbreviations (`:1332-1334`).
- Real API fields: `name`, `address`, `district`, `about`, `ai_summary`, `categories`, `dishes`, `opening_hours`. Mocked: `rating`, `reviews`, badges, hero/thumbnail images (`media.*` — R7).
- "Chia sẻ", "Lưu", and the menu CTA have no handlers today (R10).

## Do

- Keep the close button as the single dismissal affordance (`aria-label` present).
- Use the lazy-detail pattern; show `.detailStatus` during load.
- Prefer detailed dish rows when dish data exists; fall back to chips otherwise.
- Use real API fields for text; treat `media.*` as placeholder.

## Don't

- Don't render share/save/menu as enabled no-ops without wiring (R10).
- Don't rely on `font-weight:850` for the compact actions (falls back — R4).
- Don't present mocked rating/reviews/images as real (R7).

## Accessibility

- The panel is `<aside aria-label="Venue detail">`; the close button has `aria-label="Close venue detail"`.
- Section headings (`<h3>`) structure the body for screen-reader navigation.
- The rating star and action icons are decorative; adjacent text/labels carry meaning.
- Gaps: share/save/menu buttons are enabled but do nothing (no `disabled`, R10); the loading line isn't an `aria-live` region.

## Related components

- [Badge](badge.md) · [Trend Score Card](trend-score-card.md) · [Video Card](video-card.md) — embedded sections.
- [Button](button.md) — `glassAction`, close, and the `primaryButton large` CTA.
- [Icon](icon.md) — route/share/bookmark/star/close glyphs.
- [App Shell](app-shell.md) — the `rightPanel` host and `detailOpen` grid state.
- [Venue Map](venue-map.md) — routing interaction mirrored here.

## Files

- `apps/web/app/globals.css` — close `:1163-1177`; hero/body/title/actions/rating `:1179-1295`; status/about `:1351-1361`; category/dish/rows `:1363-1415`; menu CTA `:1463-1471`.
- `apps/web/components/DiscoveryExperience.tsx` — mount `:632-657`; `VenueDetail` `:1032-1178`; lazy load `:255-285`; `formatDay` `:1332-1334`.
- `apps/web/lib/api.ts` — `Venue`/`VenueDish`/`OpeningHour` `:1-61`.

## References

- [Patterns — Data-flow, UI states](../04-patterns.md#data-flow-pattern)
- [Component catalog — Detail-panel building blocks](../03-components.md#detail-panel-building-blocks)
- [Inconsistency register — R4, R7, R10, R17](../05-inconsistencies.md)

## Future improvements

- Wire up share/save/menu or disable them (R10).
- Source rating/reviews/images from the API (R7).
- Add an `aria-live` region for the lazy-load status.

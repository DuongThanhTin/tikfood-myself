# Card

> **Status: As-built.** Tokens: [`../03-color-system.md`](../03-color-system.md),
> [`../04-spacing-system.md`](../04-spacing-system.md).

## Purpose

Container surfaces that group a unit of content. Three real card variants exist:

- **Venue rail card** — the selectable venue row in the left rail.
- **Trend score card** — trend percentage + gradient track + AI summary.
- **Video card** — 9:16 social-video thumbnail.

(The glass overlays — count legend, route summary — are card-like but map-scoped; see
[`map.md`](map.md).)

## Usage

### Venue rail card — `.venueRailCard` (`globals.css:591-677`)
2-col grid (`96px minmax(0,1fr)`, `gap:16px`), borderless `<button>`, radius 14. Hover/
selected turns the title `--primary` and zooms the thumb `scale(1.08)`. Parts:
`.venueThumb` (96×96, radius 12), `.venueSummary` (name `strong` ellipsis, `small`
"cuisine · district"), `.venueMetaRow` (badge `em` + star rating). Thumb → 82px ≤640px.

### Trend score card — `.trendScoreCard` (`:1304-1349`)
Glass base (shared with `.glassAction`, `:1297-1302`), radius 18, `blur(18px)`. Header
row (label + `--primary` percentage) over `.trendTrack` (8px pill) whose fill is a
`--secondary → --trend-fire` gradient at an inline width. Body `p` = AI summary (italic).

### Video card — `.videoCard` (`:1424-1461`)
`aspect-ratio:9/16`, radius 14; image zoom on hover; dark gradient `::after`; caption
`span` (bottom-left, white) = play glyph + platform + views. Laid out in `.videoGrid`
(2-col).

## Examples

```jsx
{/* Rail card (:672-693) */}
<button className={selected ? "venueRailCard selected" : "venueRailCard"} type="button" onClick={onSelect}>
  <span className="venueThumb"><img alt={`${venue.name} preview`} src={media.image} /></span>
  <span className="venueSummary">
    <strong>{venue.name}</strong>
    <small>{media.cuisine} · {venue.district}</small>
    <span className="venueMetaRow"><em>{media.badge}</em><span><Icon name="star" />{media.rating}</span></span>
  </span>
</button>

{/* Trend score card (:1104-1113) — only the fill width is inline */}
<section className="trendScoreCard">
  <div><span>Trend Score</span><strong>{trendPercent}%</strong></div>
  <div className="trendTrack"><span style={{ width: `${trendPercent}%` }} /></div>
  <p>{venue.ai_summary}</p>
</section>
```

## Rules

- The rail card is a single `<button>`; inner elements are non-interactive `<span>`s.
- Trend fill width is the only inline style (a computed value, allowed per `apps/web/CLAUDE.md`); clamp the score 1–99 (`:1054`).
- Video cards use `aspect-ratio`, not fixed heights; show at most 2 (`social_videos.slice(0,2)`).
- **Real vs mocked data:** `venue.name/district`, `trend_score`, `ai_summary`, `social_videos[].thumbnail_url` are real; `media.rating/reviews/badge/cuisine/image` come from the hardcoded `mediaByVenue` map (mocked).

## Do / Don't

- **Do** keep the rail card one tap target; use API fields for text.
- **Do** clamp trend score; enforce the 9:16 video ratio.
- **Don't** present mocked rating/badge/images as real venue data.
- **Don't** render a "playable" video card that has no action (current gap).

## Accessibility

- Rail card is a `<button>` named by its text; star/play icons are `aria-hidden`.
- Thumbnails have descriptive `alt`; decorative fallback images use `alt=""`.
- Selected state is color-only (no `aria-pressed`/`aria-current`); trend track has no `role="progressbar"` — see [`../09-accessibility.md`](../09-accessibility.md).

## Related components

- [Badge](../pages/restaurant.md) (trend badge — Proposed page hosts detail badges) · [Icon](../07-icons.md) · [Map](map.md) · [pages/home.md](../pages/home.md)
- [Input](input.md) (filters that populate the rail).

## Files

- `apps/web/app/globals.css` — rail `:558-677`; trend card `:1297-1349`; video card `:1417-1461`.
- `apps/web/components/DiscoveryExperience.tsx` — `VenueRailCard` `:662-694`; trend `:1054, 1104-1113`; `VideoCard` `:1180-1207`; `getVenueMedia` `:1209-1211`.
- `apps/web/lib/api.ts` — `Venue`/`SocialVideo` `:1-37`.

## References

- [Color](../03-color-system.md) · [Spacing](../04-spacing-system.md) · [Accessibility](../09-accessibility.md)

## Future improvements

- Source rating/reviews/images/badge from the API (drop the mocked `mediaByVenue`).
- Add `role="progressbar"` to the trend track; make video cards actually open the source.

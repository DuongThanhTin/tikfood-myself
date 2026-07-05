# Video Card

> As-built. Tokens live in [`../01-foundations.md`](../01-foundations.md).

## Purpose

A vertical (9:16) social-video thumbnail shown in the detail panel's "Video social nổi
bật" grid. Displays a poster image with a gradient scrim and a caption line (platform ·
view count) over a play glyph.

## Usage

- `.videoGrid` (`globals.css:1417-1422`): 2-column grid, `gap:12px`, `margin-top:12px`.
- `.videoCard` (`:1424-1430`): relative, `overflow:hidden`, radius 14,
  `aspect-ratio:9/16`, `--surface-high` bg.
- `img` (`:1432-1441`): covers; zooms `scale(1.08)` on hover.
- `::after` (`:1443-1448`): bottom-weighted dark gradient scrim.
- `> span` (`:1450-1461`): absolute bottom-left, white 11px/900, `gap:4px` (play icon +
  text).

Rendered by the `VideoCard` sub-component (`DiscoveryExperience.tsx:1180-1190`); data
comes from `getDisplayVideos` (`:1192-1207`), which uses up to 2 real `social_videos`
or falls back to mocked `media.videoA/B`.

## Examples

```jsx
{/* Grid in the detail body (:1127-1134) */}
<section>
  <h3>Video social nổi bật</h3>
  <div className="videoGrid">
    {getDisplayVideos(venue, media).map((video) => (
      <VideoCard key={video.id} image={video.image} views={video.views} label={video.label} platform={video.platform} />
    ))}
  </div>
</section>

{/* The card (:1180-1190) */}
<div className="videoCard">
  <img alt={label} src={image} />
  <span>
    <Icon name="play" />
    {platform ? `${platform} · ` : ""}{views}
  </span>
</div>
```

## Rules

- Enforce the 9:16 ratio via `aspect-ratio` (not fixed pixel heights) so cards stay proportional in the 2-col grid.
- The caption is a single overlaid line; keep it to platform + formatted view count.
- View counts are pre-formatted by `formatCompactNumber` (`:1322-1330`) into `K`/`M`.
- Show at most 2 videos (`social_videos.slice(0, 2)`); beyond that is out of the current pattern.
- Image sources may be real (`thumbnail_url`) or mocked fallbacks — see [R7](../05-inconsistencies.md#r7--presentation-data-is-mocked-not-api-driven).

## Do

- Pass a descriptive `label` (used as `alt`) — video caption or "<name> video".
- Let the gradient scrim guarantee caption legibility over any image.
- Keep the grid at 2 columns; it's tuned for portrait thumbnails.

## Don't

- Don't set a fixed height; rely on `aspect-ratio`.
- Don't render more than 2 cards per venue (breaks the layout intent).
- Don't treat the fallback images as real when `social_videos` is empty (R7).

## Accessibility

- The poster `img` has `alt` from `label` (caption or generated name).
- The play `<Icon>` is decorative (`aria-hidden`); the caption text carries platform/views.
- The card is **not interactive** today — it's a static thumbnail (no link/handler), so there's no focus/label requirement, but also no way to actually play the video.

## Related components

- [Icon](icon.md) — the play glyph.
- [Venue Detail Panel](venue-detail-panel.md) — the host section.

## Files

- `apps/web/app/globals.css` — `.videoGrid` `:1417-1422`; `.videoCard` `:1424-1461`.
- `apps/web/components/DiscoveryExperience.tsx` — `VideoCard` `:1180-1190`; `getDisplayVideos` `:1192-1207`; `formatCompactNumber` `:1322-1330`.
- `apps/web/lib/api.ts` — `SocialVideo` `:27-37`.

## References

- [Component catalog — Cards](../03-components.md#cards)
- [Inconsistency register — R7](../05-inconsistencies.md)

## Future improvements

- The card looks playable (play glyph) but has no action — either link to the source video or make the affordance honest.

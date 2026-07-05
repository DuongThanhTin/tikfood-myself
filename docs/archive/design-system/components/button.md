# Button

> As-built documentation of the general-purpose button family in `apps/web`.
> Tokens (colors, radii, shadows, motion) are defined in
> [`../01-foundations.md`](../01-foundations.md) and referenced here, not repeated.

## Purpose

Trigger an action or submit intent. The button family covers the solid brand action
(`primaryButton`), the low-emphasis companion (`secondaryButton`), and the translucent
"glass" set used over the map and in the detail panel (`glassButton`,
`iconGlassButton`, `glassAction`). Purpose-specific buttons that live inside another
component (nav tabs, close, bottom action, section link) are documented with their
host — see [Related components](#related-components).

## Usage

All variants share one base rule group: `1px solid var(--outline)`,
`border-radius: 12px`, `font-weight: 800` (`globals.css:505-513`). Pick a variant by
emphasis and surface:

| Variant | Class | When to use |
|---------|-------|-------------|
| Primary | `primaryButton` | The single most important action in a view (search, sign-in) |
| Primary large | `primaryButton large` | Full-width primary that spans a grid row (detail CTA) |
| Secondary | `secondaryButton` | Companion to a primary (reset/cancel) |
| Glass | `glassButton` | Actions floating over the map |
| Glass solid | `glassButton solid` | A glass button needing more contrast (e.g. "Gần tôi") |
| Icon glass | `iconGlassButton` | Icon-only circular action over the map |
| Glass action | `glassAction` | Detail-panel action rows (route/share/save) |

`type="button"` is always set explicitly. Loading is expressed by setting `disabled`
and swapping the label text (there is no spinner).

## Examples

Real usages from `components/DiscoveryExperience.tsx`:

```jsx
{/* Primary — loading via disabled + label swap (:535) */}
<button className="primaryButton" type="button" onClick={() => void runSearch()} disabled={isLoading}>
  {isLoading ? "Đang tìm" : "Tìm kiếm"}
</button>

{/* Secondary (:538) */}
<button className="secondaryButton" type="button" onClick={clearFilters} disabled={isLoading}>
  Reset
</button>

{/* Primary large — spans the grid row (:1173) */}
<button className="primaryButton large" type="button">Xem chi tiết & Menu</button>

{/* Glass + solid over the map (:602) */}
<button className="glassButton solid" type="button" onClick={useCurrentLocation}>Gần tôi</button>

{/* Icon-only glass — must carry aria-label (:605) */}
<button className="iconGlassButton" type="button" aria-label="Help">
  <Icon name="help" />
</button>

{/* Glass action with leading icon + loading label (:1073) */}
<button className="glassAction" type="button" onClick={isActiveRoute ? onClearRoute : onRequestRoute} disabled={isRouting}>
  <Icon name="route" />
  {isRouting ? "Đang tải" : isActiveRoute ? "Xóa route" : "Chỉ đường"}
</button>
```

## Rules

- Every button sets `type="button"` (only the absent, future form-submit case would use `submit`).
- Emphasis is exclusive: one primary per view; secondary is always paired with a primary.
- Primary text color is the fixed `#4e2600` on `var(--primary)` — a hardcoded value, not a token (see [R5](../05-inconsistencies.md#r5--blue-accent-has-no-token-warm-colors-re-expressed-as-literals)).
- Glass variants require a blurred/translucent backdrop to read correctly; don't place them on opaque panels where `secondaryButton` fits better.
- Loading/disabled uses the native `disabled` attribute (primary/secondary → `opacity:.7; cursor:wait`; `glassAction` → `opacity:.58`).
- Icon + text order: the `<Icon>` precedes the label; the base rule provides `gap: 8px` on the inline-flex glass group.

## Do

- Use `primaryButton large` when the button must fill a `grid-column: 1 / -1` row (detail CTA).
- Give every icon-only button an `aria-label` describing the action.
- Keep label text short and in Vietnamese for UI copy (project convention).
- Reuse an existing variant; the base group already standardizes border, radius, weight.

## Don't

- Don't render a button with no handler and no `disabled` — it looks actionable but isn't (this exists today: "Bộ lọc", "Chia sẻ", "Lưu", help, "Đăng nhập ngay" — see [R10](../05-inconsistencies.md#secondary-findings)).
- Don't introduce a new one-off button class when a variant + modifier covers it.
- Don't hardcode the amber/orange again as `rgba(...)` — prefer the token where the value equals a token.
- Don't rely on `font-weight` above 800; 850/900 aren't loaded and fall back (see [R4](../05-inconsistencies.md#r4--font-weights-used-that-arent-loaded)).

## Accessibility

- **Icon-only buttons** (`iconGlassButton`, and the close/mini-tab buttons elsewhere) always pair with `aria-label`; the inner glyph is `aria-hidden` via `.uiIcon`.
- **Text buttons** expose their accessible name through visible text — no extra label needed.
- **State**: loading is communicated to assistive tech only via the native `disabled` attribute; there is no `aria-busy`. Buttons that do nothing yet are neither disabled nor labelled as unavailable (accessibility gap, R10).
- **Focus**: relies on the browser default focus ring (no custom `:focus-visible` style is defined for buttons).

## Related components

- [Icon](icon.md) — leading/trailing glyphs inside buttons.
- [Navigation Tabs](navigation-tabs.md) — `sideTab` / `miniTab` button-like controls.
- [Chip & Tag](chip-tag.md) — pill-shaped selectable buttons (distinct family).
- [Venue Detail Panel](venue-detail-panel.md) — hosts `glassAction` rows and the close button.

## Files

- `apps/web/app/globals.css` — base group `:505-513`; `.primaryButton` `:515-530`; `.secondaryButton` `:532-543`; glass group `:761-806`; `.glassAction` `:1298-1302, 1473-1482`.
- `apps/web/components/DiscoveryExperience.tsx` — usages at `:535, 538, 550, 569, 589, 594, 602, 605, 625, 698, 1073, 1077, 1081, 1173`.

## References

- [Foundations — Color / Typography / Motion](../01-foundations.md)
- [Component catalog — Buttons](../03-components.md#buttons)
- [Inconsistency register — R4, R5, R10](../05-inconsistencies.md)
- Frontend rules: [`apps/web/CLAUDE.md`](../../../apps/web/CLAUDE.md)

## Future improvements

Grounded in the register; recorded as opportunities only (no code change proposed by this audit):

- Tokenize the primary button text color and per-variant shadows (currently literals) — R5.
- Wire up or explicitly `disabled` the no-op buttons — R10.
- Replace the label-swap loading pattern with a spinner/`aria-busy` for clearer state.
- Add a shared `:focus-visible` treatment for keyboard users.

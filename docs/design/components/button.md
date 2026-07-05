# Button

> **Status: As-built.** Tokens: [`../03-color-system.md`](../03-color-system.md),
> [`../06-typography.md`](../06-typography.md). Radii/spacing:
> [`../04-spacing-system.md`](../04-spacing-system.md).

## Purpose

Trigger an action. Covers the solid brand action (`primaryButton`), its low-emphasis
companion (`secondaryButton`), and the translucent glass set used over the map and in
the detail panel (`glassButton`, `iconGlassButton`, `glassAction`). Purpose-specific
buttons embedded in other components (nav tabs, close, bottom action) are documented
with their host.

## Usage

Base group `.primaryButton, .secondaryButton, .glassButton, .iconGlassButton,
.glassAction` (`globals.css:505-513`): `1px solid var(--outline)`,
`border-radius:12px`, `font-weight:800`.

| Variant | Class | When |
|---------|-------|------|
| Primary | `primaryButton` | The single most important action (search, sign-in) |
| Primary large | `primaryButton large` | Full-width primary spanning a grid row |
| Secondary | `secondaryButton` | Companion to a primary (reset) |
| Glass | `glassButton` / `glassButton solid` | Actions floating over the map |
| Icon glass | `iconGlassButton` | Icon-only circular map action |
| Glass action | `glassAction` | Detail-panel action rows (route/share/save) |

`type="button"` is always set. Loading = `disabled` + label swap (no spinner).

## Examples

```jsx
<button className="primaryButton" type="button" onClick={() => void runSearch()} disabled={isLoading}>
  {isLoading ? "Đang tìm" : "Tìm kiếm"}
</button>
<button className="secondaryButton" type="button" onClick={clearFilters} disabled={isLoading}>Reset</button>
<button className="glassButton solid" type="button" onClick={useCurrentLocation}>Gần tôi</button>
<button className="iconGlassButton" type="button" aria-label="Help"><Icon name="help" /></button>
<button className="glassAction" type="button" onClick={isActiveRoute ? onClearRoute : onRequestRoute} disabled={isRouting}>
  <Icon name="route" />{isRouting ? "Đang tải" : isActiveRoute ? "Xóa route" : "Chỉ đường"}
</button>
```
(`DiscoveryExperience.tsx:535, 538, 602, 605, 1073`)

## Rules

- Always set `type="button"`.
- One primary per view; secondary is always paired with a primary.
- Primary text is `#4e2600` on `--primary` — currently a literal (no `--on-primary` token; see [color §non-tokenized](../03-color-system.md#non-tokenized-colors-as-built)).
- Glass variants require a translucent/blurred backdrop to read.
- Icon precedes label; the glass group provides `gap:8px`.

## Do / Don't

- **Do** use `primaryButton large` for a `grid-column:1/-1` CTA.
- **Do** give icon-only buttons an `aria-label`.
- **Don't** ship an enabled button with no handler (exists today: "Bộ lọc"/"Chia sẻ"/"Lưu"/help/auth).
- **Don't** rely on `font-weight` >800 (not loaded — falls back).

## Accessibility

- Icon-only buttons carry `aria-label`; the glyph is `aria-hidden`.
- Text buttons expose their name via visible text.
- Loading is signalled only by `disabled` (no `aria-busy`); inert buttons are neither wired nor disabled — see [`../09-accessibility.md`](../09-accessibility.md) (A5/A6).
- No custom `:focus-visible` (A1).

## Related components

- [Icon](../07-icons.md) · [Card](card.md) · [Input](input.md) · [Map](map.md) (topbar/bottom-action buttons)
- Nav-tab and close buttons: [Grid system](../05-grid-system.md), [pages/home.md](../pages/home.md).

## Files

- `apps/web/app/globals.css` — base `:505-513`; primary `:515-530`; secondary `:532-543`; glass `:761-806`; glass-action `:1298-1302, 1473-1482`.
- `apps/web/components/DiscoveryExperience.tsx` — `:535, 538, 550, 569, 589, 594, 602, 605, 625, 698, 1073-1084, 1173`.

## References

- [Color](../03-color-system.md) · [Typography](../06-typography.md) · [Accessibility](../09-accessibility.md)
- `apps/web/CLAUDE.md`.

## Future improvements

- Tokenize primary text color and per-variant shadows.
- Wire or `disabled` the inert buttons.
- Replace label-swap loading with a spinner + `aria-busy`; add `:focus-visible`.

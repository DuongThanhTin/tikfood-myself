# Color System

> **Status: As-built. Canonical color reference.** All values are literal from
> `apps/web/app/globals.css`. Other docs reference this file rather than restating colors.

## Tokens

Colors are CSS custom properties. Dark is default (`:root`); light is a full override
scoped to `.appShell[data-theme="light"]`.

| Token | Dark (`:root`, `:3-22`) | Light (`:72-91`) | Role |
|-------|------|------|------|
| `--background` | `#050505` | `#eee7dc` | Page / map background |
| `--surface` | `#131313` | `#fbf3e8` | Panels |
| `--surface-low` | `#1c1b1b` | `#eadfce` | Select inputs |
| `--surface-mid` | `#201f1f` | `#e2d4c1` | **Declared, 0 uses** (see [B3](../archive/design-system/improvement-proposals.md#b3-inconsistent-colors-relates-to-r5--r11)) |
| `--surface-high` | `#2a2a2a` | `#fff7ea` | Thumbnails, glass-solid |
| `--surface-highest` | `#353534` | `#eadac6` | Search bg, chip hover, trend track |
| `--on-surface` | `#e5e2e1` | `#251b14` | Primary text |
| `--on-muted` | `#ddc1ae` | `#6b5746` | Secondary text |
| `--outline` | `rgba(255,255,255,0.12)` | `rgba(74,55,40,0.22)` | Borders |
| `--glass` | `rgba(255,255,255,0.08)` | `rgba(250,239,224,0.78)` | Translucent surfaces |
| `--glass-strong` | `rgba(255,255,255,0.12)` | `rgba(248,233,214,0.92)` | Stronger translucent |
| `--primary` | `#ffb77f` | `#8f4700` | Brand, active, primary buttons |
| `--primary-strong` | `#ff8a00` | `#d06b00` | **Declared, 0 uses** (gradients use literal `#ff8a00`) |
| `--secondary` | `#feb700` | `#a26300` | Ratings, badges |
| `--trend-fire` | `#ff4d00` | `#d94f00` | Trend gradients, hot markers |
| `--danger` | `#ffb4ab` | `#922922` | Error text |
| `--shadow` | `0 24px 70px rgba(0,0,0,0.45)` | `0 24px 70px rgba(74,52,32,0.2)` | Shared elevation |

## Theming mechanism

- Theme is set via `data-theme` on `.appShell` (`DiscoveryExperience.tsx:398-406`),
  toggled in state (default `dark`), **not persisted**.
- Light mode re-declares the same token names plus ~15 per-surface overrides (glass,
  hero search, markers, map filters).
- `:root` is `color-scheme: dark`; light tokens apply **only inside `.appShell`** —
  anything outside stays dark.

## Non-tokenized colors (as-built)

The token set is not fully used; these literals appear directly:

- **Blue accent `#2f80ff`** (user location, route) — no token (×3 CSS + ×1 TSX).
- **`--primary` bypassed** as `rgba(255,183,127,…)` **17 times across 11 alpha levels**;
  `--secondary` (`rgba(254,183,0,…)` ×2) and `--trend-fire` (`rgba(255,77,0,…)` ×4) too.
- **`#ffffff` ×10** (no white/on-inverse token), `#4e2600` ×3 (primary-button text),
  `#2f271f`, `#2f1500`, `#080808`.
- The **map basemap palette** (`applyMapLayerPreferences`, `:145-216`) is a separate
  hardcoded color set.

## Rules

- Reference tokens, not raw hex/rgba, for any color that equals a token.
- Set `data-theme` on the shell; keep new themed UI inside `.appShell`.
- Text-on-primary is currently `#4e2600`; text-on-inverse is `#ffffff` — both are
  literals today (no token).

## Do / Don't

- **Do** use `--primary`/`--secondary`/`--trend-fire`/`--danger` for semantic color.
- **Do** provide both dark and light values when adding a token.
- **Don't** add new `rgba(255,183,127,…)` literals — that duplicates `--primary`.
- **Don't** introduce accents (like the blue) without a token.

## Accessibility

- Verify contrast when placing text on gradients (trend chip white-on-gradient) and on
  translucent glass, where the effective background varies.
- Light and dark must both meet contrast for the same component.

## Related

- [Spacing](04-spacing-system.md) · [Typography](06-typography.md) · [Motion](08-motion.md)
- Consistency findings: [improvement proposals §B3](../archive/design-system/improvement-proposals.md#b3-inconsistent-colors-relates-to-r5--r11)

## Files

- `apps/web/app/globals.css` — tokens `:3-22, 72-91`; overrides `:142-1041` (scattered).

## References

- `apps/web/CLAUDE.md` (styling rules).

## Future improvements

- Tokenize the blue accent, `--on-primary` (`#4e2600`), and a white/`--on-inverse`.
- Introduce alpha-tint tokens (or `color-mix`) so translucent primaries derive from `--primary`.
- Use or remove `--primary-strong` / `--surface-mid` (0 references today).

# Pattern: Loading

> **Status: As-built.** How the UI communicates in-flight work. There are **no
> skeletons or spinners** today — loading is label swaps + a text line.

## Purpose

Signal that a search, detail fetch, or routing request is in progress, and prevent
double-submits.

## Usage

Three independent loading flags drive the pattern:

| Flag | Signal | Where |
|------|--------|-------|
| `isLoading` (search) | button label → "Đang tìm"; primary/secondary `disabled` | `DiscoveryExperience.tsx:535-540` |
| `isDetailLoading` | `.detailStatus` line "Đang tải chi tiết quán..." | `:1115` |
| `isRouting` | action label → "Đang tải"; `glassAction` `disabled` | `:1073` |

Detail loads lazily: on select, `fetchVenueDetail` runs only if `dishes`/`opening_hours`
are missing (`:255-285`).

## Examples

```jsx
{isLoading ? "Đang tìm" : "Tìm kiếm"}                     {/* button label swap */}
{isDetailLoading ? <p className="detailStatus">Đang tải chi tiết quán...</p> : null}
```

## Rules

- Set the relevant flag before the request and clear it in `finally`.
- Disable the triggering control while loading (prevents duplicate requests).
- Loading, empty, and error are three distinct states — don't merge them.

## Do / Don't

- **Do** reflect loading on the specific control that triggered it.
- **Do** clear the flag in `finally` so errors don't strand a spinner state.
- **Don't** block the whole UI; loading is scoped per action.
- **Don't** show the empty state while a load is pending.

## Accessibility

- Loading is conveyed by visible text + `disabled` only. There is **no `aria-busy`** and
  **no `aria-live`**, so assistive tech isn't explicitly notified — a gap (A5/A4). See
  [`../09-accessibility.md`](../09-accessibility.md).

## Related

- [patterns/empty-state.md](empty-state.md) · [patterns/error-state.md](error-state.md) · [components/button.md](../components/button.md)

## Files

- `apps/web/app/globals.css` — `.statusText`/`.detailStatus` `:545-551, 1356-1361`.
- `apps/web/components/DiscoveryExperience.tsx` — flags `:235-237`; usages `:535, 1073, 1115`; lazy detail `:255-285`.

## References

- `apps/web/CLAUDE.md`.

## Future improvements

- Introduce skeleton/spinner primitives + `aria-busy`; standardize a loading component.

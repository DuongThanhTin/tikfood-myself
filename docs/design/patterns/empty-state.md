# Pattern: Empty State

> **Status: As-built.** What shows when a query/filter returns no venues.

## Purpose

Tell the user, plainly, that no venues match — without breaking the layout.

## Usage

Today the empty state is a **single line of muted text**, rendered when
`venues.length === 0`:

- Class `.emptyText` (shares styling with `.statusText`/`.errorText`,
  `globals.css:545-551`): `--on-muted`, 13px, `line-height:1.5`.
- Copy: **"Không có địa điểm phù hợp bộ lọc."**

There is no illustration, no suggested action, and no reset shortcut in the empty state
itself (Reset lives in the filter block).

## Examples

```jsx
{/* Rail empty state (:554) */}
{venues.length === 0 ? <p className="emptyText">Không có địa điểm phù hợp bộ lọc.</p> : null}
```

## Rules

- Trigger on `venues.length === 0` after a completed load (not during loading — see [loading](loading.md)).
- Keep copy short and Vietnamese.
- The empty message renders inside the rail section, above the (empty) list.

## Do / Don't

- **Do** distinguish "no results" from "loading" and from "error" (three separate states).
- **Don't** leave the area blank; always show the message.
- **Don't** conflate an empty result with an error (different copy/handling).

## Accessibility

- Plain text, read inline. It is not an `aria-live` region, so screen-reader users
  aren't notified when results become empty after a search — a gap (A4).

## Related

- [patterns/loading.md](loading.md) · [patterns/error-state.md](error-state.md) · [patterns/filter.md](filter.md)

## Files

- `apps/web/app/globals.css` — `.emptyText` `:545-551`.
- `apps/web/components/DiscoveryExperience.tsx` — `:554`.

## References

- `apps/web/CLAUDE.md`.

## Future improvements

- Add an illustrated empty state with a "reset filters" affordance and suggestions.
- Announce empties via `aria-live`.

# Pattern: Error State

> **Status: As-built.** How failures surface. Errors are **inline text**, per scope —
> there is no toast/notification component.

## Purpose

Tell the user when a search, detail load, or routing request fails, close to where it
was triggered, without losing the rest of the UI.

## Usage

Three scoped error strings, each rendered as inline text when non-empty:

| Error | Surface | Class | Where |
|-------|---------|-------|-------|
| `error` (search/location) | below the filter actions | `.errorText` (danger, bold) | `DiscoveryExperience.tsx:544` |
| `routeError` | in the detail panel | `.routeHint.error` (danger) | `:1092` |
| `detailError` | in the detail panel | `.routeHint.error` (danger) | `:1093` |

- `.errorText` (`globals.css:553-556`): `--danger`, `font-weight:800`.
- `.routeHint.error` (`:1491-1493`): `--danger` variant of the route hint.
- Messages come from the caught error's `message` (API `{ error }` envelope) or a
  Vietnamese fallback string (e.g. "Location permission was not granted.",
  "Không thể tải chỉ đường.").

## Examples

```jsx
{error ? <p className="errorText">{error}</p> : null}
{routeError ? <p className="routeHint error">{routeError}</p> : null}
```

## Rules

- Scope the error to the action that failed (search vs detail vs route); don't show a global banner.
- Prefer the API `error.message`; fall back to a clear Vietnamese sentence.
- Clear the error at the start of the next attempt (`setError("")` before a new request).
- Keep error distinct from empty and loading states.

## Do / Don't

- **Do** reset the relevant error before retrying.
- **Do** use `--danger` styling so errors read as errors.
- **Don't** use the empty-state copy for an error.
- **Don't** swallow the API message when one is provided.

## Accessibility

- Error text is `--danger`-colored and bold, but color/weight alone convey severity.
- It is **not** an `aria-live`/`role="alert"` region, so screen-reader users aren't
  proactively notified — a gap (A4). See [`../09-accessibility.md`](../09-accessibility.md).

## Related

- [patterns/loading.md](loading.md) · [patterns/empty-state.md](empty-state.md) · [components/map.md](../components/map.md) (routing errors)

## Files

- `apps/web/app/globals.css` — `.errorText` `:553-556`; `.routeHint.error` `:1491-1493`.
- `apps/web/components/DiscoveryExperience.tsx` — error state `:238-240`; usages `:544, 1092-1093`; messages `:296-297, 347, 391`.
- `apps/web/lib/api.ts` — `{ error }` envelope `:63-70`.

## References

- `apps/web/CLAUDE.md`.

## Future improvements

- Add `role="alert"`/`aria-live` for errors; consider a shared toast for transient failures.

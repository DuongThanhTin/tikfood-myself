# Playbook: Create a Form

> Build a form from the **existing** input primitives. There is no form library today
> (no React Hook Form / Zod) — forms are controlled React state + native controls. See
> [`../components/input.md`](../components/input.md).

## When to use

Grouping inputs for a task (filters already do this; a sign-in form would be the next).

## Steps

1. **Reuse the input primitives**: `.field` + native `<select>`/`<input>`, `.toggle`,
   the hero-search field; layout with `.filterGrid` (2-col → 1-col ≤640px).
2. **Controlled state**: each field is React state with `value` + `onChange` (mirrors the
   filter block). No form library unless real validation/mutations appear (`apps/web/CLAUDE.md`).
3. **Labels**: wrap each control in `<label>` (implicit association); caption-less controls
   get `aria-label`.
4. **Submit**: an explicit `primaryButton` (deliberate submit, like search); reflect
   in-flight via `disabled` + label swap ([loading](../patterns/loading.md)).
5. **Validation/errors**: surface with the [error-state](../patterns/error-state.md) pattern
   (`.errorText`); clear before re-submit.
6. **Data**: submit through `lib/api.ts`; respect the `{ data, error }` envelope. If it's a
   **mutation/auth**, get human approval first (`CLAUDE.md`).
7. **Copy** in Vietnamese.

## Checklist

- [ ] Reused `.field`/native controls/`.filterGrid`; controlled state.
- [ ] Every control labelled; submit is explicit; loading + error handled.
- [ ] Submit via `lib/api.ts`; mutations/auth approved.
- [ ] No form library added without justification.
- [ ] Meets [DoD](../../verification/definition-of-done.md).

## Related

- [../components/input.md](../components/input.md) · [../patterns/filter.md](../patterns/filter.md) · [../patterns/error-state.md](../patterns/error-state.md) · [create-modal.md](create-modal.md)

## References

- `apps/web/CLAUDE.md` (no RHF/Zod today) · `apps/web/lib/api.ts`.

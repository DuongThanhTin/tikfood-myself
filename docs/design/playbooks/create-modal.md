# Playbook: Create a Modal

> **Target is Proposed** — no modal exists yet ([`../components/modal.md`](../components/modal.md)).
> This guide covers building the first one to spec.

## When to use

A self-contained overlay task that shouldn't leave discovery (e.g. sign-in, menu
expansion). For a small confirm/prompt, use [create-dialog](create-dialog.md) instead.

## Steps

1. **Read the [modal spec](../components/modal.md)** and confirm scope is discovery-only
   (no cart/checkout/booking). Auth-bearing modals need human approval.
2. **Structure**: a conditionally-mounted scrim + centered `role="dialog"` `aria-modal="true"`
   card (reuse `--glass-strong`, `blur(18px)`, `radius:18px`, `--shadow`), a `.closeButton`,
   body, and a footer of [buttons](../components/button.md).
3. **Focus management**: trap focus while open, move initial focus into the modal, restore
   focus to the trigger on close, and make the background inert.
4. **Dismissal**: scrim click, Esc, and close button.
5. **Motion**: subtle fade/scale using the standard durations; gate on `prefers-reduced-motion`.
6. **Document** as As-built and remove the Proposed banner from `modal.md`; add it to the
   components index.

## Checklist

- [ ] Matches the spec; discovery-only; auth approved if applicable.
- [ ] `role="dialog"`/`aria-modal`; focus trap + restore; Esc + scrim + button dismiss.
- [ ] Reused glass tokens, `.closeButton`, buttons.
- [ ] `prefers-reduced-motion` respected.
- [ ] `modal.md` updated to As-built; meets [DoD](../../verification/definition-of-done.md).

## Related

- [../components/modal.md](../components/modal.md) · [create-dialog.md](create-dialog.md) · [create-form.md](create-form.md) · [../09-accessibility.md](../09-accessibility.md)

## References

- `apps/web/app/globals.css` (`.closeButton` `:1163-1177`; glass tokens).

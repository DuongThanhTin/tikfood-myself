# Playbook: Create a Dialog

> **Target is Proposed** — no dialog exists yet. In this system, a **dialog** is a small,
> focused confirm/prompt (e.g. "Xóa chỉ đường?"), distinct from the larger
> [modal](create-modal.md) overlay. Both share the same a11y contract.

## When to use

A short, blocking decision or acknowledgement — confirm/cancel, alert. For a multi-field
task or sign-in, use [create-modal](create-modal.md).

## Steps

1. **Keep it small**: a compact centered card (narrower than a modal), one line of prompt
   text, and 1–2 [buttons](../components/button.md) (`primaryButton` + `secondaryButton`).
2. **Semantics**: `role="alertdialog"` for confirmations (or `role="dialog"`),
   `aria-modal="true"`, `aria-labelledby`/`aria-describedby`.
3. **Focus**: initial focus on the safe/default action; trap while open; restore on close.
4. **Dismissal**: Esc + explicit cancel; be careful auto-dismissing destructive prompts.
5. **Reuse** the glass card tokens and `.closeButton` language; don't invent new styles.
6. **Copy**: Vietnamese, action-oriented ("Xóa", "Huỷ").
7. **Document** alongside the [modal spec](../components/modal.md) (dialog = compact variant).

## Checklist

- [ ] Compact, ≤2 actions, single prompt.
- [ ] `role="alertdialog"`/`aria-modal`; labelled + described; focus trap + restore.
- [ ] Reused glass tokens/buttons; Vietnamese copy.
- [ ] Destructive actions don't auto-dismiss unsafely.
- [ ] Meets [DoD](../../verification/definition-of-done.md).

## Related

- [create-modal.md](create-modal.md) · [../components/modal.md](../components/modal.md) · [../components/button.md](../components/button.md) · [../09-accessibility.md](../09-accessibility.md)

## References

- `apps/web/app/globals.css` (glass tokens; `.closeButton`).

# Modal

> **Status: Proposed — not implemented.** No modal exists in `apps/web` today. The spec
> below is a forward proposal grounded in the current tokens/components
> ([`../03-color-system.md`](../03-color-system.md), [`button`](button.md)); nothing here
> describes shipped code. See also [`dialog`](../playbooks/create-dialog.md) — in this
> system "modal" = a centered overlay panel, "dialog" = a focused confirm/prompt.

## Purpose

A centered overlay for a self-contained task that shouldn't navigate away from discovery
— the most likely first uses are **sign-in** and the **"Xem chi tiết & Menu"** expansion,
both currently inert buttons ([R10](../../archive/design-system/05-inconsistencies.md)).

## Usage (proposed)

- Overlay: full-viewport scrim (`rgba(5,5,5,0.5)` dark / warm-tinted light) above the
  detail panel's `z-index` (i.e. `> 7`).
- Surface: reuse the glass card language — `--glass-strong`, `blur(18px)`, `border-radius:18px`,
  `--shadow`, `max-width` ~480px, centered.
- Header (title + `.closeButton`-style close), body, footer actions using
  [`button`](button.md) (`primaryButton` + `secondaryButton`).
- Mount conditionally (like the detail panel), not CSS-hidden.

## Examples (proposed shape)

```jsx
{open ? (
  <div className="modalScrim" onClick={onClose}>
    <div className="modalCard" role="dialog" aria-modal="true" aria-labelledby="m-title" onClick={(e)=>e.stopPropagation()}>
      <button className="closeButton" aria-label="Đóng" onClick={onClose}><Icon name="close" /></button>
      <h2 id="m-title">Đăng nhập</h2>
      {/* body */}
      <div className="detailActions"><button className="primaryButton large">Tiếp tục</button></div>
    </div>
  </div>
) : null}
```

## Rules (proposed)

- Reuse existing tokens/classes; do not invent a new color/radius language.
- `role="dialog"` + `aria-modal="true"` + `aria-labelledby`; trap focus while open; restore focus on close.
- Close on scrim click, Esc, and the close button.
- Keep it discovery-scoped — no cart/checkout/booking content (anti-goals).

## Do / Don't

- **Do** build it as a conditionally-mounted overlay like `rightPanel`.
- **Do** reuse `button`, `.closeButton`, glass tokens.
- **Don't** stack multiple modals.
- **Don't** use it for transient messages (that's a toast — separate proposal).

## Accessibility (proposed)

- Focus trap + Esc close + focus restore; `aria-modal`; scrim not focusable.
- First focus goes to the modal (title or first control); background made inert.
- Respect `prefers-reduced-motion` for any enter/exit transition.

## Related components

- [Button](button.md) · [Icon](../07-icons.md) · [create-dialog playbook](../playbooks/create-dialog.md) · [create-modal playbook](../playbooks/create-modal.md)

## Files

- **None yet** (Proposed). Would live in `apps/web/components/` with styles in `app/globals.css`.

## References

- [Color](../03-color-system.md) · [Accessibility](../09-accessibility.md) · `CLAUDE.md` anti-goals.

## Future improvements

- N/A until implemented. When built, document as As-built and drop the Proposed banner.

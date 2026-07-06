# Authentication — UI

One-line purpose: the visual and interaction design for the authentication screens, expressed entirely in the existing TikFood design language.

> Status: Design — not yet implemented · Date: 2026-07-05 · Feature: authentication

---

## Purpose of this document

This document describes what the authentication screens look like and how they
behave — layout, responsive behavior, loading/error/disabled states,
accessibility, and Vietnamese copy. It reuses the canonical design system in
[`docs/design/`](../../design/) and introduces **no new design language**. It stays
at the design altitude: wireframes and state descriptions, not markup. Component
reuse decisions live in [`component-reuse.md`](component-reuse.md); the client-side
behavior lives in [`frontend.md`](frontend.md); the data contract in
[`api.md`](api.md).

---

## Existing implementation

The design language is already defined and in use by the discovery experience:

- **Aesthetic:** dark-first glassmorphism — near-black background with warm amber
  accents and translucent `backdrop-filter: blur()` surfaces. Light theme via
  `data-theme` on `.appShell`.
- **Color tokens (from `app/globals.css`):** `--primary` (#ffb77f, CTAs/active),
  `--danger` (#ffb4ab, errors), `--surface-low` (input backgrounds), `--surface`,
  `--surface-high`, `--glass` / `--glass-strong`, `--on-surface`, `--on-muted`,
  `--outline`.
- **Typography:** Plus Jakarta Sans (600/700/800) for display and headings;
  Be Vietnam Pro (400/500/700/800) for body/labels (full Vietnamese diacritics).
- **Spacing:** raw px, no scale tokens — 24px panel gutter (→16px on mobile), 8px
  dominant gap, radii 8–18px, 999px pills.
- **Existing controls/classes:** `.primaryButton` / `.primaryButton.large` /
  `.secondaryButton`, `.glassButton` / `.iconGlassButton`, `.field` (label + native
  input, grid), `.filterGrid` (2-col → 1-col ≤640px), `.toggle`, `.heroSearch`,
  `.errorText` (danger, bold, inline).
- **Auth UI today:** none. Only an inert `.authCta` "Đăng nhập ngay" placeholder
  button exists in the left panel of `DiscoveryExperience.tsx`.
- **Modal:** documented in [`docs/design/components/modal.md`](../../design/components/modal.md)
  as **Proposed / not built**. This feature uses **pages**, not modals (decision D3),
  so the modal spec is not exercised here.

→ Foundations: [`docs/design/`](../../design/) —
[`01-design-principles.md`](../../design/01-design-principles.md),
[`02-brand.md`](../../design/02-brand.md),
[`03-color-system.md`](../../design/03-color-system.md),
[`06-typography.md`](../../design/06-typography.md).

---

## Proposed changes

Three new pages plus a Google callback state, each a centered glass card built only
from existing tokens/classes. Copy is Vietnamese.

### Login page (`/login`)

```
┌──────────────────────────────────────────────┐
│                  ● TikFood                     │  ← wordmark (Plus Jakarta 800, --primary)
│                                                │
│   ┌────────────── glass card ─────────────┐   │
│   │  Đăng nhập                             │   │  ← h1
│   │  Chào mừng trở lại 👋                  │   │  ← subtitle (--on-muted)
│   │                                        │   │
│   │  Email                                 │   │  ← .field (label + native input)
│   │  [ ____________________________ ]      │   │
│   │  Mật khẩu                     [👁]      │   │  ← optional show/hide (OQ)
│   │  [ ____________________________ ]      │   │
│   │  ⚠ Email hoặc mật khẩu không đúng.     │   │  ← .errorText (only on error)
│   │                                        │   │
│   │  [        Đăng nhập        ]           │   │  ← .primaryButton.large
│   │  ──────────  hoặc  ──────────          │   │  ← divider
│   │  [   G   Đăng nhập với Google   ]      │   │  ← GoogleSignInButton (glass variant)
│   │                                        │   │
│   │  Chưa có tài khoản? Đăng ký            │   │  ← link → /register
│   │  Quên mật khẩu?                        │   │  ← link → /forgot-password
│   └────────────────────────────────────────┘   │
└──────────────────────────────────────────────┘
```

### Register page (`/register`)

Same card language, with an added display-name field. On wide screens the Email /
Mật khẩu row may use the existing `.filterGrid` two-column pattern.

```
┌────────────── glass card ─────────────┐
│  Đăng ký                              │
│  Tạo tài khoản TikFood                │
│                                        │
│  Tên hiển thị                          │  ← .field
│  [ ____________________________ ]      │
│  Email                                 │
│  [ ____________________________ ]      │
│  Mật khẩu (tối thiểu 8 ký tự)          │
│  [ ____________________________ ]      │
│  ⚠ <inline field error if any>         │  ← .errorText
│                                        │
│  [        Tạo tài khoản        ]       │  ← .primaryButton.large
│  ──────────  hoặc  ──────────          │
│  [   G   Đăng ký với Google   ]        │
│                                        │
│  Đã có tài khoản? Đăng nhập            │  ← link → /login
└────────────────────────────────────────┘
```

### Forgot-password placeholder (`/forgot-password`)

Clearly marked as a future feature; no working submit.

```
┌────────────── glass card ─────────────┐
│  Quên mật khẩu                        │
│  Tính năng đang được phát triển.       │  ← "future" notice
│  Vui lòng quay lại đăng nhập.          │
│  [   ← Quay lại đăng nhập   ]          │  ← .secondaryButton → /login
└────────────────────────────────────────┘
```

### Google callback state (`/auth/google/callback`)

Transient screen while the session is established, with an error fallback.

```
  Loading:  [ ◐ ] Đang hoàn tất đăng nhập…
  Error:    ⚠ Đăng nhập với Google thất bại. [ Thử lại ] → /login
```

### Responsive behavior

| Breakpoint | Layout |
|---|---|
| ≤640px (e.g. 375) | Single column; card full-width with 16px gutter; fields stacked |
| 641–1024px (e.g. 768) | Centered card ~440px; 24px gutter |
| ≥1025px (e.g. 1280) | Centered card ~440px on the dark glass background; register may use `.filterGrid` 2-col |

Reuse the existing `.filterGrid` 2-col → 1-col ≤640px rule; do not invent new
breakpoints.

### State visuals

- **Loading:** submit button shows `disabled` state with label swap
  ("Đăng nhập" → "Đang đăng nhập…"); inputs remain readable. No new spinner styles
  beyond existing conventions.
- **Error:** inline `.errorText` (danger, bold) scoped to the form; API
  `error.message` preferred with a Vietnamese fallback; error cleared before retry.
- **Disabled:** submit disabled until required fields are non-empty (and while
  submitting).

### Google button

Rendered within an **existing** button variant (glass/secondary) with a small "G"
mark and Vietnamese label — never a foreign branded widget that breaks the design
language. Follow Google brand guidance for the mark within our token palette.
Details in [`component-reuse.md`](component-reuse.md).

### Accessibility

- Every input wrapped in a `<label>` (via `.field`); labels are visible, not
  placeholder-only.
- Icon-only controls (e.g. password show/hide, close/back) carry an `aria-label`.
- Inline errors associated to inputs via `aria-describedby`; focus moves to the
  first invalid field on submit.
- `aria-live` on the error region is recommended to fix the known gap A4 for auth;
  see Open questions. Verify both dark and light `data-theme`.

### Vietnamese copy (strings)

| Key | Vietnamese |
|---|---|
| Login title | Đăng nhập |
| Register title | Đăng ký |
| Email label | Email |
| Password label | Mật khẩu |
| Display name label | Tên hiển thị |
| Password hint | Mật khẩu (tối thiểu 8 ký tự) |
| Submit login | Đăng nhập / Đang đăng nhập… |
| Submit register | Tạo tài khoản / Đang tạo… |
| Google login | Đăng nhập với Google |
| Divider | hoặc |
| To register | Chưa có tài khoản? Đăng ký |
| To login | Đã có tài khoản? Đăng nhập |
| Forgot link | Quên mật khẩu? |
| Forgot placeholder | Tính năng đang được phát triển. |
| Invalid credentials | Email hoặc mật khẩu không đúng. |
| Invalid email | Email không hợp lệ. |
| Weak password | Mật khẩu phải có ít nhất 8 ký tự. |
| Email taken | Email đã được sử dụng. |
| Google failure | Đăng nhập với Google thất bại. |
| Callback loading | Đang hoàn tất đăng nhập… |

---

## Open questions

1. **Card vs split-screen:** single centered glass card (proposed) or a two-pane
   layout (brand imagery + form) on wide screens — which better fits the map-first
   brand without adding a new pattern?
2. **Password visibility toggle:** include a show/hide control? If so it needs an
   `aria-label` and a reused icon-button variant.
3. **Validation timing:** validate per-blur, on-submit only, or hybrid? (Discovery
   search uses deliberate on-submit.)
4. **Logged-in state in the discovery shell:** exact treatment of the `.authCta`
   region when authenticated (name + logout link vs compact menu) — must stay within
   the existing language. See [`frontend.md`](frontend.md).
5. **Light-theme parity:** confirm the auth card is validated in both dark and light
   `data-theme` modes.
6. **`aria-live` adoption:** adopt for auth now (fixing gap A4 locally) or defer to a
   system-wide accessibility pass?

---

## Related repository documentation

- [`docs/design/`](../../design/) — canonical design system (as-built).
  - [`01-design-principles.md`](../../design/01-design-principles.md) · [`02-brand.md`](../../design/02-brand.md) · [`03-color-system.md`](../../design/03-color-system.md) · [`06-typography.md`](../../design/06-typography.md)
  - [`components/button.md`](../../design/components/button.md) · [`components/input.md`](../../design/components/input.md) · [`components/modal.md`](../../design/components/modal.md) (Proposed)
  - [`patterns/error-state.md`](../../design/patterns/error-state.md) · [`playbooks/create-form.md`](../../design/playbooks/create-form.md)
- [`apps/web/CLAUDE.md`](../../../apps/web/CLAUDE.md) — plain-CSS rule, reuse `globals.css`, Vietnamese copy, `aria-label` on icon buttons.
- [`/CLAUDE.md`](../../../CLAUDE.md) — naming, "reuse existing code" before finishing.
- Sibling docs: [`ui.md`](ui.md) ↔ [`component-reuse.md`](component-reuse.md), [`frontend.md`](frontend.md), [`api.md`](api.md), [`specification.md`](specification.md).

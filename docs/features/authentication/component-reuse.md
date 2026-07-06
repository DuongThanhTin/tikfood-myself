# Authentication — Component Reuse

One-line purpose: the reuse-before-create audit that keeps the auth UI inside the existing design language.

> Status: Design — not yet implemented · Date: 2026-07-05 · Feature: authentication

---

## Purpose of this document

This document enumerates what auth UI needs, maps each need to an **existing** token
or CSS class to reuse, and justifies the few genuinely new components. It enforces
the rule "reuse before create; introduce no new design language" from the feature
brief and [`/CLAUDE.md`](../../../CLAUDE.md). Visual design is in [`ui.md`](ui.md);
client behavior is in [`frontend.md`](frontend.md).

---

## Existing implementation

`apps/web` has **no shared React component library** today — reusable UI exists only
as CSS classes in `app/globals.css`, consumed inline by `DiscoveryExperience.tsx`
and `VenueList.tsx`. Relevant existing assets:

| Asset (CSS class) | Role today |
|---|---|
| `.primaryButton`, `.primaryButton.large` | Primary CTA; `.large` spans full width |
| `.secondaryButton` | Companion/cancel action |
| `.glassButton`, `.glassButton.solid`, `.iconGlassButton` | Translucent actions over the map; icon-only (needs `aria-label`) |
| `.field` | Labelled control: `<label>` + native `<input>`/`<select>`, grid layout |
| `.filterGrid` | 2-column form layout → 1-column ≤640px |
| `.toggle` | Native checkbox with `accent-color: var(--primary)` |
| `.heroSearch` | Large search field (not reused for auth) |
| `.errorText` | Inline error text (danger, bold) |
| Color/type/spacing tokens | `--primary`, `--danger`, `--surface-low`, Plus Jakarta / Be Vietnam, 24/16px gutters |

No form library (react-hook-form/Zod) and no React Query exist, per
[`apps/web/CLAUDE.md`](../../../apps/web/CLAUDE.md).

---

## Proposed changes

### Reuse-vs-create map

| Auth UI need | Decision | Reuse / Create | Why |
|---|---|---|---|
| Primary submit ("Đăng nhập"/"Tạo tài khoản") | **Reuse** | `.primaryButton.large` | Exact CTA pattern already exists |
| Back / secondary action (forgot-password) | **Reuse** | `.secondaryButton` | Matches companion-action language |
| Text/email/password input + label + inline error | **Create thin wrapper** | `FormField` **built on `.field` + `.errorText`** | Avoids repeating label/input/error markup across ≥6 fields; adds `aria-describedby` wiring; visually identical to `.field` |
| Two-column field layout (register, wide screens) | **Reuse** | `.filterGrid` | Existing responsive 2→1 col rule |
| Error message display | **Reuse** | `.errorText` | Canonical inline error pattern |
| Loading state on submit | **Reuse pattern** | `disabled` + label swap | Established convention (search button) |
| Google sign-in button | **Create thin wrapper** | `GoogleSignInButton` **on `.glassButton`/`.secondaryButton` variant** | Keeps Google mark within our token palette + Google brand rules; no foreign widget |
| Form shell (state, submit, error mapping) | **Create** | `AuthForm` | Shared controlled-state container for login/register; not a visual element, no new design tokens |
| "hoặc" divider | **Reuse or one-off CSS** | small CSS rule (candidate component) | Trivial; see Open questions |
| Card container (`.authCard`) / page wrapper (`.authPage`) | **Create CSS only** | new classes composed from existing glass/surface tokens | Layout containers; reuse `--glass-strong`, `--surface`, radii, `--shadow` — no new visual language |
| Wire the inert `.authCta` | **Reuse + wire** | existing `.authCta` element | Fixes accessibility gap A6; links to `/login` / shows logout when authenticated |

### What is reused vs newly created

- **Reused as-is:** `.primaryButton(.large)`, `.secondaryButton`, `.field`,
  `.filterGrid`, `.errorText`, `.glassButton`, all color/type/spacing tokens, the
  loading (disabled + label-swap) and error (inline `.errorText`) conventions.
- **Newly created (thin, no new design language):** `FormField`, `AuthForm`,
  `GoogleSignInButton` React components; `.authPage` / `.authCard` CSS composed from
  existing tokens; a small "hoặc" divider rule.

### Guiding rule

Reuse before create. New components are permitted only to remove duplication or wire
behavior — never to introduce new colors, fonts, spacing, or a competing visual
language. Every new class is composed from existing tokens. A design self-review
against [`docs/design/`](../../design/) runs before the frontend PRs merge.

### Promotion candidates (later, not now)

If these prove reusable beyond auth, promote them into the canonical design system
in a follow-up (not part of the feature implementation):

- `FormField` → generalize as a documented input primitive in
  [`docs/design/components/input.md`](../../design/components/input.md).
- An **auth-form / create-account playbook** under
  [`docs/design/playbooks/`](../../design/playbooks/) (extends `create-form.md`).
- Promote [`components/modal.md`](../../design/components/modal.md) from Proposed to
  as-built **only if** a modal is later introduced (this feature uses pages).

These are tracked in [`tasks.md`](tasks.md) M12 (docs & design-system
reconciliation) and gated by an actual reuse audit.

---

## Open questions

1. Promote `FormField` to a general primitive now, or keep it auth-scoped until a
   second consumer appears?
2. Do `.authPage` / `.authCard` belong in `globals.css` (single-file convention) or a
   scoped stylesheet?
3. Google mark asset: inline SVG "G" vs an existing icon glyph — which stays cleanest
   within the icon system and Google brand rules?
4. Divider ("hoặc"): one-off CSS rule vs a small reusable component — is it used
   elsewhere enough to warrant a component?

---

## Related repository documentation

- [`docs/design/`](../../design/) — design system; especially
  [`components/button.md`](../../design/components/button.md),
  [`components/input.md`](../../design/components/input.md),
  [`patterns/error-state.md`](../../design/patterns/error-state.md),
  [`playbooks/create-form.md`](../../design/playbooks/create-form.md).
- [`apps/web/CLAUDE.md`](../../../apps/web/CLAUDE.md) — reuse `globals.css`, no new CSS
  framework, no form library yet.
- [`/CLAUDE.md`](../../../CLAUDE.md) — "reuse existing code; search before adding a new
  helper/service/component."
- Sibling docs: [`ui.md`](ui.md), [`frontend.md`](frontend.md), [`tasks.md`](tasks.md).

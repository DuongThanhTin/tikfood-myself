# Dependency Management Standard

Dependencies must be intentional. A small dependency surface keeps the repo easier for humans and AI agents to reason about.

## General Rules

- Prefer standard library and existing local helpers before adding a package.
- Add a dependency only when it removes real complexity or provides a mature, well-known capability.
- Do not add packages for trivial formatting, string manipulation, or one-off helpers.
- Do not add a package that introduces delivery, cart, order, payment, booking, reservation, chat, monetization, or livestream scope.
- Update lockfiles when dependencies change.
- Run relevant tests after dependency changes.
- Mention new dependencies in the implementation summary.

## Root npm Workspace

Root `package.json` owns shared npm workspace scripts.

Allowed npm workspaces:

```text
apps/ai-code-runner
apps/web
```

When adding npm packages:

```bash
npm install <package> --workspace apps/web
npm install <package> --workspace apps/ai-code-runner
```

Do not manually edit `package-lock.json` except through npm.

## Frontend Dependencies

Current frontend stack:

- Next.js App Router
- React
- TypeScript

Preferred future frontend additions:

- Tailwind CSS for styling
- shadcn/ui for accessible UI primitives
- MapLibre GL or Mapbox GL for map rendering
- TanStack Query for client-side server state

Rules:

- Add UI dependencies only when a real component needs them.
- Add map dependencies only when implementing actual interactive map behavior.
- Keep API client types in `apps/web/lib`.
- Do not add state libraries until local React state and URL state are clearly insufficient.

## Backend Dependencies

Current backend stack:

- Go
- Standard library `net/http`

Preferred future backend additions:

- `chi` or `Gin` only if routing/middleware complexity justifies it
- PostgreSQL driver when real persistence is added
- PostGIS support through SQL queries or a mature driver layer

Rules:

- Keep MVP API on `net/http` unless route/middleware complexity grows.
- Do not add an ORM by default.
- Prefer explicit SQL for early PostgreSQL/PostGIS work.
- Keep domain logic independent from HTTP framework choices.

## Security Review For Dependencies

Before adding a dependency, check:

- Is it actively maintained?
- Does it add network behavior?
- Does it require credentials?
- Does it change build or deploy behavior?
- Does it touch protected paths?
- Does it increase bundle or runtime size meaningfully?

If yes to any sensitive area, call it out in the PR summary.

## Audit Policy

Run:

```bash
npm audit --audit-level=moderate
```

Do not run `npm audit fix --force` blindly. If a forced fix downgrades or breaks framework versions, document the risk and choose a targeted upgrade path.

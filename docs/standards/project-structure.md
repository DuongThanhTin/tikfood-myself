# Project Structure Standard

This repo is a monorepo. Keep code organized by product surface and responsibility.

## Root Layout

```text
apps/
  api/                  Go backend
  web/                  Next.js frontend
  ai-code-runner/       AI automation runner

packages/
  config/               Agent and product config
  prompts/              Runtime AI prompts
  schemas/              OpenAPI and JSON Schemas

docs/
  standards/            Engineering standards
  tikfood/              Product scope and anti-goals
  agents/               Human-readable agent docs

workflows/
  n8n/                  n8n workflow docs and examples

examples/
  feature-requests/     Example feature request JSON
```

## Backend Structure

Use `apps/api`.

```text
apps/api/
  cmd/server/           Process entrypoint
  internal/http/        HTTP routes, handlers, request parsing, response writing
  internal/discovery/   Discovery domain logic
```

Future backend modules should follow:

```text
internal/<domain>/      Domain/service logic
internal/http/          Transport layer only
internal/storage/       Database access when PostgreSQL/PostGIS is added
internal/config/        Runtime config parsing
internal/logging/       Logger setup
```

Do not put business logic directly inside HTTP handlers except thin request/response orchestration.

## Frontend Structure

Use `apps/web`.

```text
apps/web/
  app/                  Next.js App Router routes and layouts
  components/           Reusable UI components
  lib/                  API client, shared helpers, typed data access
  public/               Static assets
```

Future feature modules may use:

```text
features/discovery/
features/venue/
features/search/
```

Use feature folders only when a feature has enough state, components, and data access to justify separation.

## Automation Structure

Use `apps/ai-code-runner`.

```text
src/jobs/               Job orchestration
src/tools/              File, Git, command, context utilities
src/guards/             Security and policy guards
src/utils/              Small pure helpers
```

Runner code must remain conservative and explicit. Do not hide Git, file, or command side effects behind vague abstractions.

## Placement Rules

- New backend endpoints go in `apps/api/internal/http`.
- New backend domain behavior goes in `apps/api/internal/<domain>`.
- New frontend route pages go in `apps/web/app`.
- New frontend API calls go in `apps/web/lib`.
- New reusable frontend components go in `apps/web/components`.
- New AI prompts go in `packages/prompts`.
- New API schemas go in `packages/schemas`.
- New standards go in `docs/standards`.

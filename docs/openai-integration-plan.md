# OpenAI Integration Plan

OpenAI API usage is deferred for now to avoid API cost.

The current runner stops at `stage: "model"` after:

- Validating the feature request
- Blocking obvious MVP anti-goal requests
- Cloning the repository
- Creating a local `ai/*` branch
- Reading safe repo context

## When To Enable

Enable this only when:

- A new OpenAI API key is available.
- The key is stored in local `.env` or a secret manager, never committed.
- Budget limits and usage monitoring are defined.
- Git branch and protected-path guards are already working.

## Step 1: Planning Only

Do not let the model edit files first.

Flow:

```text
feature request
-> repo context summary
-> planning prompt
-> OpenAI
-> JSON plan
-> validate plan
```

Expected JSON shape:

```json
{
  "success": true,
  "summary": "string",
  "files_to_read": [],
  "files_to_change": [],
  "commands_to_run": [],
  "risks": [],
  "requires_human_approval": false
}
```

Block the plan if it includes:

- Secrets or `.env` paths
- Protected paths without approval
- Push to `main` or `master`
- Delivery, cart, order, payment, checkout, booking, reservation, chat, monetization, or livestream scope

## Step 2: Editing JSON

After a valid plan, ask for edit instructions as JSON.

```json
{
  "success": true,
  "edits": [
    {
      "path": "apps/api/internal/http/router.go",
      "operation": "replace",
      "content": "full file content"
    }
  ],
  "commands_to_run": ["npm test"]
}
```

The runner, not the model, applies edits.

## Step 3: Verification

Run allowlisted commands only:

- `npm test`
- `npm run build`
- `npm run api:test`
- `npm run web:typecheck`
- `go test ./...`

If verification fails, allow a small fixed number of model repair attempts.

## Step 4: Review

Send `git diff`, check results, feature request, and repo context to the reviewer prompt.

Commit and push only if reviewer approves.

## Step 5: Commit And Push

Only push branches under `ai/`.

Never push `main` or `master`.

Return success JSON matching `packages/schemas/runner-success-response.schema.json` so n8n can create the PR.

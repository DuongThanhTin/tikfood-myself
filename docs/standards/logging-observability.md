# Logging And Observability Standard

Logging should make debugging easier without leaking sensitive data.

## General Rules

- Logs must not include secrets, tokens, API keys, private keys, `.env` values, or credentials.
- Logs should include enough context to debug request flow.
- Prefer structured fields over large unstructured blobs.
- Avoid logging full request or response bodies unless explicitly safe and small.
- Never log user private data unless the feature has an approved privacy policy.

## Backend Logging

Current backend uses minimal startup logging.

Future backend logs should include:

- request id
- method
- path
- status code
- duration
- domain operation
- non-sensitive error reason

Recommended shape:

```json
{
  "level": "info",
  "service": "api",
  "request_id": "req_123",
  "method": "GET",
  "path": "/api/v1/map/venues",
  "status": 200,
  "duration_ms": 12
}
```

Do not log:

- `.env`
- Authorization headers
- cookies
- API keys
- raw SQL with sensitive values
- full social ingestion payloads unless approved

## Frontend Logging

Frontend should keep console logging minimal.

Allowed:

- Development-only debugging
- Non-sensitive error summaries

Not allowed:

- Tokens
- Raw API keys
- Secret env values
- Full user/session payloads

User-facing errors should be readable and calm. Do not expose stack traces in the UI.

## Runner Logging

`ai-code-runner` logs and responses must be especially careful because they can include repo context.

Allowed:

- workspace path
- repo path
- branch
- base commit
- safe files read
- commands run
- test result summaries

Not allowed:

- contents of `.env`
- secret-like file contents
- credentials
- private keys
- model prompts containing secrets
- full command output if it includes secret-like values

## Metrics Direction

Future metrics:

- request count
- request latency
- API error count
- runner job count
- runner job stage failures
- test pass/fail counts
- AI token usage when OpenAI integration is enabled

Do not add metrics dependencies until the app needs operational monitoring beyond local MVP.

## Audit Direction

Runner audit logs should record:

- feature id
- repo
- branch
- commit SHA
- safe files read
- commands run
- reviewer result
- PR URL

Audit logs must not include secrets.

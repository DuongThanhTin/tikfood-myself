# AI Feature To PR Workflow

## Flow

```text
Manual Trigger
-> Set Feature Request
-> HTTP Request: POST /jobs/feature
-> IF success
-> GitHub: Create Pull Request
-> Notification
```

## Feature Request Source

Use payloads from `examples/feature-requests/`.

Recommended first test:

- `examples/feature-requests/tikfood-realtime-map.json`

## HTTP Request Configuration

- Method: `POST`
- URL inside Docker: `http://ai-code-runner:8080/jobs/feature`
- URL from host: `http://localhost:8080/jobs/feature`
- Send body: enabled
- Body content type: JSON

## GitHub PR Configuration

Only run this node when `success` is true.

- Repository: parse or configure from `repo`
- Base branch: request `base_branch`
- Head branch: response `branch`
- Title: response `pr_title`
- Body: response `pr_body`

## Failure Behavior

If `success` is false, send the runner failure response to a notification channel and require human review. Do not create a PR.

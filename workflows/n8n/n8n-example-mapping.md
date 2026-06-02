# n8n Example Mapping

Map Set node fields directly into the runner request:

| n8n field | Runner field |
| --- | --- |
| `feature_id` | `feature_id` |
| `repo` | `repo` |
| `base_branch` | `base_branch` |
| `title` | `title` |
| `description` | `description` |
| `acceptance_criteria` | `acceptance_criteria` |
| `mode` | `mode` |

Map runner success output into the GitHub PR node:

| Runner field | GitHub PR field |
| --- | --- |
| `repo` | Repository |
| `base_branch` | Base branch |
| `branch` | Head branch |
| `pr_title` | Title |
| `pr_body` | Body |

Map runner failure output into notifications:

| Runner field | Notification usage |
| --- | --- |
| `feature_id` | Feature identifier |
| `stage` | Failed stage |
| `error` | Main failure message |
| `logs` | Trimmed logs |
| `recommendation` | Next step |

# Usage

## Basic CLI Usage

Run octoslash with a GitHub webhook event:

```bash
octoslash --event-name=issue_comment --event-path=./event.json
```

## Environment Variables

Octoslash uses the following environment variables:

- `GITHUB_TOKEN`: GitHub Personal Access Token or GitHub App token
- `GITHUB_EVENT_NAME`: GitHub event name (automatically set in GitHub Actions)
- `GITHUB_EVENT_PATH`: Path to GitHub event JSON file (automatically set in GitHub Actions)

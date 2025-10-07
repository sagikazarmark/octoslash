# Octoslash GitHub Action

This is a composite GitHub Action that downloads and runs the [octoslash](https://github.com/sagikazarmark/octoslash) binary with a GitHub token for authentication.

## Usage

### Basic Usage

```yaml
- name: Run Octoslash
  uses: ./action  # or sagikazarmark/octoslash/action if using from the repository
  with:
    github-token: ${{ secrets.GITHUB_TOKEN }}
```

### Specify Version

```yaml
- name: Run Octoslash
  uses: ./action
  with:
    version: 'v0.0.2'
    github-token: ${{ secrets.GITHUB_TOKEN }}
```

### Use Latest Version (Default)

```yaml
- name: Run Octoslash
  uses: ./action
  with:
    version: 'latest'  # This is the default
    github-token: ${{ secrets.GITHUB_TOKEN }}
```

## Inputs

| Name | Description | Required | Default |
|------|-------------|----------|---------|
| `version` | Version of octoslash to download (e.g., `v0.0.2` or `latest`) | No | `latest` |
| `github-token` | GitHub token for authentication | Yes | - |

## Supported Platforms

This action supports Linux runners with the following architectures:
- x86_64 (amd64)
- arm64 (aarch64)
- i386 (i686)

## How it Works

1. **Version Resolution**: If `version` is set to `latest`, the action queries the GitHub API to get the latest release version. Otherwise, it uses the specified version.

2. **Architecture Detection**: The action detects the runner's architecture using `uname -m` and maps it to the appropriate binary variant.

3. **Download**: Downloads the appropriate octoslash binary archive from the GitHub releases.

4. **Extraction**: Extracts the binary from the tar.gz archive and makes it executable.

5. **Execution**: Runs octoslash with the provided GitHub token set as the `GITHUB_TOKEN` environment variable.

## Example Workflow

```yaml
name: Run Octoslash

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  octoslash:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4

      - name: Run Octoslash
        uses: ./action
        with:
          version: 'latest'
          github-token: ${{ secrets.GITHUB_TOKEN }}
```

## Requirements

- Linux runner (GitHub Actions Ubuntu runners work)
- GitHub token with appropriate permissions for octoslash operations
- Internet access to download the binary from GitHub releases

## Troubleshooting

### Unsupported Architecture

If you see an error about unsupported architecture, ensure you're running on a Linux runner with one of the supported architectures (x86_64, arm64, i386).

### Download Failures

If the download fails, check:
- The specified version exists in the [releases](https://github.com/sagikazarmark/octoslash/releases)
- The runner has internet access
- GitHub API rate limits haven't been exceeded

### Permission Issues

Ensure the provided GitHub token has the necessary permissions for octoslash to perform its operations.
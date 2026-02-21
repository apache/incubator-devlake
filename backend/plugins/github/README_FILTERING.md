# GitHub Plugin - Bot Filtering

## Overview

The GitHub plugin supports filtering bot-generated PRs, reviews, and comments from data collection to prevent them from skewing metrics like lead time for changes and PR pickup time.

## Configuration

Set the `GITHUB_PR_EXCLUDELIST` environment variable with a comma-separated list of bot usernames to exclude:

```bash
export GITHUB_PR_EXCLUDELIST="renovate[bot],dependabot[bot],github-actions[bot]"
```

## Common Bot Usernames

- `renovate[bot]` - Renovate dependency updates
- `dependabot[bot]` - GitHub Dependabot
- `github-actions[bot]` - GitHub Actions automated PRs
- `sonarcloud[bot]` - SonarCloud code analysis
- `codecov[bot]` - Codecov coverage reports

## What Gets Filtered

When a username is in the exclusion list, the following entities are NOT collected:

1. **Pull Requests** - PRs authored by bots
2. **PR Reviews** - Reviews submitted by bots
3. **PR Review Comments** - Comments on PR reviews by bots
4. **Issue Comments** - Comments on issues by bots

## How It Works

- Filtering happens at the **extraction** layer
- Raw API responses are still saved (in `_raw_github_api_*` tables)
- Filtered entities never reach the tool layer tables
- Metrics queries only see non-bot entities

## Matching Rules

- **Case-insensitive**: `renovate[bot]` matches `Renovate[bot]` and `RENOVATE[BOT]`
- **Exact match**: Must match the full username
- **Whitespace trimmed**: Extra spaces in the config are ignored

## Examples

### Docker Compose

```yaml
services:
  devlake:
    environment:
      - GITHUB_PR_EXCLUDELIST=renovate[bot],dependabot[bot]
```

### Kubernetes

```yaml
env:
  - name: GITHUB_PR_EXCLUDELIST
    value: "renovate[bot],dependabot[bot],github-actions[bot]"
```

### Local Development

```bash
# .env file
GITHUB_PR_EXCLUDELIST=renovate[bot],dependabot[bot]
```

## Updating the Exclusion List

Changes to `GITHUB_PR_EXCLUDELIST` require a DevLake restart. After updating:

1. Restart DevLake
2. Trigger re-collection for affected repositories
3. Previously collected bot data remains in the database
4. New collections will respect the updated filter

## Verification

Check logs for filtering activity:

```
DEBUG: Skipping PR #123 from bot user: renovate[bot]
DEBUG: Skipping review #456 from bot user: dependabot[bot]
```

## Troubleshooting

**Bot PRs still appearing in metrics:**

1. Verify `GITHUB_PR_EXCLUDELIST` is set correctly
2. Check DevLake logs for "Skipping" messages
3. Ensure username matches exactly (case-insensitive)
4. Restart DevLake after config changes
5. Re-run collection for the repository

**How to find bot usernames:**

Check GitHub PR/comment authors in the web UI - bot usernames typically end with `[bot]`.

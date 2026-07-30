<!--
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->
# Cursor Plugin (Usage & Cost)

This plugin ingests **Cursor team usage, billing, and adoption metrics** from the [Cursor Admin API](https://cursor.com/docs/account/teams/admin-api) and stores them in DevLake tool-layer tables for Grafana dashboards and SQL analysis.

It follows the same structure and patterns as other DevLake AI usage plugins (notably `backend/plugins/gh-copilot`).

## What it collects

**Cursor Admin API endpoints:**

| Endpoint | Method | Data |
|----------|--------|------|
| `/teams/members` | GET | Team roster |
| `/teams/spend` | POST | Per-user billing cycle spend |
| `/teams/filtered-usage-events` | POST | Event-level usage and charges |
| `/teams/daily-usage-data` | POST | Per-user per-day adoption metrics |

**Stored data (tool layer):**

| Table | Description |
|-------|-------------|
| `_tool_cursor_members` | Team member roster (email, name, role) |
| `_tool_cursor_usage_events` | Billable usage events with model, tokens, and charged amounts |
| `_tool_cursor_user_spend` | Per-user spend for the current billing cycle (on-demand and included) |
| `_tool_cursor_daily_usage` | Daily adoption metrics: completions, requests by feature, tab acceptance, line edits |

Data is collected in the **Raw → Tool** layers only. There is no domain-layer converter in this plugin; Grafana dashboards query `_tool_cursor_*` tables directly.

## Data flow

```mermaid
flowchart LR
  API[Cursor Admin API]
  RAW[(Raw tables\n_raw_cursor_*)]
  TOOL[(Tool tables\n_tool_cursor_*)]
  GRAF[Grafana Dashboards]

  API --> RAW --> TOOL --> GRAF
```

**Pipeline subtasks (in order):**

1. `collectMembers` → `extractMembers`
2. `collectUsageEvents` → `extractUsageEvents`
3. `collectUserSpend` → `extractUserSpend`
4. `collectDailyUsage` → `extractDailyUsage`

## Repository layout

- `api/` — REST layer for connections, scopes, and scope configs
- `impl/` — plugin meta, blueprint v200, connection helpers
- `models/` — tool-layer models and migration scripts
- `tasks/` — collectors, extractors, and pipeline registration
- `service/` — connection test logic (Admin API permission probes)
- `e2e/` — E2E fixtures and golden CSV assertions

## Setup

### Prerequisites

- A Cursor **Team or Business** plan with Admin API access
- A **Team Admin API key** from the Cursor dashboard (Dashboard → API Keys)
- Do **not** use a User API key from Settings → Integrations — user keys cannot access `/teams/*` endpoints

Authentication uses HTTP Basic auth: the API key as username with an empty password.

### 1) Create a connection

1. DevLake UI → **Connections → Add Connection → Cursor**
2. Fill in:
   - **Name**: e.g. `Cursor Production Team`
   - **Endpoint**: defaults to `https://api.cursor.com`
   - **Token**: Team Admin API key
   - **Rate Limit**: defaults to 1,200 requests/hour (Cursor documents 20 requests/minute)
3. Click **Test Connection**. DevLake probes `/teams/members`, `/teams/spend`, and `/teams/filtered-usage-events` and reports which endpoints the key can access.
4. Save the connection.

When updating an existing connection, omit the token field to keep the encrypted value already stored in DevLake.

### 2) Add a scope

Cursor data is team-level. Add a **Team** scope for the connection. The default scope ID is `team`.

### 3) Create a blueprint

Use a blueprint plan like:

```json
[
  [
    {
      "plugin": "cursor",
      "options": {
        "connectionId": 1,
        "scopeId": "team"
      }
    }
  ]
]
```

Run the blueprint on a daily schedule to keep usage and cost data current.

### Collection behavior

- **Initial backfill**: usage events and daily usage collect up to **90 days** of history on the first run.
- **Incremental runs**: subsequent runs use the pipeline sync policy / collector state (`LatestSuccessStart`) and **rewind 7 calendar days** so recently-missed or partial days are re-fetched (same idea as the gh-copilot report lookback).
- **Date range chunking**: both `/teams/daily-usage-data` and `/teams/filtered-usage-events` requests are split into **30-day** chunks (API limit for daily usage; applied to usage events for resilience).
- **Extract**: `extractUsageEvents` and `extractDailyUsage` use a cursor-local stateful extractor (incremental by default). Incremental runs also promote any raw rows with `id > MAX(_raw_data_id)` already in the tool table, so collected-but-unpromoted raw data is healed without a full refresh. A config version bump (`extractorVersion`) triggers a one-time full re-extract after upgrade.
- **Dashboards**: query `_tool_cursor_*` tables only. Legacy domain tables (`cursor_usage`, `cursor_team_events`) are unrelated and must not be used.
- **Rate limiting**: collectors honor `Retry-After` response headers and respect the configured `rateLimitPerHour`.

## Dashboards

Grafana dashboard JSON lives under `grafana/dashboards/mysql/`:

| Dashboard | File | UID |
|-----------|------|-----|
| Cursor Usage & Cost | `cursor-usage.json` | `cursor_usage` |
| AI Cost Efficiency (Cursor panels) | `ai-cost-efficiency.json` | — |
| Multi-AI Comparison (Cursor panels) | `multi-ai-comparison.json` | — |

## Error handling

| Symptom | Likely cause |
|---------|--------------|
| **401 Unauthorized** on test connection | Invalid API key, or a User API key instead of a Team Admin key |
| **403 Forbidden** on spend or usage events | Key lacks permission for that Admin API endpoint |
| **429 Too Many Requests** | Rate limit exceeded — lower `rateLimitPerHour` or wait for `Retry-After` |
| Empty usage events after successful run | Selected time range has no billable events, or sync policy excludes the date range |
| Tool max date stuck while raw advances | Deploy plugin with lookback + extract repair, then full-refresh the blueprint; verify `MAX(_raw_data_id)` in `_tool_cursor_*` catches up to `MAX(id)` in `_raw_cursor_*` |
| Reconciliation delta on dashboard | Expected when comparing event-level charges to billing-cycle spend snapshots; see dashboard notes |

Tokens are sanitized before persisting. Connection test results include a `permissions` object showing which Admin API endpoints succeeded.

## Not collected (Enterprise AI Code Tracking)

The following endpoints from the [AI Code Tracking API](https://cursor.com/docs/account/teams/ai-code-tracking-api) are **Enterprise plan only** and are **not implemented** in this plugin (no collector, extractor, or `_tool_cursor_*` table):

| Endpoint | Purpose |
|----------|---------|
| `GET /analytics/ai-code/commits` | Per-commit AI line attribution (TAB vs Composer vs non-AI) |
| `GET /analytics/ai-code/changes` | Granular accepted AI changes |
| `GET /analytics/ai-code/commits.csv` / `changes.csv` | Bulk CSV exports of the above |

Team/Business Admin API keys typically receive **401** on these routes. Connection test optionally probes `GET /analytics/team/dau` (`probeEnterpriseAnalytics`) to detect Enterprise access; it does **not** ingest analytics data.

**Metrics this would unlock** (shown on the Cursor native dashboard but absent from DevLake today):

- AI share of **committed** code (not editor line acceptance)
- Per-commit `commitSource`: IDE, CLI, or cloud
- Repo name and primary-branch filters
- TAB vs Composer line attribution on commits

The **Cursor Usage** Grafana dashboard (`grafana/dashboards/mysql/cursor-usage.json`) uses Admin API proxies from `_tool_cursor_daily_usage` and `_tool_cursor_usage_events` instead; panel descriptions note where metrics are approximate.

## Limitations

- **Team/Business Admin API only** — Enterprise-only Analytics API endpoints (`/analytics/*`) are not collected in this plugin (see [Not collected (Enterprise AI Code Tracking)](#not-collected-enterprise-ai-code-tracking) above).
- **Tool layer only** — no domain-layer tables; cross-plugin joins (Jira, GitHub PRs, etc.) are done in Grafana SQL or separate tooling.
- **Team-level scope** — one scope per connection represents the whole team; per-team multi-tenant collection is not supported.
- **Beta** — the plugin is marked beta in Config UI while the Admin API surface continues to evolve.

## Testing

```sh
# Unit tests
cd backend && go test ./plugins/cursor/...

# E2E (requires E2E_DB_URL)
make e2e-test
```

E2E fixtures live in `backend/plugins/cursor/e2e/raw_tables/` and `e2e/snapshot_tables/`.

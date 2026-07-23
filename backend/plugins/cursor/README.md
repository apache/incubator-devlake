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
- **Incremental runs**: subsequent runs use the pipeline sync policy (`TimeAfter`) to collect only new data.
- **Date range chunking**: both `/teams/daily-usage-data` and `/teams/filtered-usage-events` requests are split into **30-day** chunks (API limit for daily usage; applied to usage events for resilience).
- **Extract**: `extractUsageEvents` and `extractDailyUsage` use **StatefulApiExtractor** (incremental by default). A config version bump triggers a one-time full re-extract after upgrade.
- **Rate limiting**: collectors honor `Retry-After` response headers and respect the configured `rateLimitPerHour`.

## Dashboards

Grafana dashboard JSON lives under `grafana/dashboards/mysql/`:

| Dashboard | File | UID |
|-----------|------|-----|
| Cursor Usage & Cost | `cursor-usage.json` | `cursor_usage` |
| AI Cost Efficiency (Cursor panels) | `ai-cost-efficiency.json` | — |
| Multi-AI Comparison (Cursor panels) | `multi-ai-comparison.json` | — |

See `grafana/dashboards/mysql/CursorREADME.md` for dashboard prerequisites, variables, and panel descriptions.

## Error handling

| Symptom | Likely cause |
|---------|--------------|
| **401 Unauthorized** on test connection | Invalid API key, or a User API key instead of a Team Admin key |
| **403 Forbidden** on spend or usage events | Key lacks permission for that Admin API endpoint |
| **429 Too Many Requests** | Rate limit exceeded — lower `rateLimitPerHour` or wait for `Retry-After` |
| Empty usage events after successful run | Selected time range has no billable events, or sync policy excludes the date range |
| Reconciliation delta on dashboard | Expected when comparing event-level charges to billing-cycle spend snapshots; see dashboard notes |

Tokens are sanitized before persisting. Connection test results include a `permissions` object showing which Admin API endpoints succeeded.

## Limitations

- **Team/Business Admin API only** — Enterprise-only Analytics API endpoints (`/analytics/*`) are not collected in this plugin.
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

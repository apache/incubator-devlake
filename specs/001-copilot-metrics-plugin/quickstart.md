# Quickstart — GitHub Copilot Adoption Plugin

Date: December 11, 2025  
Branch: `001-copilot-metrics-plugin`

## 1. Prerequisites
- DevLake server running from `docker-compose-dev.yml` with MySQL/PostgreSQL and Grafana.
- GitHub Copilot Business or Enterprise subscription enabled for the target organization.
- Personal Access Token (classic) with `manage_billing:copilot` scope (or org-level Copilot Business fine-grained token).
- DevLake CLI/Config UI access to create plugin connections.

## 2. Install Plugin (after implementation)
1. Pull latest DevLake code with the Copilot plugin branch.
2. Run `make dep` (installs Go + UI dependencies if needed).
3. Rebuild the server: `make build` or `make dev`.
4. Confirm Grafana dashboards sync via `make grafana-dashboard-sync` (or deploy JSON manually).

## 3. Configure Copilot Connection
1. Navigate to **Data Integrations → Add Connection → GitHub Copilot** (new entry).
2. Provide:
   - **Name**: e.g., `Copilot Octodemo`
   - **Endpoint**: default `https://api.github.com`
   - **Organization**: GitHub org slug (e.g., `octodemo`)
   - **Token**: PAT with Copilot billing scope
3. Click **Test Connection** and ensure the response shows plan type, active seats, and success status.
4. Save the connection.

## 4. Define Scope
1. In the newly created connection, open **Scopes**.
2. Add the organization scope (same slug). Leave `implementationDate` blank for Phase 1 (Phase 2 will use it).
3. Save scopes; blueprint JSON will reference `scopeId = <orgSlug>`.

## 5. Create Blueprint
```json
[
  [
    {
      "plugin": "copilot",
      "options": {
        "connectionId": 1,
        "scopeId": "octodemo"
      }
    }
  ]
]
```
- Schedule the blueprint to run daily to stay within the 100-day lookback.

## 6. Run Collection
- Run the blueprint immediately (`Run Now`).
- Monitor DevLake logs for rate-limit warnings or privacy-threshold messages.
- Verify that tables `_tool_copilot_org_metrics`, `_tool_copilot_language_metrics`, and `_tool_copilot_seats` contain new records.

## 7. Explore Grafana Adoption Dashboard
1. Open Grafana → `DevLake Copilot Adoption` dashboard (new entry).
2. Select connection scope (org) via dashboard variables.
3. Panels available:
   - Active vs Engaged users over time
   - Acceptance rate (acceptances / suggestions)
   - Copilot IDE vs GitHub.com chat usage
   - Top 10 languages & editor distribution
   - Seat adoption timeline (cumulative assignments)
4. Adjust time range (e.g., last 90 days) and verify panels refresh successfully.

## 8. Troubleshooting
- **403 Forbidden**: Ensure PAT includes `manage_billing:copilot`.
- **422 Metrics Disabled**: Copilot Metrics API must be enabled in GitHub organization settings.
- **Empty datasets**: Organization may not meet the ≥5 engaged user privacy threshold; plugin will log warnings.
- **Rate limit**: Respect Retry-After headers; rerun pipeline if necessary.

## 9. Next Steps (Phase 2 Preview)
- Configure `implementationDate` once Phase 2 is available to unlock Impact Dashboard comparisons.
- Evaluate team-level metrics demand before opting into Phase 2 feature branch.

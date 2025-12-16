# Copilot Plugin (Adoption Metrics)

This directory mirrors the structure of other data-source plugins (e.g., `backend/plugins/q_dev`) and will host the GitHub Copilot adoption collector implementation.

## Structure
- `api/` – REST layer for connections, scopes, and helper endpoints
- `impl/` – Plugin metadata, options, and helper utilities
- `models/` – Tool-layer models and migration scripts
- `tasks/` – Collectors, extractors, and pipeline registrations
- `e2e/` – End-to-end fixtures and assertions for Copilot data flows

Additional documentation and code will be added in subsequent phases.

## Connection Setup

1. Navigate to **Data Integrations → Add Connection → GitHub Copilot** (once the plugin is compiled).
2. Supply the following fields:
	- **Name** – Friendly label (e.g., `Copilot Octodemo`).
	- **Endpoint** – Defaults to `https://api.github.com`; leave blank to use the default.
	- **Organization** – GitHub organization slug that has Copilot metrics enabled.
	- **Token** – Personal access token with the `manage_billing:copilot` scope (classic PAT) or equivalent fine-grained scope.
3. Click **Test Connection**. The backend calls `GET /orgs/{org}/copilot/billing` and returns plan/seat data on success.
4. Save the connection and configure scopes. Each scope corresponds to an organization slug and may carry optional adoption metadata.

### Error Handling Guidance

- **403 Forbidden** → Token is missing `manage_billing:copilot` or the organization lacks Copilot access.
- **404 Not Found** → Organization slug is incorrect or Copilot metrics are not enabled; verify settings in GitHub.
- **422 Unprocessable Entity** → Copilot metrics are disabled for the organization; enable metrics in GitHub Copilot Business/Enterprise admin settings.
- **429 Too Many Requests** → GitHub rate-limited the request; respect the `Retry-After` header before retrying.

Tokens are sanitized before persisting. When patching an existing connection, omit the token to retain the encrypted value already stored in DevLake.

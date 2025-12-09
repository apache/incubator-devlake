# Feature Specification: GitHub Copilot Plugin - Adoption Metrics (Phase 1)

**Feature Branch**: `001-copilot-metrics-plugin`  
**Created**: December 5, 2025  
**Status**: Draft  
**Phase**: 1 of 2 (see `002-copilot-impact-dashboard` for Phase 2)  
**Input**: User description: "Build a GitHub Copilot plugin for Apache DevLake that collects Copilot usage metrics from GitHub's REST API (org-level and team-level) and provides an Adoption Dashboard"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Configure Copilot Connection (Priority: P1)

As a DevLake administrator, I want to configure a connection to GitHub's Copilot API so that I can start collecting Copilot usage metrics for my organization.

**Why this priority**: Without a valid connection, no metrics can be collected. This is the foundational capability that enables all other features.

**Independent Test**: Can be fully tested by creating a connection with a valid GitHub PAT, testing the connection, and verifying successful authentication response with org/billing information.

**Acceptance Scenarios**:

1. **Given** I have a GitHub PAT with `manage_billing:copilot` scope, **When** I create a new Copilot connection with org name and token, **Then** the system validates the connection and shows org details (plan type, seat count).

2. **Given** I have an existing Copilot connection, **When** I test the connection, **Then** I see a success message with current seat utilization (total seats, active seats).

3. **Given** I provide an invalid token or wrong org name, **When** I test the connection, **Then** I see a clear error message explaining the issue (permission denied, org not found, Copilot not enabled).

---

### User Story 2 - Collect Daily Copilot Usage Metrics (Priority: P1)

As an engineering manager, I want DevLake to automatically collect daily Copilot usage metrics so that I can track adoption and usage patterns across my organization.

**Why this priority**: Core data collection is essential for any dashboard or analysis. Without metrics data, no value can be delivered to users.

**Independent Test**: Can be fully tested by triggering a data collection job and verifying that metrics appear in the database tables with expected fields (active users, acceptance rates, chat usage).

**Acceptance Scenarios**:

1. **Given** a configured Copilot connection and scope, **When** I run a data collection pipeline, **Then** the system collects daily metrics including code completion suggestions/acceptances, active users, and chat usage.

2. **Given** the system has previously collected metrics, **When** I run collection again, **Then** only new days since the last collection are fetched (incremental collection respecting 100-day API lookback limit).

3. **Given** the organization has fewer than 5 active Copilot users on a team, **When** metrics collection runs, **Then** the system logs a warning about privacy threshold and continues without failing.

4. **Given** a valid scope with team filtering enabled, **When** I run collection, **Then** only metrics for the specified team are collected.

---

### User Story 3 - View Copilot Adoption Dashboard (Priority: P1)

As a DevOps leader, I want to view a dashboard showing Copilot adoption trends so that I can understand how the tool is being used across the organization.

**Why this priority**: Visualizing adoption data is the primary way stakeholders consume Copilot metrics. Without dashboards, raw database data has limited value.

**Independent Test**: Can be fully tested by loading the Grafana dashboard after metrics collection and verifying all panels display data correctly (user trends, acceptance rates, language breakdown).

**Acceptance Scenarios**:

1. **Given** Copilot metrics have been collected, **When** I open the Adoption Dashboard, **Then** I see active and engaged user counts over time as a trend line.

2. **Given** Copilot metrics exist for multiple languages, **When** I view the language breakdown panel, **Then** I see a bar chart of the top languages by suggestions/acceptances.

3. **Given** metrics include chat usage data, **When** I view the chat usage panel, **Then** I see IDE chat vs GitHub.com chat usage as separate trend lines.

4. **Given** I filter by date range, **When** the dashboard refreshes, **Then** all panels update to show data only within the selected range.

---

### User Story 4 - Collect Language and Editor Breakdown (Priority: P2)

As an engineering leader, I want to see Copilot usage broken down by programming language and IDE so that I can understand adoption patterns across different tech stacks.

**Why this priority**: Language/editor breakdown provides useful insights but is optional detail beyond core usage metrics.

**Independent Test**: Can be fully tested by viewing language breakdown panels after collection and verifying data for multiple languages/editors appears.

**Acceptance Scenarios**:

1. **Given** metrics collection has run, **When** I view language breakdown, **Then** I see suggestions and acceptances per programming language.

2. **Given** developers use multiple editors (VS Code, JetBrains, Neovim), **When** I view editor distribution, **Then** I see a pie chart showing usage share by editor.

---

### User Story 5 - Track PR Summary Usage by Repository (Priority: P3)

As a repository owner, I want to see which repositories use Copilot PR summaries so that I can encourage adoption in lower-usage repos.

**Why this priority**: PR summary tracking is a nice-to-have feature that depends on core metrics collection.

**Independent Test**: Can be fully tested by viewing the PR summaries panel and verifying repository-level counts appear.

**Acceptance Scenarios**:

1. **Given** repositories have PR summary data, **When** I view the PR summaries panel, **Then** I see a table of repositories ranked by PR summaries created.

---

### Edge Cases

- What happens when the GitHub API rate limit is exceeded? The system should respect `Retry-After` headers and resume collection without losing progress.
- How does the system handle organizations that disable the Copilot Metrics API? Display a clear error message (422 response) indicating the feature must be enabled in GitHub settings.
- What if no Copilot data exists for the specified date range? Display "No data available" in dashboard panels rather than errors.
- How are gaps in collection handled (e.g., DevLake was down for 3 days)? Backfill missing days up to the 100-day API lookback limit on next run.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow users to create, update, and delete Copilot connections with GitHub PAT and organization name.
- **FR-002**: System MUST validate connections by testing GitHub API authentication and Copilot billing access.
- **FR-003**: System MUST collect daily Copilot usage metrics including active users, engaged users, suggestions, and acceptances.
- **FR-004**: System MUST collect chat usage metrics for both IDE and GitHub.com chat interactions.
- **FR-005**: System MUST collect PR summary counts per repository from the Copilot API.
- **FR-006**: System MUST support incremental collection to respect the 100-day API lookback limit.
- **FR-007**: System MUST allow configuration of organization-level or team-level scope for metric collection.
- **FR-008**: System MUST provide an Adoption Dashboard showing usage trends (active & engaged users), acceptance rates, and user counts.
- **FR-009**: System MUST handle API errors gracefully (403, 404, 422, 429) with clear user-facing messages.
- **FR-010**: System MUST respect team privacy thresholds (5+ users required for data) without failing collection.
- **FR-011**: System MUST store language and editor breakdown metrics for granular usage analysis.
- **FR-012**: System MUST support seat assignment collection to track adoption timeline.

### Key Entities

- **CopilotConnection**: Authentication credentials linking DevLake to a GitHub organization's Copilot subscription. Contains token, organization name, and API endpoint.
- **CopilotScope**: Defines the data collection boundary (organization or team level).
- **CopilotOrgMetrics**: Daily aggregate usage data including active users, suggestions, acceptances, and chat interactions. Primary entity for adoption dashboards.
- **CopilotLanguageMetrics**: Breakdown of usage by programming language and editor. Supports granular adoption analysis.
- **CopilotSeat**: Individual seat assignment record tracking when users received Copilot access. Enables adoption timeline visualization.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can configure a Copilot connection and complete initial metrics collection within 10 minutes.
- **SC-002**: Adoption Dashboard displays current-week Copilot usage within 24 hours of data availability from GitHub.
- **SC-003**: System handles 100-day historical collection without timeout or resource exhaustion.
- **SC-004**: 90% of users can identify Copilot adoption trends by viewing the Adoption Dashboard without additional training.
- **SC-005**: System recovers gracefully from API rate limits and completes collection within the same pipeline run.
- **SC-006**: Clear error messages displayed for all API failure scenarios (invalid token, missing permissions, API disabled).

## Assumptions

- GitHub Copilot Business or Enterprise subscription is active for the target organization.
- Users have access to a GitHub PAT with `manage_billing:copilot` or equivalent fine-grained permissions.
- The organization has enabled the Copilot Metrics API feature (may require opt-in in GitHub settings).
- Standard web application performance expectations apply (pages load within 3 seconds, dashboards render within 5 seconds).
- Data retention follows DevLake's existing data management policies.

## Out of Scope (Phase 2)

The following capabilities are deferred to Phase 2 (`002-copilot-impact-dashboard`):
- Implementation date configuration for before/after analysis
- Baseline period configuration
- Impact Dashboard with before/after velocity comparisons
- Correlation with `project_pr_metrics` data
- Before/after deployment frequency analysis

# Feature Specification: GitHub Copilot Impact Dashboard (Phase 2)

**Feature Branch**: `002-copilot-impact-dashboard`  
**Created**: December 5, 2025  
**Status**: Draft  
**Phase**: 2 of 2  
**Depends On**: `001-copilot-metrics-plugin` (Copilot plugin and Adoption Dashboard)  
**Input**: User description: "Build an Impact Dashboard that compares engineering velocity metrics before vs after Copilot adoption using implementation date configuration and correlation with existing PR metrics"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Configure Copilot Implementation Date (Priority: P1)

As an engineering manager, I want to specify when Copilot was rolled out to my organization/team so that I can compare productivity metrics before and after adoption.

**Why this priority**: Implementation date is the foundational configuration required for all before/after impact analysis. Without it, no comparisons can be made.

**Independent Test**: Can be fully tested by configuring an implementation date in scope settings and verifying the date is persisted and displayed correctly.

**Acceptance Scenarios**:

1. **Given** I am editing a Copilot scope configuration, **When** I set an implementation date, **Then** the date is saved and displayed in scope details.

2. **Given** I set a baseline period (e.g., 90 days before implementation to compare), **When** I save the configuration, **Then** the baseline period is stored and used in impact calculations.

3. **Given** I have not set an implementation date, **When** I try to access the Impact Dashboard, **Then** I see a message prompting me to configure the implementation date first.

---

### User Story 2 - View Before/After PR Cycle Time Comparison (Priority: P1)

As a VP of Engineering, I want to compare PR cycle time before and after Copilot adoption so that I can measure the tool's impact on development velocity.

**Why this priority**: PR cycle time is the most requested metric for measuring engineering productivity impact. This is the core value proposition of the Impact Dashboard.

**Independent Test**: Can be fully tested by viewing the Impact Dashboard with a configured implementation date and verifying before/after PR cycle time panels show meaningful comparison data.

**Acceptance Scenarios**:

1. **Given** an implementation date is configured and PR metrics exist in DevLake, **When** I open the Impact Dashboard, **Then** I see average PR cycle time for the baseline period vs post-implementation period.

2. **Given** sufficient historical data exists, **When** I view the impact summary, **Then** I see percentage change metrics (e.g., "PR cycle time reduced by 23%").

3. **Given** I view the time series panel, **When** implementation date is set, **Then** I see a vertical annotation line marking the Copilot rollout date.

4. **Given** I hover over the before/after comparison, **When** I view details, **Then** I see breakdown of cycle time components (coding time, pickup time, review time).

---

### User Story 3 - View Deployment Frequency Impact (Priority: P2)

As a DevOps leader, I want to compare deployment frequency before and after Copilot adoption so that I can assess impact on release velocity.

**Why this priority**: Deployment frequency is a key DORA metric that provides insight into overall team velocity beyond just PR speed.

**Independent Test**: Can be fully tested by viewing the deployment frequency panel and verifying before/after deployment counts are displayed.

**Acceptance Scenarios**:

1. **Given** deployment data exists in DevLake, **When** I view the deployment frequency panel, **Then** I see deployment counts per week for baseline period vs post-implementation period.

2. **Given** I view the deployment comparison, **When** data spans the implementation date, **Then** I see side-by-side bar charts showing before vs after deployment rates.

---

### User Story 4 - Correlate Copilot Adoption with PR Speed (Priority: P2)

As an engineering leader, I want to see how Copilot adoption levels correlate with PR processing speed so that I can identify whether higher usage leads to faster delivery.

**Why this priority**: Correlation analysis provides deeper insight into the relationship between Copilot usage and outcomes, beyond simple before/after comparison.

**Independent Test**: Can be fully tested by viewing the correlation panel and verifying scatter plot or trend data appears linking Copilot metrics to PR metrics.

**Acceptance Scenarios**:

1. **Given** both Copilot metrics and PR metrics exist for overlapping dates, **When** I view the correlation panel, **Then** I see a visualization showing the relationship between daily Copilot active users and average PR cycle time.

2. **Given** I filter by date range, **When** the dashboard refreshes, **Then** the correlation calculation updates to reflect only the selected period.

---

### User Story 5 - View Code Review Time Impact (Priority: P3)

As a tech lead, I want to see if Copilot adoption has reduced code review time so that I can understand its effect on the review process.

**Why this priority**: Code review time is a specific component of PR cycle time that may show direct impact from Copilot-assisted coding.

**Independent Test**: Can be fully tested by viewing the review time panel and verifying before/after comparison data appears.

**Acceptance Scenarios**:

1. **Given** PR review time data exists, **When** I view the review time impact panel, **Then** I see average review time before vs after Copilot implementation.

2. **Given** I view review time distribution, **When** comparing periods, **Then** I see box plots or histograms showing the distribution shift.

---

### Edge Cases

- What if implementation date is set in the future? Display a message indicating the before/after comparison will be available after the implementation date passes.
- What if there is no PR data before the implementation date? Display "Insufficient baseline data" with guidance on minimum data requirements.
- What if implementation date is more than 100 days ago (beyond Copilot API lookback)? The Impact Dashboard can still function using existing PR metrics; only Copilot adoption trend data may be incomplete.
- How are outliers handled in cycle time calculations? Use median in addition to mean, and provide filtering options for extreme values.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow users to configure an implementation date in the Copilot scope configuration.
- **FR-002**: System MUST allow users to configure a baseline period (default: 90 days) for before/after comparison.
- **FR-003**: System MUST provide an Impact Dashboard comparing engineering velocity metrics before vs after implementation date.
- **FR-004**: System MUST display PR cycle time comparison (average, median) for baseline vs post-implementation periods.
- **FR-005**: System MUST display PR cycle time component breakdown (coding time, pickup time, review time) for both periods.
- **FR-006**: System MUST display deployment frequency comparison for baseline vs post-implementation periods.
- **FR-007**: System MUST calculate and display percentage change for key metrics.
- **FR-008**: System MUST display a vertical annotation line on time series charts marking the implementation date.
- **FR-009**: System MUST correlate Copilot metrics with PR metrics using date-based joins on existing `project_pr_metrics` table.
- **FR-010**: System MUST provide a correlation visualization showing relationship between Copilot adoption and PR speed.
- **FR-011**: System MUST display informative messages when insufficient data exists for comparison.
- **FR-012**: System MUST support date range filtering on Impact Dashboard panels.

### Key Entities

- **CopilotScope** (extended): Adds implementation date and baseline period configuration to the scope model from Phase 1.
- **project_pr_metrics** (existing): DevLake's existing PR metrics table containing cycle time, coding time, pickup time, review time, and deploy time. Used via SQL joins for impact analysis.
- **cicd_deployment_commits** (existing): DevLake's existing deployment table used for deployment frequency analysis.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can configure implementation date and view Impact Dashboard within 5 minutes of completing Phase 1 setup.
- **SC-002**: Impact Dashboard accurately shows before/after comparison when implementation date is set and 90+ days of PR data exists on both sides.
- **SC-003**: Before/after impact metrics match manual calculations within 2% margin when validated against raw data.
- **SC-004**: 90% of users can identify whether Copilot improved their PR cycle time by viewing the Impact Dashboard summary.
- **SC-005**: Dashboard panels load within 5 seconds for organizations with up to 10,000 PRs.
- **SC-006**: Clear guidance displayed when data requirements are not met (insufficient baseline, no implementation date).

## Assumptions

- Phase 1 (`001-copilot-metrics-plugin`) is complete and Copilot metrics are being collected.
- DevLake already has PR data from GitHub/GitLab plugins for the repositories being analyzed.
- DevLake has deployment data (via CI/CD plugins) for deployment frequency analysis; if not, that panel will show "No data available."
- The `project_pr_metrics` domain table contains cycle time components (pr_cycle_time, pr_coding_time, pr_pickup_time, pr_review_time, pr_deploy_time).
- Users understand that correlation does not imply causation; dashboard should include appropriate disclaimers.

## Dependencies

- **Hard Dependency**: `001-copilot-metrics-plugin` must be merged and deployed before this feature can function.
- **Soft Dependency**: Existing PR metrics from GitHub/GitLab plugin for the target repositories.
- **Soft Dependency**: Existing CICD deployment data for deployment frequency panels.

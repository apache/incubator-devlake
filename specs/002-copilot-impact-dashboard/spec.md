# Feature Specification: GitHub Copilot Impact Dashboard (Phase 2)

**Feature Branch**: `002-copilot-impact-dashboard`  
**Created**: December 5, 2025  
**Status**: Draft  
**Phase**: 2 of 2  
**Depends On**: `001-copilot-metrics-plugin` (Copilot plugin and Adoption Dashboard)  
**Input**: User description: "Build an Impact Dashboard that continuously correlates GitHub Copilot adoption levels with DORA/engineering velocity metrics, with optional implementation date annotations for milestone comparisons"

## Design Philosophy

**Correlation-First Approach**: Rather than requiring a single "implementation date" for before/after analysis, this dashboard primarily shows **continuous correlation** between Copilot adoption metrics (active users, acceptance rate) and DORA metrics (PR cycle time, deployment frequency, CFR, MTTR). This reflects real-world phased rollouts where:

1. Pilot teams adopt first, then gradual expansion
2. Adoption levels vary week-over-week
3. Impact correlates with usage intensity, not just time

**Optional Milestones**: Users may optionally configure an implementation date to add annotation markers to time series charts, enabling point-in-time "before/after" comparisons as a secondary analysis mode.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Correlate Copilot Adoption with DORA Metrics (Priority: P1) 🎯 PRIMARY

As an engineering leader, I want to see how Copilot adoption levels (active users, acceptance rate) correlate with DORA metrics over time so that I can understand the continuous impact of Copilot regardless of a specific rollout date.

**Why this priority**: Phased rollouts mean there's no single "implementation date." Continuous correlation shows impact as adoption grows organically over time and captures the relationship between usage intensity and productivity outcomes.

**Independent Test**: Can be fully tested by viewing the Impact Dashboard with Copilot metrics data and verifying correlation visualizations appear without any implementation date configuration.

**Acceptance Scenarios**:

1. **Given** Copilot metrics are being collected and PR metrics exist, **When** I open the Impact Dashboard without configuring an implementation date, **Then** I see a dual-axis time series showing Copilot active users alongside PR cycle time trend.

2. **Given** adoption varies week-over-week, **When** I view the correlation panel, **Then** I see a scatter plot correlating "Copilot Active Users %" with "Avg PR Cycle Time" for each time period.

3. **Given** I want to quantify correlation strength, **When** I view the correlation details, **Then** I see a Pearson correlation coefficient (e.g., r = -0.72) with interpretation guidance (strong negative = higher adoption correlates with faster PRs).

4. **Given** I filter by date range, **When** the dashboard refreshes, **Then** the correlation coefficient and visualizations update to reflect only the selected period.

5. **Given** I want to see correlation with other DORA metrics, **When** I view the full dashboard, **Then** I see correlation analysis for Deployment Frequency, CFR, and MTTR alongside PR cycle time.

---

### User Story 2 - Configure Rollout Milestone (Priority: P2, Optional)

As an engineering manager, I want to optionally specify rollout milestone dates so that I can add visual markers to trend charts and enable supplementary before/after comparisons.

**Why this priority**: While correlation analysis is primary, some organizations want to highlight specific rollout milestones (pilot launch, full rollout, etc.) on charts for context.

**Independent Test**: Can be fully tested by configuring an implementation date in scope settings and verifying the annotation appears on time series charts.

**Acceptance Scenarios**:

1. **Given** I am editing a Copilot scope configuration, **When** I set an implementation date, **Then** the date is saved and displayed in scope details.

2. **Given** I set an implementation date, **When** I view time series charts, **Then** I see a vertical annotation line marking the rollout milestone.

3. **Given** I set a baseline period (e.g., 90 days before milestone), **When** I save the configuration, **Then** supplementary before/after comparison panels become available.

4. **Given** I have NOT set an implementation date, **When** I access the Impact Dashboard, **Then** the dashboard still functions fully using correlation analysis (no blocking message).

---

### User Story 3 - View PR Cycle Time Correlation (Priority: P1)

As a VP of Engineering, I want to see how PR cycle time correlates with Copilot adoption levels so that I can measure the tool's impact on development velocity.

**Why this priority**: PR cycle time is the most requested metric for measuring engineering productivity impact. This is the core value proposition of the Impact Dashboard.

**Independent Test**: Can be fully tested by viewing the Impact Dashboard with Copilot and PR metrics overlapping, verifying correlation panels show meaningful data.

**Acceptance Scenarios**:

1. **Given** Copilot metrics and PR metrics exist for overlapping dates, **When** I open the Impact Dashboard, **Then** I see a dual-axis time series showing Copilot active users % and average PR cycle time trending together.

2. **Given** I view the correlation scatter plot, **When** data is available, **Then** I see weekly data points plotting adoption % vs cycle time with a trend line.

3. **Given** an implementation date is optionally configured, **When** I view the time series panel, **Then** I see a vertical annotation line marking the rollout milestone.

4. **Given** an implementation date is configured, **When** I view supplementary panels, **Then** I see before/after average comparison with percentage change.

5. **Given** I hover over the correlation details, **When** I view breakdown, **Then** I see cycle time components (coding time, pickup time, review time) correlated separately.

---

### User Story 4 - View Deployment Frequency Correlation (Priority: P2)

As a DevOps leader, I want to see how deployment frequency correlates with Copilot adoption so that I can assess impact on release velocity.

**Why this priority**: Deployment frequency is a key DORA metric that provides insight into overall team velocity beyond just PR speed.

**Independent Test**: Can be fully tested by viewing the deployment frequency panels and verifying correlation with Copilot adoption is displayed.

**Acceptance Scenarios**:

1. **Given** deployment data and Copilot metrics exist, **When** I view the deployment frequency panel, **Then** I see a dual-axis chart showing weekly deployments alongside Copilot adoption %.

2. **Given** I view the correlation analysis, **When** data spans multiple weeks, **Then** I see correlation coefficient between adoption and deployment frequency.

3. **Given** an implementation date is optionally configured, **When** I view supplementary panels, **Then** I see before/after deployment count comparison.

---

### User Story 5 - View Adoption Intensity Analysis (Priority: P2)

As an engineering leader, I want to see detailed breakdowns of how different Copilot adoption intensities correlate with outcomes so that I can identify optimal usage patterns.

**Why this priority**: Beyond simple correlation, understanding adoption intensity (low/medium/high weeks) helps identify thresholds where Copilot impact becomes significant.

**Independent Test**: Can be fully tested by viewing the intensity analysis panel showing bucketed adoption levels.

**Acceptance Scenarios**:

1. **Given** Copilot metrics span multiple weeks with varying adoption, **When** I view the intensity panel, **Then** I see metrics grouped by adoption tier (e.g., <25%, 25-50%, 50-75%, >75% active users).

2. **Given** I view intensity buckets, **When** comparing tiers, **Then** I see average PR cycle time, deployment frequency for each adoption tier.

3. **Given** I want trend analysis, **When** viewing the panel, **Then** I see how metrics improve (or not) as adoption tier increases.

---

### User Story 6 - View Code Review Time Correlation (Priority: P3)

As a tech lead, I want to see if higher Copilot adoption correlates with reduced code review time so that I can understand its effect on the review process.

**Why this priority**: Code review time is a specific component of PR cycle time that may show direct impact from Copilot-assisted coding.

**Independent Test**: Can be fully tested by viewing the review time panel and verifying correlation with Copilot adoption appears.

**Acceptance Scenarios**:

1. **Given** PR review time data and Copilot metrics exist, **When** I view the review time panel, **Then** I see correlation between Copilot adoption % and average review time.

2. **Given** I view review time distribution, **When** comparing high vs low adoption periods, **Then** I see distribution shift visualization.

---

### User Story 7 - View Change Failure Rate Correlation (Priority: P2)

As a DevOps leader, I want to see how Change Failure Rate correlates with Copilot adoption so that I can validate whether higher Copilot usage improves code quality in production.

**Why this priority**: Change Failure Rate is a key DORA metric that directly measures production code quality. This helps assess whether Copilot-assisted code is more or less reliable.

**Independent Test**: Can be fully tested by viewing the CFR panel and verifying correlation with Copilot adoption is displayed.

**Acceptance Scenarios**:

1. **Given** deployment, incident, and Copilot data exists, **When** I view the CFR panel, **Then** I see dual-axis chart showing CFR % alongside Copilot adoption %.

2. **Given** I view correlation analysis, **When** sufficient data spans exist, **Then** I see correlation coefficient between adoption and CFR (expecting negative correlation = higher adoption, lower failures).

3. **Given** an implementation date is optionally configured, **When** I view the CFR trend, **Then** I see a vertical annotation line marking the rollout milestone.

4. **Given** no incident data is available, **When** I view the CFR panel, **Then** I see "N/A. Please check if you have collected incidents" with configuration guidance.

---

### User Story 8 - View MTTR/Recovery Time Correlation (Priority: P3)

As an incident manager, I want to see if higher Copilot adoption correlates with faster incident recovery time so that I can understand if Copilot helps teams resolve production issues faster.

**Why this priority**: Mean Time to Recovery (MTTR) shows how quickly teams can fix production issues. This may improve with Copilot assistance during incident response.

**Independent Test**: Can be fully tested by viewing the MTTR panel and verifying correlation with Copilot adoption.

**Acceptance Scenarios**:

1. **Given** incident data and Copilot metrics exist, **When** I view the MTTR panel, **Then** I see correlation between Copilot adoption % and median recovery time.

2. **Given** I view recovery time trend, **When** comparing high vs low adoption periods, **Then** I see distribution comparison.

3. **Given** I hover over the MTTR details, **When** viewing breakdown, **Then** I see incident counts and recovery time statistics for different adoption levels.

---

### User Story 9 - View Code Quality Metrics Correlation (Priority: P3, Optional)

As a tech lead, I want to see if code quality metrics (bugs, vulnerabilities, code smells, complexity) correlate with Copilot adoption so that I can assess Copilot's impact on code health.

**Why this priority**: Code quality metrics provide objective measures of code health. This is optional because it requires SonarQube or similar code quality tool integration.

**Independent Test**: Can be fully tested by viewing the Code Quality panel and verifying correlation with Copilot adoption when SonarQube data is available.

**Acceptance Scenarios**:

1. **Given** SonarQube and Copilot metrics exist, **When** I view the Code Quality panel, **Then** I see correlation between Copilot adoption and bugs/vulnerabilities/code smells trends.

2. **Given** code complexity data exists, **When** I view the complexity correlation, **Then** I see relationship between adoption and complexity metrics over time.

3. **Given** code coverage data exists, **When** I view the quality panel, **Then** I see coverage trend alongside Copilot adoption trend.

4. **Given** no SonarQube data is available, **When** I view the Code Quality panel, **Then** I see "N/A. Please configure SonarQube integration" with documentation links.

5. **Given** an implementation date is optionally configured, **When** I view quality trends, **Then** I see annotation marking the rollout milestone.

---

### Edge Cases

**Correlation Analysis:**
- What if no Copilot metrics exist yet? Display "Awaiting Copilot data. Please run a collection first." with guidance.
- What if Copilot and DORA metrics don't have overlapping dates? Display correlation panels as "Insufficient overlapping data" and show individual trends separately.
- How many data points are needed for meaningful correlation? Minimum 2 weeks of overlapping data recommended; display warning if less.
- What if correlation is weak (|r| < 0.3)? Display "Weak or no correlation detected" with explanation that other factors may be influencing metrics.
- What if correlation is counter-intuitive (positive when negative expected)? Display result honestly with disclaimer about confounding variables.

**Optional Rollout Milestone:**
- What if implementation date is set in the future? Display annotation as "Planned Rollout" and continue showing correlation analysis.
- What if there is no PR data before the implementation date? Before/after supplementary panels show "Insufficient baseline data" but correlation panels still work.
- What if implementation date is more than 100 days ago (beyond Copilot API lookback)? Correlation analysis may have gaps; annotate affected periods.

**Adoption Intensity:**
- What if adoption is consistently low (<10%)? Show intensity tiers but note "Limited high-adoption data available."
- How are adoption tiers calculated? Default: <25%, 25-50%, 50-75%, >75% of total seats active.

**Data Quality:**
- How are outliers handled in cycle time calculations? Use median in addition to mean, and provide filtering options for extreme values.
- What if incident data is not configured? CFR and MTTR panels display "N/A. Please check if you have collected incidents" with documentation links.
- What if deployments exist but no incidents are linked? CFR shows 0% (perfect success rate) with a note about validation.
- How are incidents linked to deployments for CFR? Use DevLake's existing `project_incident_deployment_relationships` table.
- What if SonarQube is not integrated? Code Quality panel displays "N/A. Please configure SonarQube integration" with setup documentation links.
- How are code quality metrics aggregated? Use average per-file metrics from `cq_file_metrics` table, weighted by lines of code when available.
- What if code quality data exists but is sparse? Display panels that have data and show "Insufficient data" for missing metrics.

## Requirements *(mandatory)*

### Functional Requirements

**Core Correlation Analysis (Primary):**
- **FR-001**: System MUST provide an Impact Dashboard that functions without requiring an implementation date configuration.
- **FR-002**: System MUST display continuous correlation analysis between Copilot adoption metrics and DORA metrics.
- **FR-003**: System MUST calculate and display Pearson correlation coefficient between Copilot active users % and each DORA metric.
- **FR-004**: System MUST show dual-axis time series charts overlaying Copilot adoption trend with each DORA metric trend.
- **FR-005**: System MUST provide scatter plot visualizations correlating adoption levels with metric values.
- **FR-006**: System MUST display correlation coefficient interpretation guidance (strong/moderate/weak, positive/negative).
- **FR-007**: System MUST support date range filtering on all Impact Dashboard panels.

**Adoption Intensity Analysis:**
- **FR-008**: System MUST display adoption intensity analysis with metrics grouped by adoption tier (<25%, 25-50%, 50-75%, >75%).
- **FR-009**: System MUST calculate and display average DORA metrics for each adoption tier.

**Optional Rollout Milestone:**
- **FR-010**: System SHOULD allow users to optionally configure an implementation/rollout date in the Copilot scope configuration.
- **FR-011**: System SHOULD allow users to configure a baseline period (default: 90 days) for supplementary before/after comparison.
- **FR-012**: System SHOULD display a vertical annotation line on time series charts marking the rollout milestone when configured.
- **FR-013**: System SHOULD provide supplementary before/after comparison panels when implementation date is configured.

**PR Cycle Time Correlation:**
- **FR-014**: System MUST display PR cycle time correlation with Copilot adoption (average, median).
- **FR-015**: System MUST display PR cycle time component breakdown (coding time, pickup time, review time) correlated with adoption.
- **FR-016**: System MUST calculate and display percentage change for key metrics when comparing high vs low adoption periods.

**Deployment Frequency Correlation:**
- **FR-017**: System MUST display deployment frequency correlation with Copilot adoption.

**Change Failure Rate Correlation:**
- **FR-018**: System MUST display Change Failure Rate correlation with Copilot adoption.
- **FR-019**: System MUST display failed deployment counts and total deployment counts alongside adoption.
- **FR-020**: System MUST link incidents to deployments using existing DevLake incident-deployment relationship data.

**MTTR Correlation:**
- **FR-021**: System MUST display MTTR/Recovery Time correlation with Copilot adoption (median).

**Code Quality Correlation (Optional):**
- **FR-022**: System SHOULD display code quality metrics correlation (bugs, vulnerabilities, code smells) with Copilot adoption when SonarQube data is available.
- **FR-023**: System SHOULD display code complexity correlation (cyclomatic, cognitive complexity) when available.
- **FR-024**: System SHOULD display code coverage correlation when available.
- **FR-025**: System SHOULD aggregate code quality metrics by averaging per-file metrics, weighted by lines of code (NCLOC) when available. If NCLOC is null or zero, use unweighted average.

**Graceful Degradation:**
- **FR-026**: System MUST display informative messages when insufficient overlapping data exists for correlation.
- **FR-027**: System MUST gracefully handle missing incident data by displaying "N/A" messages with configuration guidance.
- **FR-028**: System MUST gracefully handle missing SonarQube data by displaying "N/A. Please configure SonarQube integration" with documentation links.

### Key Entities

- **CopilotScope** (extended): Adds implementation date and baseline period configuration to the scope model from Phase 1.
- **project_pr_metrics** (existing): DevLake's existing PR metrics table containing cycle time, coding time, pickup time, review time, and deploy time. Used via SQL joins for impact analysis.
- **cicd_deployment_commits** (existing): DevLake's existing deployment table used for deployment frequency analysis.
- **incidents** (existing): DevLake's domain table for production incidents, containing `lead_time_minutes` for MTTR calculation and `resolution_date` for tracking when incidents were resolved.
- **project_incident_deployment_relationships** (existing): DevLake's crossdomain table linking incidents to the deployments that caused them, used for Change Failure Rate calculation.
- **cq_file_metrics** (existing, optional): DevLake's code quality domain table from SonarQube plugin, containing per-file metrics including `Bugs`, `Vulnerabilities`, `CodeSmells`, `Complexity`, `CognitiveComplexity`, `Coverage`, and `Ncloc` (lines of code). Used for code quality impact analysis when SonarQube is integrated.
- **cq_issues** (existing, optional): DevLake's code quality issues table containing individual code quality violations with severity, type, and technical debt metrics.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can view correlation-based Impact Dashboard within 2 minutes of completing Phase 1 setup (no implementation date configuration required).
- **SC-002**: Impact Dashboard accurately shows correlation analysis when Copilot metrics and DORA metrics have 2+ weeks of overlapping data.
- **SC-003**: Correlation coefficients match manual calculations within 2% margin when validated against raw data.
- **SC-004**: 90% of users can determine correlation strength (strong/moderate/weak) and direction (positive/negative) by viewing the correlation summary.
- **SC-005**: Dashboard panels load within 5 seconds for organizations with up to 10,000 PRs and 1,000 incidents.
- **SC-006**: Clear guidance displayed when data requirements are not met (insufficient overlapping data, no incident data, no SonarQube integration).
- **SC-007**: CFR calculations match DORA dashboard calculations within 2% margin when incident and deployment data is configured.
- **SC-008**: Code quality panels appear only when SonarQube data is available; otherwise, graceful "N/A" messages are shown without breaking the dashboard.
- **SC-009**: Optional rollout milestone annotation correctly appears on all time series charts when configured.
- **SC-010**: Adoption intensity analysis accurately groups weeks into correct adoption tiers.

## Assumptions

- Phase 1 (`001-copilot-metrics-plugin`) is complete and Copilot metrics are being collected.
- Copilot metrics include daily active users count and total seats, enabling adoption percentage calculation.
- DevLake already has PR data from GitHub/GitLab plugins for the repositories being analyzed.
- Copilot metrics and DORA metrics have at least 2 weeks of overlapping date ranges for meaningful correlation analysis.
- DevLake has deployment data (via CI/CD plugins) for deployment frequency analysis; if not, those panels will show "No data available."
- The `project_pr_metrics` domain table contains cycle time components (pr_cycle_time, pr_coding_time, pr_pickup_time, pr_review_time, pr_deploy_time).
- For CFR and MTTR panels: Incident data may or may not be configured. If not configured, panels will display "N/A" messages.
- Incident-to-deployment relationships are configured via data transformations in DevLake blueprints.
- For Code Quality panels: SonarQube integration is optional. If not configured, panels will display "N/A" messages.
- Users understand that correlation does not imply causation; dashboard includes appropriate disclaimers.
- Implementation date configuration is optional and serves only as a visual annotation/milestone, not a requirement for dashboard functionality.

## Dependencies

- **Hard Dependency**: `001-copilot-metrics-plugin` must be merged and deployed before this feature can function.
- **Soft Dependency**: Existing PR metrics from GitHub/GitLab plugin for the target repositories.
- **Soft Dependency**: Existing CICD deployment data for deployment frequency panels.
- **Soft Dependency**: Existing incident data (from Jira, GitHub Issues, PagerDuty, etc.) for Change Failure Rate and MTTR panels. If not configured, these panels will show "N/A" messages.
- **Soft Dependency**: SonarQube or SonarQube Cloud integration for Code Quality panels. If not configured, these panels will show "N/A" messages with setup guidance.

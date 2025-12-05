# GitHub Copilot Plugin - Specification Document

## Version: 1.0 (Draft)
## Date: December 5, 2025

---

## 1. Overview

### 1.1 Purpose
Create a DevLake plugin that collects GitHub Copilot usage metrics and enables correlation with existing engineering productivity metrics (DORA, PR cycle time, etc.) to measure AI coding assistant impact.

### 1.2 Strategic Approach
**Option B: Repository/Project-Level Analysis**

Instead of per-user productivity tracking, we measure:
- Aggregate Copilot adoption at org/team level
- Aggregate productivity metrics at repo level
- Before vs. after Copilot implementation date comparison

### 1.3 Key Stakeholder Needs
1. **Adoption Dashboard**: Track Copilot rollout and usage trends
2. **Impact Dashboard**: Measure effect on engineering velocity (PR cycle time, deployment frequency, etc.)

---

## 2. Data Sources

### 2.1 GitHub Copilot APIs

| Endpoint | Purpose | Scope |
|----------|---------|-------|
| `GET /orgs/{org}/copilot/billing` | Seat breakdown, settings | Org |
| `GET /orgs/{org}/copilot/metrics` | Daily usage metrics | Org |
| `GET /orgs/{org}/team/{team}/copilot/metrics` | Team-level metrics | Team |
| `GET /orgs/{org}/copilot/billing/seats` | Seat assignments | Org |

### 2.2 Required Permissions
- Personal Access Token (Classic): `manage_billing:copilot` OR `read:org`
- Fine-grained PAT: "GitHub Copilot Business" org permissions (read)

### 2.3 API Limitations
- **100-day lookback limit** - Must collect incrementally
- **5+ users privacy threshold** - Small teams won't return data
- **Daily processing** - Data available next day

---

## 3. Data Models

### 3.1 Connection Model

```go
// _tool_copilot_connections
type CopilotConnection struct {
    helper.BaseConnection
    Token        string `json:"token"`        // GitHub PAT
    Organization string `json:"organization"` // Target org name
    Endpoint     string `json:"endpoint"`     // API endpoint (default: api.github.com)
    RateLimitPerHour int `json:"rateLimitPerHour" default:"5000"`
}
```

### 3.2 Scope Model

```go
// _tool_copilot_scopes (org or team)
type CopilotScope struct {
    common.Scope
    Id               string    `gorm:"primaryKey"`
    Organization     string    `json:"organization"`
    Team             string    `json:"team,omitempty"`  // Optional team slug
    ImplementationDate time.Time `json:"implementationDate"` // Copilot rollout date
    BaselinePeriodDays int      `json:"baselinePeriodDays" default:"90"`
}
```

### 3.3 Org/Team Daily Metrics

```go
// _tool_copilot_org_metrics
type CopilotOrgMetrics struct {
    ConnectionId     uint64
    ScopeId          string    `gorm:"primaryKey"`
    Date             time.Time `gorm:"primaryKey"`
    
    // Engagement
    TotalActiveUsers   int
    TotalEngagedUsers  int
    
    // Code Completions (aggregated)
    CompletionSuggestions    int
    CompletionAcceptances    int
    CompletionLinessuggested int
    CompletionLinesAccepted  int
    
    // IDE Chat
    ChatTotalConversations  int
    ChatInsertionEvents     int
    ChatCopyEvents          int
    ChatEngagedUsers        int
    
    // GitHub.com Chat  
    DotcomChatTotal        int
    DotcomChatEngagedUsers int
    
    // PR Summaries
    PRSummariesCreated     int
    PREngagedUsers         int
}
```

### 3.4 Language Breakdown (Optional Detail)

```go
// _tool_copilot_language_metrics
type CopilotLanguageMetrics struct {
    ConnectionId     uint64
    ScopeId          string    `gorm:"primaryKey"`
    Date             time.Time `gorm:"primaryKey"`
    Editor           string    `gorm:"primaryKey;type:varchar(50)"`
    Language         string    `gorm:"primaryKey;type:varchar(50)"`
    
    EngagedUsers         int
    CodeSuggestions      int
    CodeAcceptances      int
    CodeLinesSuggested   int
    CodeLinesAccepted    int
}
```

### 3.5 Repository PR Metrics

```go
// _tool_copilot_repo_pr_metrics
type CopilotRepoPRMetrics struct {
    ConnectionId     uint64
    ScopeId          string    `gorm:"primaryKey"`
    Date             time.Time `gorm:"primaryKey"`
    RepoFullName     string    `gorm:"primaryKey;type:varchar(255)"`
    
    PRSummariesCreated int
    EngagedUsers       int
}
```

### 3.6 Seat Snapshot (For Adoption Tracking)

```go
// _tool_copilot_seats
type CopilotSeat struct {
    ConnectionId       uint64
    Organization       string    `gorm:"primaryKey"`
    UserLogin          string    `gorm:"primaryKey"`
    UserId             int64
    CreatedAt          time.Time // When seat was assigned
    LastActivityAt     *time.Time
    LastActivityEditor string
    PlanType           string
}
```

---

## 4. Subtasks

### 4.1 Task Flow

```
collectCopilotBilling → collectCopilotMetrics → collectCopilotSeats → extractCopilotData
```

| Task | Description | Dependencies |
|------|-------------|--------------|
| `collectCopilotBilling` | Get org billing/seat summary | None |
| `collectCopilotMetrics` | Collect daily metrics from API | Billing |
| `collectCopilotSeats` | Get seat assignments (optional) | None |
| `extractCopilotData` | Parse and store metrics | Metrics collected |

### 4.2 Incremental Collection

Since API only provides 100-day lookback:
- Store last collected date in state
- On each run, collect from last_date to today
- Handle gaps gracefully

---

## 5. API Endpoints

### 5.1 Connection CRUD

```
POST   /plugins/copilot/connections
GET    /plugins/copilot/connections
GET    /plugins/copilot/connections/:connectionId
PATCH  /plugins/copilot/connections/:connectionId
DELETE /plugins/copilot/connections/:connectionId
POST   /plugins/copilot/connections/:connectionId/test
```

### 5.2 Scope CRUD

```
PUT    /plugins/copilot/connections/:connectionId/scopes
GET    /plugins/copilot/connections/:connectionId/scopes
GET    /plugins/copilot/connections/:connectionId/scopes/:scopeId
PATCH  /plugins/copilot/connections/:connectionId/scopes/:scopeId
DELETE /plugins/copilot/connections/:connectionId/scopes/:scopeId
```

### 5.3 Test Connection Response

```json
{
  "success": true,
  "message": "Connection successful",
  "organization": "octodemo",
  "plan_type": "enterprise",
  "total_seats": 872,
  "active_seats": 627
}
```

---

## 6. Grafana Dashboards

### 6.1 Dashboard 1: Copilot Adoption & Usage

**Panels:**
1. **Active Users Over Time** (time series)
2. **Engaged Users Over Time** (time series)
3. **Acceptance Rate Trend** (line chart: acceptances/suggestions)
4. **Top Languages by Usage** (bar chart)
5. **Editor Distribution** (pie chart)
6. **Chat Usage Trend** (stacked area: IDE vs dotcom)
7. **PR Summaries by Repo** (table)

### 6.2 Dashboard 2: Copilot Impact on Velocity

**Panels:**
1. **Before/After Summary** (stat cards)
   - Avg PR Cycle Time: Before → After (% change)
   - Deployment Frequency: Before → After
   - Mean Time to Merge: Before → After

2. **PR Cycle Time Trend** (time series with annotation)
   - Vertical line at implementation_date
   - Highlight before/after regions

3. **Deployment Frequency Comparison** (bar chart)
   - Side-by-side: 90 days before vs 90 days after

4. **Code Review Time Impact** (box plot)
   - Distribution before vs after

5. **Correlation: Copilot Adoption vs PR Speed** (scatter)
   - X: Daily active Copilot users
   - Y: Same-day avg PR cycle time

### 6.3 Sample SQL for Impact Analysis

```sql
-- PR Cycle Time Before vs After Copilot Implementation
-- Uses project_pr_metrics domain table which contains:
--   pr_cycle_time: Total time from first commit to merge (minutes)
--   pr_coding_time: Time from first commit to PR created (minutes)  
--   pr_pickup_time: Time from PR created to first review (minutes)
--   pr_review_time: Time from first review to merge (minutes)
--   pr_deploy_time: Time from merge to deployment (minutes)

WITH copilot_config AS (
    SELECT implementation_date, baseline_period_days
    FROM _tool_copilot_scopes
    WHERE organization = '${org}'
),
pr_data AS (
    SELECT 
        CASE 
            WHEN pr.created_date < cc.implementation_date THEN 'Before Copilot'
            ELSE 'After Copilot'
        END AS period,
        ppm.pr_cycle_time,
        ppm.pr_coding_time,
        ppm.pr_pickup_time,
        ppm.pr_review_time,
        ppm.pr_deploy_time
    FROM project_pr_metrics ppm
    JOIN pull_requests pr ON ppm.id = pr.id
    CROSS JOIN copilot_config cc
    WHERE pr.created_date >= DATE_SUB(cc.implementation_date, INTERVAL cc.baseline_period_days DAY)
      AND pr.merged_date IS NOT NULL  -- Only merged PRs
)
SELECT 
    period,
    COUNT(*) as pr_count,
    AVG(pr_cycle_time)/60 as avg_cycle_time_hours,
    AVG(pr_coding_time)/60 as avg_coding_time_hours,
    AVG(pr_pickup_time)/60 as avg_pickup_time_hours,
    AVG(pr_review_time)/60 as avg_review_time_hours
FROM pr_data
GROUP BY period;
```

### 6.4 Deployment Frequency Impact SQL

```sql
-- Deployment Frequency Before vs After Copilot
WITH copilot_config AS (
    SELECT implementation_date, baseline_period_days
    FROM _tool_copilot_scopes
    WHERE organization = '${org}'
),
deployments AS (
    SELECT 
        cdc.finished_date,
        CASE 
            WHEN cdc.finished_date < cc.implementation_date THEN 'Before Copilot'
            ELSE 'After Copilot'
        END AS period
    FROM cicd_deployment_commits cdc
    CROSS JOIN copilot_config cc
    WHERE cdc.result = 'SUCCESS'
      AND cdc.environment = 'PRODUCTION'
      AND cdc.finished_date >= DATE_SUB(cc.implementation_date, INTERVAL cc.baseline_period_days DAY)
)
SELECT 
    period,
    COUNT(*) as total_deployments,
    COUNT(*) / baseline_period_days as daily_avg
FROM deployments
GROUP BY period;
```

---

## 7. Configuration

### 7.1 Blueprint Configuration

```json
[
  [
    {
      "plugin": "copilot",
      "subtasks": null,
      "options": {
        "connectionId": 1,
        "scopeId": "octodemo"
      }
    }
  ]
]
```

### 7.2 Scope Configuration

```json
{
  "organization": "octodemo",
  "team": null,
  "implementationDate": "2023-08-29",
  "baselinePeriodDays": 90
}
```

---

## 8. Error Handling

| Error | Handling |
|-------|----------|
| 403 Forbidden | Check PAT permissions, log clear message |
| 404 Not Found | Org doesn't exist or no Copilot subscription |
| 422 Metrics Disabled | Copilot Metrics API not enabled for org |
| Empty response (team) | Log warning, skip (< 5 users privacy) |
| Rate limit (429) | Respect Retry-After header |

---

## 9. Testing Requirements

### 9.1 Unit Tests
- Connection validation
- API response parsing
- Metric aggregation logic

### 9.2 E2E Tests
- Full collection with mock API responses
- CSV fixture comparisons

### 9.3 Integration Tests (with real API)
- Test org: `octodemo`
- Verify data structure matches expectations

---

## 10. Open Questions

1. **Team-level granularity priority?**
   - Adds complexity, requires team→repo mapping
   - Skip for MVP?

2. **Language metrics storage?**
   - Could explode table size (30+ languages × 5+ editors × 100 days)
   - Store aggregated only, or detailed?

3. **Seat data collection frequency?**
   - Daily? Weekly? Only on first run?
   - Needed for adoption timeline

4. **Domain layer conversion?**
   - Stay in tool layer like Q Dev?
   - Or map to domain models for DORA integration?

---

## 11. Implementation Phases

### Phase 1: Core Plugin (MVP)
- [ ] Connection model & CRUD
- [ ] Scope model (org-level only)
- [ ] Metrics collector (daily aggregates)
- [ ] Basic Grafana dashboard (adoption)

### Phase 2: Impact Analysis
- [ ] Seat data collection for adoption date
- [ ] Implementation date configuration
- [ ] Impact dashboard with before/after SQL

### Phase 3: Enhanced Features
- [ ] Team-level metrics (optional)
- [ ] Language/editor breakdown storage
- [ ] Repository-level PR correlation

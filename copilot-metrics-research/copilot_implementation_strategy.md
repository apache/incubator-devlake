# GitHub Copilot Plugin - Implementation Strategy (Option B: Repo/Project-Level)

## Date: December 5, 2025

---

## Core Philosophy

> "Developer productivity is better measured at the **team/repo/project level**"

Instead of correlating individual user Copilot usage with their personal productivity, we:
1. Track **aggregate Copilot adoption** at org/team level
2. Track **aggregate productivity metrics** (PR cycle time, deployment frequency) at repo level
3. Compare **before vs. after Copilot implementation date**

---

## The "Copilot Implementation Date" Problem

This is the **critical variable** for measuring impact:

```
Timeline:  [------------ No Copilot ------------|--- Copilot Enabled ---|]
                                                 ^
                                    Copilot Implementation Date
                                    
Metrics:   [--- Baseline Period ---|--- Post-Copilot Period ---|]
```

### How to Determine Implementation Date

**Option A: API-Derived (Automatic)**
```go
// Use seat management API to find earliest seat creation
// GET /orgs/{org}/copilot/billing/seats
// Look for MIN(seat.created_at) across all seats
```

**Option B: User-Configured (Recommended for MVP)**
```go
// Add to connection or scope config:
type CopilotScopeConfig struct {
    ImplementationDate time.Time `json:"implementationDate"`
    BaselinePeriodDays int       `json:"baselinePeriodDays" default:"90"` // Compare 90 days before
}
```

**Option C: Hybrid**
- Auto-detect from API as suggestion
- Allow manual override in UI

---

## Data Model Design

### Level 1: Organization/Team Daily Metrics
```go
// _tool_copilot_org_metrics
type CopilotOrgMetrics struct {
    ConnectionId          uint64
    OrgName               string    `gorm:"primaryKey"`
    Date                  time.Time `gorm:"primaryKey"`
    
    // Engagement
    TotalActiveUsers      int
    TotalEngagedUsers     int
    
    // Code Completions (aggregated across all editors/languages)
    CompletionSuggestions int
    CompletionAcceptances int
    CompletionLinessuggested int
    CompletionLinesAccepted int
    CompletionAcceptanceRate float64 // calculated
    
    // Chat
    ChatTotal             int
    ChatInsertions        int
    ChatCopies            int
    
    // PR Features
    PRSummariesCreated    int
}
```

### Level 2: Repository-Level PR Metrics (from Copilot API)
```go
// _tool_copilot_repo_pr_metrics
type CopilotRepoPRMetrics struct {
    ConnectionId        uint64
    RepoFullName        string    `gorm:"primaryKey"` // e.g., "demo/repo1"
    Date                time.Time `gorm:"primaryKey"`
    
    TotalEngagedUsers   int
    PRSummariesCreated  int
}
```

### Level 3: Adoption Timeline (for before/after analysis)
```go
// _tool_copilot_adoption_config
type CopilotAdoptionConfig struct {
    ConnectionId        uint64    `gorm:"primaryKey"`
    OrgName             string    `gorm:"primaryKey"`
    ImplementationDate  time.Time // When Copilot was rolled out
    BaselinePeriodDays  int       // How many days before to compare
}
```

---

## The Impact Correlation Strategy

### No User-Level Joins Required!

Instead of joining on users, we join on **time periods**:

```sql
-- Example: PR Cycle Time Before vs After Copilot
WITH copilot_dates AS (
    SELECT implementation_date
    FROM _tool_copilot_adoption_config
    WHERE org_name = 'octodemo'
),
pr_metrics AS (
    SELECT 
        CASE 
            WHEN pr.created_at < cd.implementation_date THEN 'Before Copilot'
            ELSE 'After Copilot'
        END AS period,
        AVG(ppm.pr_cycle_time) AS avg_cycle_time,
        AVG(ppm.pr_coding_time) AS avg_coding_time,
        AVG(ppm.pr_review_time) AS avg_review_time,
        COUNT(*) as pr_count
    FROM project_pr_metrics ppm
    JOIN pull_requests pr ON ppm.pr_id = pr.id
    CROSS JOIN copilot_dates cd
    GROUP BY period
)
SELECT * FROM pr_metrics;
```

### Time-Series Overlay Dashboard

```sql
-- Copilot adoption curve overlaid with DORA metrics
SELECT 
    cm.date,
    cm.total_engaged_users,
    cm.completion_acceptance_rate,
    -- Join with DORA daily snapshots
    d.deployment_frequency,
    d.lead_time_for_changes,
    d.change_failure_rate
FROM _tool_copilot_org_metrics cm
LEFT JOIN dora_daily_metrics d ON cm.date = d.date AND cm.org_name = d.project
ORDER BY cm.date;
```

---

## Implementation Changes from Q Dev Pattern

| Q Dev Approach | Copilot Approach |
|----------------|------------------|
| Data from S3 CSV files | Data from GitHub REST API |
| User-level granularity | Org/Team-level (Option B) |
| AWS credentials | GitHub PAT with org read |
| No time correlation | Before/After implementation date |
| Standalone dashboards | Correlation with existing DORA |

### Plugin Structure
```
backend/plugins/copilot/
├── api/
│   ├── connection.go       # CRUD for connections
│   ├── init.go
│   └── test_connection.go  # Verify PAT and org access
├── impl/
│   └── impl.go             # Plugin implementation
├── models/
│   ├── connection.go       # GitHub API credentials
│   ├── org_metrics.go      # Daily org-level metrics
│   ├── repo_pr_metrics.go  # Repo-level PR features
│   ├── adoption_config.go  # Implementation date config
│   └── migrationscripts/
├── tasks/
│   ├── api_client.go       # GitHub API wrapper
│   ├── metrics_collector.go    # Collect from /copilot/metrics
│   ├── seats_collector.go      # Collect seat info (for adoption date)
│   └── task_data.go
```

---

## Key Differences in Connection Model

```go
// Q Dev Connection (AWS)
type QDevConn struct {
    AccessKeyId     string
    SecretAccessKey string
    Region          string
    Bucket          string
    IdentityStoreId string
}

// Copilot Connection (GitHub)
type CopilotConn struct {
    Token        string `json:"token"`        // PAT or GitHub App token
    Organization string `json:"organization"` // Target org name
    // Optional: for team-level granularity
    Teams        []string `json:"teams,omitempty"`
    // Rate limiting
    RateLimitPerHour int `json:"rateLimitPerHour" default:"5000"`
}
```

---

## Dashboard Vision

### Dashboard 1: "Copilot Adoption & Usage" (standalone)
- Daily active/engaged users trend
- Acceptance rate over time
- Language breakdown
- Editor breakdown
- PR summaries created per repo

### Dashboard 2: "Copilot Impact on Engineering Velocity" 🎯
| Panel | Visualization |
|-------|---------------|
| **Before/After Summary** | Stat cards: "Avg PR Cycle Time: 4.2 days → 2.8 days (-33%)" |
| **PR Cycle Time Trend** | Time series with vertical line at implementation date |
| **Deployment Frequency** | Bar chart before/after comparison |
| **Mean Time to Merge** | Before/after histogram |
| **Code Review Time** | Trend with Copilot adoption overlay |

### Dashboard 3: "Repository-Level Copilot Impact"
- Select specific repo
- Compare that repo's PR metrics before/after
- Copilot PR summaries usage for that repo

---

## Open Questions

1. **Team Granularity**: Should we support team-level metrics (`/orgs/{org}/team/{team_slug}/copilot/metrics`)? Would need to map teams to repos.

2. **Repo Attribution**: The Copilot metrics API only shows repo-level data for PR summaries. For code completions, there's no repo breakdown - only language/editor. Is this sufficient?

3. **Multiple Implementation Dates**: What if different teams adopted Copilot at different times? Do we need per-team adoption dates?

4. **API Limits**: The metrics API has a **100-day lookback limit**. For long-term trending, we need to collect and store data incrementally.

---

## Next Steps

1. ✅ Test the Copilot Metrics API with `octodemo` org
2. Validate the actual response structure matches documentation
3. Design the full data model
4. Decide on team-level vs org-level granularity
5. Build the plugin skeleton

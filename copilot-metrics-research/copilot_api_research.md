# GitHub Copilot API Research Notes

## Date: December 5, 2025

---

## API Endpoints Summary

There are **3 main API categories** for GitHub Copilot data:

### 1. Copilot Metrics API (Aggregated - Best for Option B)
**Endpoint**: `GET /orgs/{org}/copilot/metrics` or `GET /enterprises/{enterprise}/copilot/metrics`

**Key Features**:
- Returns **aggregated daily metrics** (not per-user)
- Data for up to **100 days prior**
- Processed once per day for previous day
- Requires 5+ licensed users to return data (privacy threshold)
- **PERFECT for repo/team-level analysis (Option B)**

**Response Structure** (per day):
```json
{
  "date": "2024-06-24",
  "total_active_users": 24,
  "total_engaged_users": 20,
  "copilot_ide_code_completions": {
    "total_engaged_users": 20,
    "languages": [
      { "name": "python", "total_engaged_users": 10 },
      { "name": "ruby", "total_engaged_users": 10 }
    ],
    "editors": [
      {
        "name": "vscode",
        "total_engaged_users": 13,
        "models": [
          {
            "name": "default",
            "is_custom_model": false,
            "total_engaged_users": 13,
            "languages": [
              {
                "name": "python",
                "total_engaged_users": 6,
                "total_code_suggestions": 249,
                "total_code_acceptances": 123,
                "total_code_lines_suggested": 225,
                "total_code_lines_accepted": 135
              }
            ]
          }
        ]
      }
    ]
  },
  "copilot_ide_chat": {
    "total_engaged_users": 13,
    "editors": [
      {
        "name": "vscode",
        "models": [
          {
            "total_chats": 45,
            "total_chat_insertion_events": 12,
            "total_chat_copy_events": 16
          }
        ]
      }
    ]
  },
  "copilot_dotcom_chat": {
    "total_engaged_users": 14,
    "models": [{ "total_chats": 38 }]
  },
  "copilot_dotcom_pull_requests": {
    "total_engaged_users": 12,
    "repositories": [
      {
        "name": "demo/repo1",
        "total_engaged_users": 8,
        "models": [
          {
            "total_pr_summaries_created": 6,
            "total_engaged_users": 8
          }
        ]
      }
    ]
  }
}
```

**Scopes**: Team-level (`/orgs/{org}/team/{team_slug}/copilot/metrics`) and Org-level available

**Authentication**: 
- PAT (classic): `manage_billing:copilot`, `read:org`, or `read:enterprise`
- Fine-grained: "GitHub Copilot Business" org permissions (read)

---

### 2. Copilot Usage Metrics API (JSON Downloads - Enterprise only)
**Endpoints**:
- `GET /enterprises/{enterprise}/copilot/metrics/reports/enterprise-28-day/latest` - 28-day aggregate
- `GET /enterprises/{enterprise}/copilot/metrics/reports/enterprise-1-day?day=YYYY-MM-DD` - Single day
- `GET /enterprises/{enterprise}/copilot/metrics/reports/users-28-day/latest` - User-level 28-day
- `GET /enterprises/{enterprise}/copilot/metrics/reports/users-1-day?day=YYYY-MM-DD` - User-level single day

**Response**: Returns **download links** to JSON files (signed URLs with expiration)
```json
{
  "download_links": [
    "https://example.com/copilot-usage-report-1.json",
    "https://example.com/copilot-usage-report-2.json"
  ],
  "report_start_day": "2025-07-01",
  "report_end_day": "2025-07-28"
}
```

**Note**: Reports available starting **October 10, 2025**, historical data up to 1 year

---

### 3. Copilot User Management API (Seat/Adoption data)
**Endpoints**:
- `GET /orgs/{org}/copilot/billing` - Org-level seat summary
- `GET /orgs/{org}/copilot/billing/seats` - List all seat assignments
- `GET /orgs/{org}/members/{username}/copilot` - Individual user seat details

**Seat Summary Response**:
```json
{
  "seat_breakdown": {
    "total": 12,
    "added_this_cycle": 9,
    "pending_invitation": 0,
    "pending_cancellation": 0,
    "active_this_cycle": 12,
    "inactive_this_cycle": 11
  },
  "seat_management_setting": "assign_selected",
  "ide_chat": "enabled",
  "platform_chat": "enabled",
  "cli": "enabled",
  "public_code_suggestions": "block",
  "plan_type": "business"
}
```

**Individual Seat Response**:
```json
{
  "created_at": "2021-08-03T18:00:00-06:00",
  "updated_at": "2021-09-23T15:00:00-06:00",
  "pending_cancellation_date": null,
  "last_activity_at": "2021-10-14T00:53:32-06:00",
  "last_activity_editor": "vscode/1.77.3/copilot/1.86.82",
  "plan_type": "business",
  "assignee": {
    "login": "octocat",
    "id": 1
  },
  "assigning_team": {
    "name": "Justice League",
    "slug": "justice-league"
  }
}
```

---

## Key Metrics Available

### IDE Code Completions
| Metric | Description |
|--------|-------------|
| `total_code_suggestions` | Number of suggestions shown |
| `total_code_acceptances` | Number of suggestions accepted |
| `total_code_lines_suggested` | Lines of code suggested |
| `total_code_lines_accepted` | Lines of code accepted |
| `total_engaged_users` | Users who interacted |

### IDE Chat
| Metric | Description |
|--------|-------------|
| `total_chats` | Total chat conversations |
| `total_chat_insertion_events` | Code inserted from chat |
| `total_chat_copy_events` | Code copied from chat |

### GitHub.com Chat
| Metric | Description |
|--------|-------------|
| `total_chats` | Chats on github.com |
| `total_engaged_users` | Users engaged |

### Pull Requests (Copilot for PRs)
| Metric | Description |
|--------|-------------|
| `total_pr_summaries_created` | PR summaries generated |
| `total_engaged_users` | Users using PR features |
| **By Repository** | Breakdown available per repo! |

---

## Critical Insight: Repository-Level Data Available! 🎯

The `copilot_dotcom_pull_requests` section includes **repository breakdown**:
```json
"repositories": [
  {
    "name": "demo/repo1",
    "total_engaged_users": 8,
    "models": [{ "total_pr_summaries_created": 6 }]
  }
]
```

This is PERFECT for Option B (repo-level correlation)!

---

## Implementation Considerations

### For Option B (Repo/Project-Level Analysis):

1. **Primary Data Source**: Copilot Metrics API (`/orgs/{org}/copilot/metrics`)
   - Aggregated data, no individual user tracking needed
   - Team-level granularity available
   - Repository-level PR metrics available

2. **Scope Definition**: 
   - Scope = Organization (or Team within org)
   - Store daily aggregates with date range

3. **Copilot Implementation Date Tracking**:
   - Use seat management API to find earliest `created_at` for seats
   - Or allow manual configuration of "Copilot rollout date"

4. **Correlation Strategy**:
   - Join Copilot daily metrics with DevLake's `project_pr_metrics`
   - Compare:
     - PR cycle time before/after Copilot adoption date
     - Deployment frequency trends
     - Code review time changes

### Required Permissions
- PAT scope: `manage_billing:copilot` OR `read:org`
- Org owner or billing manager access needed

---

## Next Steps

1. ✅ Test API with `octodemo` org using provided GH_PAT
2. Design data model based on metrics response structure
3. Determine how to store/flatten the nested JSON structure
4. Plan Grafana dashboard SQL for before/after comparisons

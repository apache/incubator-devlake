# GitHub Copilot Plugin - Research Summary

## Status: Ready for Spec Driven Development

All research is complete. This document summarizes findings across all research files.

---

## 1. API Research Complete

### Tested Endpoints (octodemo org)

| Endpoint | Status | Notes |
|----------|--------|-------|
| `GET /orgs/{org}/copilot/billing` | ✅ Works | Returns seat breakdown, plan type |
| `GET /orgs/{org}/copilot/metrics` | ✅ Works | Daily metrics with detailed breakdown |
| `GET /orgs/{org}/copilot/billing/seats` | ✅ Works | Lists all seat assignments |
| `GET /orgs/{org}/team/{team}/copilot/metrics` | ⚠️ Empty | Requires 5+ licensed users per team |

### Key API Findings

1. **100-day lookback limit** - Historical data beyond 100 days not available
2. **Privacy threshold** - Team metrics require 5+ users
3. **Daily data** - Metrics processed next-day
4. **PR data includes repo names** - Critical for repo-level correlation
5. **Seat `created_at` dates** - Can derive implementation date automatically

### Sample Response Structure (Metrics API)

```json
{
  "copilot_ide_code_completions": {
    "total_engaged_users": 184,
    "languages": [{"name": "typescript", ...}],
    "editors": [{"name": "vscode", ...}]
  },
  "copilot_ide_chat": {
    "total_engaged_users": 33,
    "editors": [{"name": "vscode", "models": [...]}]
  },
  "copilot_dotcom_chat": {
    "total_engaged_users": 27,
    "models": [{"name": "default", ...}]
  },
  "copilot_dotcom_pull_requests": {
    "total_engaged_users": 12,
    "repositories": [{"name": "repo-name", ...}]
  },
  "date": "2024-12-02"
}
```

---

## 2. DevLake Domain Models for Correlation

### Available PR Metrics (project_pr_metrics table)

| Field | Description | Unit |
|-------|-------------|------|
| `pr_cycle_time` | First commit → PR merged | Minutes |
| `pr_coding_time` | First commit → PR created | Minutes |
| `pr_pickup_time` | PR created → First review | Minutes |
| `pr_review_time` | First review → PR merged | Minutes |
| `pr_deploy_time` | PR merged → Deployed | Minutes |

### DORA Metrics Available

- Deployment Frequency (from `cicd_deployment_commits`)
- Lead Time for Changes (uses `pr_cycle_time` from PRs deployed)
- Change Failure Rate (from incidents linked to deployments)
- Mean Time to Recovery (from incident resolution)

---

## 3. Strategic Decisions Made

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Data granularity | Org/Team level (not user) | Better for team productivity, privacy-friendly |
| Correlation method | Time-based (before/after) | Avoids user-level joins |
| Implementation date | Configurable per scope | Enables baseline comparison |
| MVP scope | Org-level only | Team-level adds complexity |
| Language breakdown | Optional detail table | Avoid table bloat |

---

## 4. Files Created

| File | Purpose |
|------|---------|
| [copilot_api_research.md](copilot_api_research.md) | API documentation notes |
| [copilot_implementation_strategy.md](copilot_implementation_strategy.md) | Option B strategy details |
| [copilot_api_actual_responses.md](copilot_api_actual_responses.md) | Live API test results |
| [copilot_api_schemas.md](copilot_api_schemas.md) | **Complete API response schemas** |
| [copilot_plugin_spec.md](copilot_plugin_spec.md) | **Main spec document** |
| [avocado_corp_enterprise_28_day.json](avocado_corp_enterprise_28_day.json) | Sample enterprise metrics data |
| [avocado_corp_users_28_day.json](avocado_corp_users_28_day.json) | Sample user metrics data |

---

## 5. Reference Implementation

The Q Dev plugin (`backend/plugins/q_dev/`) provides the closest pattern:

- AWS S3-based data collection (we use GitHub REST API)
- User-level metrics (we aggregate to org/team level)
- No domain layer conversion (stays in tool tables)
- Two subtasks: collect + extract (we'll have billing + metrics + seats)

---

## 6. Open Questions (To Discuss in SDD Session)

1. **Team-level priority?** - Skip for MVP, add later if needed
2. **Language detail storage?** - Start with aggregates only
3. **Seat collection frequency?** - Once per run is sufficient
4. **Domain layer mapping?** - Keep in tool layer like Q Dev

---

## 7. Next Steps

1. **Start Spec Driven Development session**
2. **Define acceptance criteria per phase**
3. **Implement Phase 1 (MVP)**:
   - Connection model & CRUD
   - Scope model (org-level)
   - Metrics collector
   - Basic adoption dashboard

The [copilot_plugin_spec.md](copilot_plugin_spec.md) file contains all implementation details.

# GitHub Copilot API - Actual Response Analysis

## Date: December 5, 2025
## Test Org: `octodemo` (Enterprise plan, 872 seats)

---

## API Test Results

### 1. Copilot Billing API (`GET /orgs/{org}/copilot/billing`)

**Response:**
```json
{
  "seat_breakdown": {
    "pending_invitation": 0,
    "pending_cancellation": 9,
    "added_this_cycle": 6,
    "total": 872,
    "active_this_cycle": 627,
    "inactive_this_cycle": 245
  },
  "seat_management_setting": "assign_all",
  "plan_type": "enterprise",
  "public_code_suggestions": "allow",
  "ide_chat": "enabled",
  "cli": "enabled",
  "platform_chat": "enabled"
}
```

**Key Fields for Plugin:**
- `seat_breakdown.total` - Total licensed seats
- `seat_breakdown.active_this_cycle` - Active users this billing cycle
- `seat_breakdown.inactive_this_cycle` - Inactive users
- `plan_type` - "business" or "enterprise"
- Feature flags: `ide_chat`, `cli`, `platform_chat`

---

### 2. Copilot Metrics API (`GET /orgs/{org}/copilot/metrics`)

**Top-Level Properties per Day:**
| Property | Description |
|----------|-------------|
| `date` | ISO date string (YYYY-MM-DD) |
| `total_active_users` | Users with any Copilot activity |
| `total_engaged_users` | Users who interacted with suggestions |
| `copilot_ide_code_completions` | Code completion metrics |
| `copilot_ide_chat` | IDE chat metrics |
| `copilot_dotcom_chat` | GitHub.com chat metrics |
| `copilot_dotcom_pull_requests` | PR summary metrics |

**Sample Daily Data (Dec 4, 2025):**
- Active Users: 362
- Engaged Users: 332
- PR Summary Users: 19

**7-Day Trend:**
```
2025-08-27: Active=362, Engaged=332, PRUsers=19
2025-08-28: Active=365, Engaged=338, PRUsers=16
2025-08-29: Active=318, Engaged=282, PRUsers=10
2025-08-30: Active=78,  Engaged=68,  PRUsers=3   (Weekend)
2025-08-31: Active=62,  Engaged=57,  PRUsers=1   (Weekend)
2025-09-01: Active=173, Engaged=158, PRUsers=12
2025-09-02: Active=338, Engaged=301, PRUsers=16
```

---

### 3. Code Completions Detail

**Structure:** `copilot_ide_code_completions.editors[].models[].languages[]`

**Per-Language Metrics:**
```json
{
  "name": "python",
  "total_engaged_users": 11,
  "total_code_suggestions": 315,
  "total_code_acceptances": 119,
  "total_code_lines_suggested": 465,
  "total_code_lines_accepted": 92
}
```

**Editors Seen:**
- vscode (dominant)
- JetBrains
- VisualStudio
- Neovim
- Vim

**Languages (30+ languages tracked):**
- typescript, typescriptreact, javascript
- python, go, ruby, java, csharp
- markdown (surprisingly high usage)
- yaml, json, css, html
- Many more...

---

### 4. IDE Chat Metrics

**Structure:** `copilot_ide_chat.editors[].models[]`

```json
{
  "name": "default",
  "total_chats": 4704,
  "is_custom_model": false,
  "total_engaged_users": 148,
  "total_chat_copy_events": 195,
  "total_chat_insertion_events": 9
}
```

**Key Insight:** `total_chat_insertion_events` (code inserted directly) vs `total_chat_copy_events` (copied manually)

---

### 5. GitHub.com Chat Metrics

**Structure:** `copilot_dotcom_chat.models[]`

```json
{
  "name": "default",
  "total_chats": 1229,
  "is_custom_model": false,
  "total_engaged_users": 183
}
```

---

### 6. PR Summary Metrics (Critical for Option B!)

**Structure:** `copilot_dotcom_pull_requests.repositories[]`

```json
{
  "repositories": [
    {
      "name": "octodemo/copilot_agent_mode-glorious-octo-waffle",
      "models": [
        {
          "name": "default",
          "is_custom_model": false,
          "total_engaged_users": 1,
          "total_pr_summaries_created": 1
        }
      ],
      "total_engaged_users": 1
    },
    {
      "name": "",  // NOTE: Empty string for unattributed PRs
      "models": [...],
      "total_engaged_users": 8
    }
  ],
  "total_engaged_users": 10
}
```

**Key Insight:** Repository-level breakdown available! Some PRs have empty repo name (privacy/attribution issue?).

---

### 7. Seat Assignment Details (`GET /orgs/{org}/copilot/billing/seats`)

**Per-Seat Response:**
```json
{
  "created_at": "2023-08-29T02:50:42+03:00",   // When seat was assigned
  "assignee": {
    "login": "nathos",
    "id": 4215,
    "type": "User"
  },
  "pending_cancellation_date": null,
  "plan_type": "enterprise",
  "last_authenticated_at": "2025-12-04T15:53:22Z",
  "updated_at": "2024-02-01T03:00:00+03:00",
  "last_activity_at": "2025-11-06T19:12:15+03:00",
  "last_activity_editor": "copilot_pr_review"  // or "vscode/1.106.3/copilot-chat/0.33.4"
}
```

**Key Fields for Adoption Date Detection:**
- `created_at` - When this user got Copilot access
- `last_activity_at` - Last usage timestamp
- `last_activity_editor` - Which feature/editor used last

**For octodemo org:** Earliest seat `created_at` is **2023-08-29** (Copilot implementation date)

---

### 8. Team-Level Metrics

**Endpoint:** `GET /orgs/{org}/team/{team_slug}/copilot/metrics`

**Tested Result:** Empty response for tested teams (need 5+ licensed users per team per day)

**Note:** Team granularity works but requires teams with sufficient Copilot users.

---

## Data Model Requirements Summary

### Core Tables Needed:

1. **`_tool_copilot_org_daily_metrics`** - Daily org-level aggregates
   - date, org, total_active, total_engaged
   - completion_suggestions, completion_acceptances
   - completion_lines_suggested, completion_lines_accepted
   - chat_total, chat_insertions, chat_copies
   - dotcom_chat_total
   - pr_summaries_total, pr_engaged_users

2. **`_tool_copilot_language_metrics`** - Per-language breakdown (optional)
   - date, org, language, editor
   - suggestions, acceptances, lines metrics

3. **`_tool_copilot_repo_pr_metrics`** - Repo-level PR summaries
   - date, org, repo_full_name
   - pr_summaries_created, engaged_users

4. **`_tool_copilot_seat_info`** - For adoption tracking
   - org, user_login, created_at, last_activity_at
   - (Could be used to derive implementation date)

5. **`_tool_copilot_org_config`** - Plugin configuration
   - org, implementation_date (user-configurable)
   - baseline_period_days

---

## API Limits & Considerations

| Limit | Value |
|-------|-------|
| Historical data | 100 days max |
| Privacy threshold | 5+ users per day/team |
| Rate limit | Standard GitHub rate limits |
| Data freshness | Daily, processed overnight |

---

## Key Insights for Implementation

1. **No custom model data** - `is_custom_model` always `false` in octodemo
2. **PR repo names can be empty** - Need to handle gracefully
3. **Weekend patterns visible** - Active users drops significantly on weekends
4. **Chat insertions are low** - Most users copy rather than insert
5. **Markdown is heavily used** - Copilot used a lot for documentation
6. **Seat data gives adoption timeline** - `created_at` tells us when each user got access

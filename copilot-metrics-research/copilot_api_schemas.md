# GitHub Copilot API Response Schemas

## Generated from Actual API Responses

---

## 1. Enterprise Daily Metrics Schema

**Endpoint:** `GET /enterprises/{enterprise}/copilot/metrics` (with download link)

```json
{
  "report_start_day": "string (YYYY-MM-DD)",
  "report_end_day": "string (YYYY-MM-DD)",
  "enterprise_id": "string",
  "created_at": "string (timestamp)",
  "day_totals": [
    {
      "day": "string (YYYY-MM-DD)",
      "enterprise_id": "string",
      "daily_active_users": "integer",
      "weekly_active_users": "integer",
      "monthly_active_users": "integer",
      "monthly_active_chat_users": "integer",
      "monthly_active_agent_users": "integer",
      "user_initiated_interaction_count": "integer",
      "code_generation_activity_count": "integer",
      "code_acceptance_activity_count": "integer",
      "loc_suggested_to_add_sum": "integer",
      "loc_suggested_to_delete_sum": "integer",
      "loc_added_sum": "integer",
      "loc_deleted_sum": "integer",
      "totals_by_ide": "array<TotalsByIde>",
      "totals_by_feature": "array<TotalsByFeature>",
      "totals_by_language_feature": "array<TotalsByLanguageFeature>",
      "totals_by_language_model": "array<TotalsByLanguageModel>",
      "totals_by_model_feature": "array<TotalsByModelFeature>"
    }
  ]
}
```

---

## 2. User Daily Metrics Schema

**Endpoint:** `GET /enterprises/{enterprise}/copilot/metrics/users` (with download link)

```json
{
  "report_start_day": "string (YYYY-MM-DD)",
  "report_end_day": "string (YYYY-MM-DD)",
  "day": "string (YYYY-MM-DD)",
  "enterprise_id": "string",
  "user_id": "integer",
  "user_login": "string",
  "user_initiated_interaction_count": "integer",
  "code_generation_activity_count": "integer",
  "code_acceptance_activity_count": "integer",
  "loc_suggested_to_add_sum": "integer",
  "loc_suggested_to_delete_sum": "integer",
  "loc_added_sum": "integer",
  "loc_deleted_sum": "integer",
  "used_agent": "boolean",
  "used_chat": "boolean",
  "totals_by_ide": "array<UserTotalsByIde>",
  "totals_by_feature": "array<TotalsByFeature>",
  "totals_by_language_feature": "array<TotalsByLanguageFeature>",
  "totals_by_language_model": "array<TotalsByLanguageModel>",
  "totals_by_model_feature": "array<TotalsByModelFeature>"
}
```

---

## 3. Nested Object Schemas

### TotalsByIde

```json
{
  "ide": "string (e.g., 'vscode', 'intellij', 'neovim', 'visualstudio')",
  "user_initiated_interaction_count": "integer",
  "code_generation_activity_count": "integer",
  "code_acceptance_activity_count": "integer",
  "loc_suggested_to_add_sum": "integer",
  "loc_suggested_to_delete_sum": "integer",
  "loc_added_sum": "integer",
  "loc_deleted_sum": "integer"
}
```

### UserTotalsByIde (extends TotalsByIde)

```json
{
  "ide": "string",
  "user_initiated_interaction_count": "integer",
  "code_generation_activity_count": "integer",
  "code_acceptance_activity_count": "integer",
  "loc_suggested_to_add_sum": "integer",
  "loc_suggested_to_delete_sum": "integer",
  "loc_added_sum": "integer",
  "loc_deleted_sum": "integer",
  "last_known_plugin_version": {
    "sampled_at": "string (ISO 8601 timestamp)",
    "plugin": "string (e.g., 'copilot-chat', 'copilot', 'copilot-intellij')",
    "plugin_version": "string"
  },
  "last_known_ide_version": {
    "sampled_at": "string (ISO 8601 timestamp)",
    "ide_version": "string"
  }
}
```

### TotalsByFeature

```json
{
  "feature": "string (enum: see Feature Values below)",
  "user_initiated_interaction_count": "integer",
  "code_generation_activity_count": "integer",
  "code_acceptance_activity_count": "integer",
  "loc_suggested_to_add_sum": "integer",
  "loc_suggested_to_delete_sum": "integer",
  "loc_added_sum": "integer",
  "loc_deleted_sum": "integer"
}
```

### TotalsByLanguageFeature

```json
{
  "language": "string (e.g., 'typescript', 'python', 'ruby', 'go', 'others')",
  "feature": "string (enum: see Feature Values below)",
  "code_generation_activity_count": "integer",
  "code_acceptance_activity_count": "integer",
  "loc_suggested_to_add_sum": "integer",
  "loc_suggested_to_delete_sum": "integer",
  "loc_added_sum": "integer",
  "loc_deleted_sum": "integer"
}
```

### TotalsByLanguageModel

```json
{
  "language": "string",
  "model": "string (e.g., 'claude-4.5-sonnet', 'gpt-5.0', 'unknown', 'auto')",
  "code_generation_activity_count": "integer",
  "code_acceptance_activity_count": "integer",
  "loc_suggested_to_add_sum": "integer",
  "loc_suggested_to_delete_sum": "integer",
  "loc_added_sum": "integer",
  "loc_deleted_sum": "integer"
}
```

### TotalsByModelFeature

```json
{
  "model": "string",
  "feature": "string (enum: see Feature Values below)",
  "user_initiated_interaction_count": "integer",
  "code_generation_activity_count": "integer",
  "code_acceptance_activity_count": "integer",
  "loc_suggested_to_add_sum": "integer",
  "loc_suggested_to_delete_sum": "integer",
  "loc_added_sum": "integer",
  "loc_deleted_sum": "integer"
}
```

---

## 4. Enumeration Values

### Feature Values (observed in actual data)

| Feature | Description |
|---------|-------------|
| `code_completion` | Inline code completions/suggestions |
| `chat_panel_agent_mode` | Chat panel in agent mode (agentic coding) |
| `chat_panel_ask_mode` | Chat panel in ask/question mode |
| `chat_panel_custom_mode` | Chat panel in custom mode |
| `chat_panel_edit_mode` | Chat panel in edit mode |
| `chat_panel_unknown_mode` | Chat panel mode unknown (often IntelliJ) |
| `agent_edit` | Agent-driven code edits (applied changes) |
| `chat_inline` | Inline chat (Ctrl+I / Cmd+I) |

### IDE Values (observed in actual data)

| IDE | Plugin Name |
|-----|-------------|
| `vscode` | copilot, copilot-chat |
| `intellij` | copilot-intellij |
| `visualstudio` | (Visual Studio) |
| `neovim` | copilot.vim |

### Model Values (observed in actual data)

| Model | Description |
|-------|-------------|
| `unknown` | Model not specified/tracked |
| `auto` | Auto-selected model |
| `claude-4.5-sonnet` | Anthropic Claude 4.5 Sonnet |
| `claude-4.5-haiku` | Anthropic Claude 4.5 Haiku |
| `claude-4.0-sonnet` | Anthropic Claude 4.0 Sonnet |
| `claude-opus-4` | Anthropic Claude Opus 4 |
| `claude-opus-4.5` | Anthropic Claude Opus 4.5 |
| `gpt-4.1` | OpenAI GPT-4.1 |
| `gpt-4o` | OpenAI GPT-4o |
| `gpt-5.0` | OpenAI GPT-5.0 |
| `gpt-5.1` | OpenAI GPT-5.1 |
| `gpt-5.1-codex` | OpenAI GPT-5.1 Codex |
| `gpt-5-mini` | OpenAI GPT-5 Mini |
| `gpt-5-codex` | OpenAI GPT-5 Codex |
| `gemini-3.0-pro` | Google Gemini 3.0 Pro |
| `others` | Other models |

### Language Values (sample observed in actual data)

- `typescript`, `typescriptreact`, `javascript`
- `python`, `ruby`, `go`, `java`, `csharp`, `kotlin`
- `markdown`, `json`, `jsonc`, `yaml`, `xml`
- `css`, `html`, `tsx`
- `bash`, `shellscript`, `powershell`, `zsh`
- `dockerfile`, `terraform`, `bicep`
- `sql`, `kusto`, `proto3`
- `unknown`, `others`

---

## 5. Organization Billing Schema

**Endpoint:** `GET /orgs/{org}/copilot/billing`

```json
{
  "seat_breakdown": {
    "pending_invitation": "integer",
    "pending_cancellation": "integer",
    "added_this_cycle": "integer",
    "total": "integer",
    "active_this_cycle": "integer",
    "inactive_this_cycle": "integer"
  },
  "seat_management_setting": "string (e.g., 'assign_all', 'assign_selected')",
  "plan_type": "string ('business' | 'enterprise')",
  "public_code_suggestions": "string ('allow' | 'block')",
  "ide_chat": "string ('enabled' | 'disabled')",
  "cli": "string ('enabled' | 'disabled')",
  "platform_chat": "string ('enabled' | 'disabled')"
}
```

---

## 6. Seat Assignment Schema

**Endpoint:** `GET /orgs/{org}/copilot/billing/seats`

```json
{
  "total_seats": "integer",
  "seats": [
    {
      "created_at": "string (ISO 8601 timestamp)",
      "assignee": {
        "login": "string",
        "id": "integer",
        "type": "string ('User' | 'Team')"
      },
      "pending_cancellation_date": "string | null",
      "plan_type": "string ('business' | 'enterprise')",
      "last_authenticated_at": "string (ISO 8601 timestamp)",
      "updated_at": "string (ISO 8601 timestamp)",
      "last_activity_at": "string (ISO 8601 timestamp)",
      "last_activity_editor": "string"
    }
  ]
}
```

---

## 7. Organization Metrics Schema

**Endpoint:** `GET /orgs/{org}/copilot/metrics`

```json
[
  {
    "date": "string (YYYY-MM-DD)",
    "total_active_users": "integer",
    "total_engaged_users": "integer",
    "copilot_ide_code_completions": {
      "total_engaged_users": "integer",
      "editors": [
        {
          "name": "string",
          "total_engaged_users": "integer",
          "models": [
            {
              "name": "string",
              "is_custom_model": "boolean",
              "total_engaged_users": "integer",
              "languages": [
                {
                  "name": "string",
                  "total_engaged_users": "integer",
                  "total_code_suggestions": "integer",
                  "total_code_acceptances": "integer",
                  "total_code_lines_suggested": "integer",
                  "total_code_lines_accepted": "integer"
                }
              ]
            }
          ]
        }
      ]
    },
    "copilot_ide_chat": {
      "total_engaged_users": "integer",
      "editors": [
        {
          "name": "string",
          "total_engaged_users": "integer",
          "models": [
            {
              "name": "string",
              "is_custom_model": "boolean",
              "total_engaged_users": "integer",
              "total_chats": "integer",
              "total_chat_copy_events": "integer",
              "total_chat_insertion_events": "integer"
            }
          ]
        }
      ]
    },
    "copilot_dotcom_chat": {
      "total_engaged_users": "integer",
      "models": [
        {
          "name": "string",
          "is_custom_model": "boolean",
          "total_engaged_users": "integer",
          "total_chats": "integer"
        }
      ]
    },
    "copilot_dotcom_pull_requests": {
      "total_engaged_users": "integer",
      "repositories": [
        {
          "name": "string",
          "total_engaged_users": "integer",
          "models": [
            {
              "name": "string",
              "is_custom_model": "boolean",
              "total_engaged_users": "integer",
              "total_pr_summaries_created": "integer"
            }
          ]
        }
      ]
    }
  }
]
```

---

## 8. Key Metrics Definitions

| Metric | Definition |
|--------|------------|
| `daily_active_users` | Users who had any Copilot activity on that day |
| `weekly_active_users` | Users active in the past 7 days (rolling) |
| `monthly_active_users` | Users active in the past 28 days (rolling) |
| `user_initiated_interaction_count` | Chat messages, inline prompts initiated by user |
| `code_generation_activity_count` | Number of code suggestions/generations made |
| `code_acceptance_activity_count` | Number of suggestions accepted by user |
| `loc_suggested_to_add_sum` | Lines of code suggested for addition |
| `loc_suggested_to_delete_sum` | Lines of code suggested for deletion |
| `loc_added_sum` | Lines of code actually added (accepted) |
| `loc_deleted_sum` | Lines of code actually deleted (accepted) |
| `used_agent` | User utilized agent mode during the day |
| `used_chat` | User utilized chat features during the day |

---

## 9. Data Availability Notes

- **Historical data:** Up to 28 days for downloadable metrics
- **Data freshness:** Daily, processed overnight
- **Privacy threshold:** Metrics require 5+ users for aggregation
- **Download format:** JSON Lines (one record per line for user metrics)
- **Rate limits:** Standard GitHub API rate limits apply

# Phase 0 Research — GitHub Copilot Adoption Plugin

Date: December 11, 2025  
Branch: `001-copilot-metrics-plugin`

## Decisions

### 1. Adoption Metrics Granularity
- **Decision**: Track both active and engaged user counts as primary adoption trend lines.
- **Rationale**: The Copilot Metrics API exposes `total_active_users` and `total_engaged_users`; engaged users better reflect productive usage.
- **Alternatives Considered**: Active-only trends (insufficient insight into value); seat counts (lagging indicator, not usage-driven).

### 2. Language Breakdown Scope
- **Decision**: Limit default language chart to the top 10 languages per date range.
- **Rationale**: Balances readability and insight; avoids clutter from >30 languages while retaining secondary stack visibility.
- **Alternatives Considered**: Top 5 (too coarse for polyglot orgs); show all languages (hurts dashboard usability).

### 3. Scope Level for Phase 1
- **Decision**: Support organization-level metrics only; defer team-level collection to Phase 2.
- **Rationale**: Team metrics require additional mapping, and GitHub enforces a ≥5 user privacy threshold which often returns empty responses.
- **Alternatives Considered**: Implement org + team concurrently (increases complexity/risk without MVP value).

### 4. Seat Assignment Identifiers
- **Decision**: Store raw GitHub usernames in `_tool_copilot_seats` for adoption timeline analytics.
- **Rationale**: Aligns with existing DevLake data (raw usernames for commits/PRs) and enables correlation without lookup tables.
- **Alternatives Considered**: Hashing/anonymization (breaks correlation); configurable hashing (adds complexity in MVP phase).

### 5. PR Summary Metrics
- **Decision**: Exclude PR summary metrics from Phase 1 scope.
- **Rationale**: Focus on adoption; PR summaries correspond to code review agent usage, which is lower priority and outside adoption MVP.
- **Alternatives Considered**: Include PR summaries (requires extra storage + dashboards, dilutes MVP timeline).

### 6. Coding Agent & Code Review Agent Metrics
- **Decision**: Do not surface dedicated "Coding Agent" or "Code Review Agent" metrics in Phase 1.
- **Rationale**: GitHub Copilot Metrics API does not expose explicit agent identifiers; interactions appear as chat or completion metrics.
- **Alternatives Considered**: Infer agent usage from chat metadata (unreliable, unsupported by documented API).

## Supporting Notes

- API endpoints confirmed: `GET /orgs/{org}/copilot/metrics`, `GET /orgs/{org}/copilot/billing`, `GET /orgs/{org}/copilot/billing/seats`.
- Collection model: stateful incremental collector with 100-day lookback; use stored `last_collected_date`.
- Privacy threshold: if API omits data due to <5 engaged users, log warning and store no record for that day.
- Rate limiting: reuse DevLake API helper with per-request token bucket and Retry-After handling.
- Grafana dashboards: will source data from tool tables using SQL queries referencing `copilot_org_metrics`, `copilot_language_metrics`, and `copilot_seats`.

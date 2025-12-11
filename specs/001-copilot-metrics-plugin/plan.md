# Implementation Plan: GitHub Copilot Plugin - Adoption Metrics (Phase 1)

**Branch**: `001-copilot-metrics-plugin` | **Date**: December 11, 2025 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-copilot-metrics-plugin/spec.md`

## Summary

Deliver a new GitHub Copilot data source plugin that ingests organization-level adoption metrics (code completions, chat usage, seat activity) from the Copilot REST API, persists them in DevLake tool-layer tables, and exposes an Adoption Dashboard in Grafana with active/engaged user trends, top languages, editor distribution, and chat usage panels. The plugin will support connection CRUD, org-scope configuration, incremental collection within the 100-day API window, and seat tracking using raw GitHub usernames for accurate adoption analysis. Each collector/API surface will implement explicit error handling for 403/404/422/429 responses with retry/backoff guidance so user-facing workflows remain resilient. Recently documented enterprise and per-user download schemas introduce additional metrics (daily/weekly/monthly actives, agent usage, LOC deltas), but these remain out-of-scope for Phase 1 and are captured as future expansion notes below.

## API Surface & Scope Notes

- **Supported in Phase 1**: `GET /orgs/{org}/copilot/metrics`, `GET /orgs/{org}/copilot/billing`, `GET /orgs/{org}/copilot/billing/seats` (JSON responses, 100-day rolling window). These endpoints provide active/engaged users, editor/language breakdowns, chat usage, and seat assignments.
- **Deferred**: Enterprise download endpoints (`/enterprises/{enterprise}/copilot/metrics` and `/metrics/users`) that supply daily/weekly/monthly actives, agent mode usage, LOC change counts, and per-user breakdowns. Integration requires handling JSONL downloads and 28-day history limits; captured for future roadmap (Phase 2+).
- **Feature Mapping**: Recorded schema enums (e.g., `code_completion`, `chat_panel_agent_mode`, `agent_edit`) inform naming conventions for Grafana panels and potential filter chips but will not be fully persisted until enterprise support is prioritized.
- **Error Handling**: Connections and collectors must surface clear messaging for 403/404/422/429 responses and honor `Retry-After` headers; unit/E2E tests will simulate these cases to keep pipelines stable.

## Technical Context

**Language/Version**: Go 1.20 (DevLake backend); Grafana JSON dashboards; SQL for visualization queries  
**Primary Dependencies**: DevLake plugin helper packages (`helpers/pluginhelper/apihelper`, `helpers/pluginhelper/dal`), `net/http` for GitHub REST calls, `gorm.io/gorm`, `github.com/tidwall/gjson` for JSON parsing, Grafana dashboard definitions  
**Storage**: MySQL or PostgreSQL via DevLake DAL (Raw `_raw_copilot_*` tables, Tool `_tool_copilot_*` tables)  
**Testing**: Go unit tests (`go test ./...`), plugin E2E CSV fixtures under `backend/plugins/copilot/e2e`, Grafana dashboard JSON lint via `yarn dashlint` (existing workflow)  
**Target Platform**: DevLake backend server (Docker-compose / Linux) with Grafana front-end  
**Project Type**: Backend data source plugin + Grafana dashboards (monorepo)  
**Performance Goals**: Collect up to 100 days of daily metrics per run (<5 minutes per org, <500 API calls) and render dashboard panels within 5 seconds for 1-year time ranges  
**Constraints**: Respect GitHub rate limits (5,000 requests/hour), handle privacy threshold (≥5 engaged users per day) gracefully, no storage of team-level metrics in Phase 1  
**Scale/Scope**: Support organizations with 10k+ Copilot seats, storing one daily record per metric category (~5 tables) and language/editor breakdown for 30+ languages

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Plugin Independence (Pass)**: Plan introduces `backend/plugins/copilot` with complete interface coverage (meta, task, model, migration, source). No cross-plugin imports planned.
- **Three-Layer Data Model (Pass)**: Collection flows Raw → Tool with incremental collectors, extractors, and optional converters (none to domain in Phase 1). Migrations will create `_tool_copilot_*` tables.
- **Test-Driven Development (Pass)**: Commit to unit tests for collectors/extractors and CSV-backed E2E tests replicating sample GitHub responses. Table info registration coverage ensured.
- **Migration-First Schema Changes (Pass)**: All schema definitions added through new migration scripts registered in `models/migrationscripts/register.go`.
- **Apache Compliance (Pass)**: New Go files include ASF license header; dependencies remain compliant (GitHub REST only). No violations expected.
- **Post-Design Recheck**: No new violations introduced during Phase 1 design; constitution gates remain satisfied.

## Project Structure

### Documentation (this feature)

```text
specs/001-copilot-metrics-plugin/
├── plan.md          # Implementation plan (this file)
├── research.md      # Phase 0 research consolidation
├── data-model.md    # Phase 1 entity definitions
├── quickstart.md    # Setup walkthrough for plugin adoption
├── contracts/
│   └── copilot.openapi.yaml  # REST contract for plugin endpoints
└── tasks.md         # Generated later via /speckit.tasks (not in this step)
```

### Source Code (repository root)

```text
backend/
└── plugins/
    └── copilot/
        ├── api/
        ├── impl/
        ├── models/
        │   └── migrationscripts/
        ├── tasks/
        ├── e2e/
        └── README.md

grafana/
└── dashboards/
    └── copilot/
        └── adoption.json

backend/helpers/
└── pluginhelper/
    └── (reuse existing helpers; no structural change)
```

**Structure Decision**: Implement a dedicated `backend/plugins/copilot` Go plugin mirroring existing data-source patterns—**explicitly modeled on** the `backend/plugins/q_dev` (Quick Development) plugin for structure, interface wiring, migration layout, and task registration. Dashboard assets live under `grafana/dashboards/copilot/` (SQL snippets embedded in `adoption.json` rather than a separate `queries/` directory). No additional top-level services are required.

## Complexity Tracking

No constitution violations identified; complexity tracking not required.

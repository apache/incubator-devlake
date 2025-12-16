# Tasks: GitHub Copilot Plugin - Adoption Metrics (Phase 1)

**Input**: Design documents from `/specs/001-copilot-metrics-plugin/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Establish plugin scaffold aligned with the `q_dev` data-source pattern

- [x] T001 Create `backend/plugins/copilot/` directory skeleton (api/, impl/, models/, tasks/, e2e/, README.md) mirroring `backend/plugins/q_dev/`
- [x] T002 [P] Add `grafana/dashboards/copilot/` folder with placeholder `adoption.json` seeded from q_dev dashboard conventions

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core plumbing required before implementing any user story

- [x] T003 Implement plugin meta + task registration in `backend/plugins/copilot/impl/impl.go` modeled on `backend/plugins/q_dev/impl/impl.go`
- [x] T004 [P] Register Copilot plugin in `backend/plugins/register/register.go` and ensure inclusion in build pipeline (`backend/plugins/all_plugins.go`)
- [x] T005 Define option/connection helper structs in `backend/plugins/copilot/impl/options.go` and shared constants referencing GitHub endpoints
- [x] T006 Create tool-layer models for `CopilotConnection`, `CopilotScope`, `CopilotOrgMetrics`, `CopilotLanguageMetrics`, `CopilotSeat` in `backend/plugins/copilot/models/`
- [x] T007 Author initial migration script `backend/plugins/copilot/models/migrationscripts/20250100_initialize.go` and register via `models/migrationscripts/register.go`
- [x] T008 [P] Populate `backend/plugins/copilot/models/models.go` with `GetTablesInfo()` metadata and add unit coverage similar to `backend/plugins/q_dev/models/models_test.go`
- [x] T009 Establish task registry bootstrap in `backend/plugins/copilot/tasks/register.go` (define subtasks, meta ordering, blueprint default)

**Checkpoint**: Copilot plugin compiles, migrations apply, and plugin appears in DevLake plugin registry

---

## Phase 3: User Story 1 – Configure Copilot Connection (Priority: P1) 🎯 MVP

**Goal**: Administrators can create, test, and manage Copilot connections/scopes via REST

**Independent Test**: Using DevLake API/UI, create a Copilot connection with PAT and confirm `Test Connection` returns billing data

### Implementation

- [x] T010 [P] [US1] Implement API bootstrap in `backend/plugins/copilot/api/init.go` (connection helper, validators) following q_dev init pattern
- [x] T011 [P] [US1] Create connection + scope REST handlers in `backend/plugins/copilot/api/connection.go` and `api/scope.go` with DTOs from `contracts/copilot.openapi.yaml`
- [x] T012 [P] [US1] Implement `TestConnection` helper hitting `/orgs/{org}/copilot/billing` in `backend/plugins/copilot/impl/connection_helper.go`
- [x] T013 [US1] Add validation/unit tests mirroring `backend/plugins/q_dev/api/connection_test.go` to cover PAT/org edge cases
- [x] T014 [US1] Update plugin README (`backend/plugins/copilot/README.md`) and quickstart connection steps to document required PAT scopes
- [x] T014a [US1] Implement graceful error handling for connection and `TestConnection` endpoints, including unit tests simulating 403/404/422/429 responses with user-facing messages

**Checkpoint**: Connection CRUD & test APIs operational; no data collection yet

---

## Phase 4: User Story 2 – Collect Daily Copilot Usage Metrics (Priority: P1)

**Goal**: Automatically ingest daily adoption metrics (active/engaged users, completions, chat usage, seats)

**Independent Test**: Run blueprint → confirm `_tool_copilot_org_metrics` and `_tool_copilot_seats` populated for new days only

### Implementation

- [ ] T015 [P] [US2] Build metrics collector (`backend/plugins/copilot/tasks/metrics_collector.go`) using `helper.NewStatefulApiCollector` against `/orgs/{org}/copilot/metrics`
- [x] T015 [P] [US2] Build metrics collector (`backend/plugins/copilot/tasks/metrics_collector.go`) using `helper.NewStatefulApiCollector` against `/orgs/{org}/copilot/metrics`
- [x] T016 [P] [US2] Implement seat assignment collector (`backend/plugins/copilot/tasks/seat_collector.go`) calling `/orgs/{org}/copilot/billing/seats`
- [x] T017 [P] [US2] Create extractor/convertor mapping raw payloads to `CopilotOrgMetrics` + `CopilotSeat` in `backend/plugins/copilot/tasks/metrics_extractor.go`
- [x] T018 [US2] Register subtasks in `backend/plugins/copilot/tasks/register.go` and wire into plugin pipeline order (collector → extractor)
- [x] T019 [P] [US2] Add E2E fixture set under `backend/plugins/copilot/e2e/metrics/` with mocked JSON (org metrics + seats) and golden CSVs
- [x] T020 [US2] Write unit tests for state bookmarking + rate limit handling in `backend/plugins/copilot/tasks/metrics_collector_test.go`
- [x] T020a [US2] Harden collectors with Retry-After/backoff logic and verify graceful handling of 403/404/422/429 responses via unit tests and log assertions

**Checkpoint**: Metrics & seat data captured incrementally and ready for dashboard consumption

---

## Phase 5: User Story 3 – View Copilot Adoption Dashboard (Priority: P1)

**Goal**: Provide Grafana dashboard visualizing adoption trends (active vs engaged users, acceptance rate, chat usage, seat timeline)

**Independent Test**: Grafana `Copilot Adoption` dashboard loads, panels render data for collected metrics, and date filters work

### Implementation

- [ ] T021 [P] [US3] Author `grafana/dashboards/copilot/adoption.json` with panels for active/engaged trends, acceptance rate, chat usage, seat timeline, leveraging DevLake SQL macros
- [ ] T022 [P] [US3] Add supporting SQL query snippets (e.g., acceptance rate CTEs) directly within `grafana/dashboards/copilot/adoption.json` or existing macro snippets, avoiding extra directory structure
- [ ] T023 [US3] Update `specs/001-copilot-metrics-plugin/quickstart.md` with dashboard navigation + variable instructions and verify via local Grafana sync

**Checkpoint**: Adoption dashboard deliverable complete; foundation for Phase 2 impact analytics

---

## Phase 6: User Story 4 – Collect Language & Editor Breakdown (Priority: P2)

**Goal**: Surface language/editor usage insights in dashboards using detailed Copilot metrics

**Independent Test**: Language + editor panels display top 10 languages, editor share, and respond to date filters after rerunning pipeline

### Implementation

- [ ] T024 [P] [US4] Extend extractor (`backend/plugins/copilot/tasks/metrics_extractor.go`) to persist `CopilotLanguageMetrics` from `editors[].languages[]`
- [ ] T025 [US4] Enhance dashboard JSON (`grafana/dashboards/copilot/adoption.json`) with Top 10 language bar chart + editor distribution pie chart
- [ ] T026 [US4] Add test coverage (E2E fixture + assertions) for language/editor rows in `backend/plugins/copilot/e2e/metrics/language_breakdown.csv`

**Checkpoint**: Detailed language/editor analytics available while core adoption metrics remain stable

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Final quality, documentation, and release readiness

- [ ] T027 [P] Refresh `backend/plugins/copilot/README.md` and `specs/001-copilot-metrics-plugin/quickstart.md` with final instructions + diagrams
- [ ] T028 [P] Publish blueprint recipe + sample dashboard screenshot in `backend/plugins/copilot/README.md`
- [ ] T029 Execute `make unit-test`, `make e2e-test-go-plugins`, and `make grafana-dashboard-sync` to validate end-to-end
- [ ] T030 Capture upgrade notes + limitations in `specs/001-copilot-metrics-plugin/plan.md` and ensure docs mention deferred enterprise endpoints

---

## Dependencies & Execution Order

### Phase Dependencies
- **Setup (Phase 1)** → prerequisite for Foundational work
- **Foundational (Phase 2)** → must finish before any user story tasks (Phases 3–6)
- **User Stories** → execute in priority order (P1 stories first) or parallel once Phase 2 complete
- **Polish (Phase 7)** → after desired user stories are completed

### User Story Dependencies
- **US1**: Requires Foundational migration + API scaffolding
- **US2**: Depends on US1 connection data structures for credentials
- **US3**: Depends on US2 metrics tables being populated
- **US4**: Extends US2 extractor and US3 dashboard; can proceed once base metrics + dashboard exist

### Within Each User Story
- Develop collectors/extractors before wiring dashboards or docs
- Tests (unit/e2e) should be created alongside implementation for regression safety
- Each story is independently verifiable (API endpoints, data ingestion, dashboards)

### Parallel Opportunities
- Setup tasks T001–T002 can run concurrently
- Foundational tasks marked [P] (T004, T005, T007, T009) can be parallelized across developers
- Post-Foundational, US1/US2 can be developed in parallel with clear ownership (API vs ingestion)
- Dashboard (US3) effort can begin once sample metrics exist (use fixtures for local development)

---

## Parallel Example: User Story 2

```bash
# Parallel collectors
Task T015   # metrics collector implementation
Task T016   # seat collector implementation

# Parallel validation
Task T019   # E2E fixtures
Task T020   # collector unit tests
```

---

## Implementation Strategy

### MVP First (Phase 3 focus)
1. Complete Setup + Foundational (Phases 1–2)
2. Ship User Story 1 (connection CRUD + test) → validates API integration
3. Ship User Story 2 (metrics + seats) → enables basic data ingestion
4. Ship User Story 3 (dashboard) → user-facing MVP

### Incremental Enhancements
- User Story 4 adds language/editor insights without altering prior deliverables
- Polish phase finalizes documentation, blueprints, and testing automation

### Future Expansion (Out of Scope)
- Enterprise download endpoints (`/enterprises/{enterprise}/copilot/metrics`, `/users`) requiring JSONL ingest + additional datasets
- Before/after impact analysis (Phase 2 feature `002-copilot-impact-dashboard`)
- Team-level aggregation once GitHub privacy thresholds and mapping are addressed

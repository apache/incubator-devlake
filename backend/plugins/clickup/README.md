<!--
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# ClickUp

The ClickUp plugin collects issues, boards, and sprints from ClickUp so they
feed DevLake's issue-tracking and DORA/velocity metrics, modeled on the Jira
and Linear connectors.

## Authentication

ClickUp authenticates with a **personal API token** (ClickUp → Settings → Apps
→ API Token, starts with `pk_`). The token is sent verbatim in the
`Authorization` header. OAuth is not yet supported.

Create a connection with:

- **Endpoint** — `https://api.clickup.com/api/v2/`
- **Token** — your `pk_...` personal token

## Data scope: the folder is the board

Unlike a raw list, the **scope you select is a ClickUp folder** (e.g. a team's
`Dev Team` / `Sprint Folder`). This mirrors a Jira board: selecting the folder
collects every list inside it on each sync, so rolling and archived sprint
lists never need to be re-scoped.

Domain mapping:

| ClickUp | DevLake domain |
| --- | --- |
| Folder | `board` |
| Sprint list (name-matched) | `sprint` + `board_sprint` |
| Task | `issue` + `board_issue` |
| Task in a sprint list | `sprint_issue` |
| Folder member | `account` |

## Sprints

A list is treated as a sprint when its name matches the **sprint name pattern**
(default `(?i)sprint\s*\d+`, e.g. `v4.3.0 Sprint 40 (7/6/26 - 7/19/26)`). The
start/end dates are parsed from the parenthesised date span in the name;
`M/D/YY` vs `D/M/YY` ordering is disambiguated automatically (a component > 12
must be the day). Lists that don't match are collected as plain board issues.
Archived sprint lists are collected too, so historical velocity is retained.

## Story points

Story points default to ClickUp's native sprint **`points`** field. To read
them from a custom field instead (e.g. a Fibonacci "LOE" field), set
**Story point field** in the scope config to the custom-field name.

## Incidents (DORA Change Failure Rate / MTTR)

ClickUp tasks have no universal "type" field, so incidents are modeled by
**scope**: add the folder that holds incidents (e.g. a "Security Incidents"
folder) as its own board and set the scope config's **Force issue type** to
`INCIDENT`. Every issue on that board is then classified as an incident.

## Scope configuration (transformation)

Per board, the scope config lets you override:

- **Sprint name pattern** — which lists are sprints.
- **Story point field** — native `points` (blank) or a custom field name.
- **Force issue type** — flag the whole board as `REQUIREMENT` / `BUG` /
  `INCIDENT` (leave blank to detect per task).
- **Issue type patterns** — RegExes matched against a task's type; precedence
  is `INCIDENT` > `BUG` > `REQUIREMENT`.
- **Status mapping** — ClickUp statuses are auto-mapped by their `type`
  (`open`/`unstarted` → `TODO`, `custom` → `IN_PROGRESS`, `done`/`closed` →
  `DONE`); list raw status names to override where a team's workflow differs.

## Metrics

A branded **ClickUp** Grafana dashboard (`grafana/dashboards/{mysql,postgresql}/ClickUp.json`)
renders issue throughput and lead-time metrics; DORA and sprint/velocity
metrics are computed from the source-agnostic domain layer.

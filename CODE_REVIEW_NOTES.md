# Code Review Notes: jira ScopeConfigId fix / ExtraJQL ProjectName feature

Review date: 2026-07-30. Scope: two branches based on `main` @ `14a4e5bcd`:

- `fix/jira-blueprint-scope-config` (commit `1a07d1062`)
- `feature/expose-devlake-project-extrajql-variable` (commit `90239e1d9`)

This file records what was found and fixed, so future agents don't have to
re-derive the same context from scratch. All fixes described below were
applied as uncommitted working-tree changes on their respective branch as of
this review; check `git log`/`git status` on each branch to see whether
they've since been committed.

## fix/jira-blueprint-scope-config

**What it does:** `makeDataSourcePipelinePlanV200` in
`backend/plugins/jira/api/blueprint_v200.go` now threads the resolved
`ScopeConfigId` into `JiraTaskOptions`, the way `ArgoCD`/`Linear` already do
in their own `blueprint_v200.go`. Previously Jira relied solely on a runtime
fallback in `impl.go`'s `PrepareTaskData` that re-looks-up the board by
`connection_id + board_id` and copies `scope.ScopeConfigId` if the task
option was zero. That fallback is fragile (e.g. it does nothing if
`op.BoardId == 0`, and is generally an indirect way to get a value that's
already known at plan-build time), so passing it explicitly is a real,
justified fix and matches established plugin convention.

**Findings (fixed):**
1. `gofmt` failure — the new `ScopeConfigId` struct field broke alignment
   of the `JiraTaskOptions{}` literal. The repo's `.golangci.yaml` enables
   the `gofmt` formatter and CI runs `golangci-lint run`, so this would have
   failed CI. Fixed by running `gofmt -w`.
2. Missing regression test — sibling plugins (`linear`, `bitbucket`,
   `gitlab`, `azuredevops_go`) all have a `blueprint_v200_test.go`; Linear's
   even has `TestMakePipelinePlanV200PassesScopeConfigId`, testing exactly
   this scenario. Jira had no such test. Added
   `backend/plugins/jira/api/blueprint_v200_test.go` modeled on Linear's,
   covering: ScopeConfig.ID takes priority over scope.ScopeConfigId, and
   scope.ScopeConfigId is used when ScopeConfig is nil/unconfigured.

## feature/expose-devlake-project-extrajql-variable

**What it does:** Exposes the DevLake *project* name (the logical grouping
of scopes across data sources) as `{{.ProjectName}}` inside the `ExtraJQL`
scope-config template, alongside the existing `{{.BoardName}}` /
`{{.BoardId}}`.

**Important design change made during this review — read before touching
`ProjectName` resolution again:**

The original implementation resolved `ProjectName` in `impl.go`'s
`PrepareTaskData` by reverse-looking-up `project_mapping WHERE
table='boards' AND row_id=<domain id of the board>`. This is broken for the
feature's actual stated purpose (per the requester): **the same Jira board
attached to multiple Devlake projects, each project filtering its own
tickets via ExtraJQL.** `project_mapping`'s primary key is
`(project_name, table, row_id)`, so a shared board has one row *per
project*, and the reverse lookup can't tell which project the *current
pipeline run* belongs to — it just grabs an arbitrary/first row. That means
every project sharing a board would get the same (wrong, for all but one of
them) `{{.ProjectName}}` value, silently.

We looked at fixing this "properly" by threading `projectName` through
`plugin.DataSourcePluginBlueprintV200.MakeDataSourcePipelinePlanV200(...)`,
the interface every datasource plugin implements — but that's a framework
interface change with a large blast radius (~15+ plugins), and was
rejected as too big for this fix.

**Actual fix implemented (small blast radius, no interface changes):**
`backend/server/services/blueprint_makeplan_v200.go`'s `GeneratePlanJsonV200`
already receives `projectName` as a parameter — it's the only place in the
framework where "this pipeline run belongs to project X" and "these are its
scopes" are known together unambiguously. After building `sourcePlans` (the
per-connection plans returned by each plugin's
`MakeDataSourcePipelinePlanV200`), it now generically injects
`task.Options["projectName"] = projectName` into every task, for every
plugin, when `projectName != ""`. This is safe because the mapstructure
decoding used everywhere (`Decode`/`DecodeMapStruct`) silently ignores
unknown map keys by default — plugins that don't declare a matching struct
field simply never see it. (Note: `dora`/`refdiff` already manually put a
`"projectName"` key into their own task options today, via
`MakeMetricPluginPipelinePlanV200`'s own `projectName` parameter — this
generic injection follows the same naming convention, just for
`DataSourcePluginBlueprintV200` plugins that don't get `projectName` on
their interface.)

On the Jira side this actually *simplified* the implementation:
- `backend/plugins/jira/tasks/task_data.go`: `JiraOptions` gained a
  `ProjectName string` field (mapstructure tag `projectName,omitempty`),
  populated purely by decoding task options — no plugin-side logic needed.
- `backend/plugins/jira/impl/impl.go`: `PrepareTaskData` now just does
  `ProjectName: op.ProjectName` when building `JiraTaskData`. The
  `project_mapping`/`crossdomain`/`didgen` reverse-lookup code and its
  imports were deleted entirely — there's no more ambiguity to handle.
- Test coverage: `backend/server/services/blueprint_makeplan_v200_test.go`
  gained `TestMakePlanV200InjectsProjectNameIntoDataSourceTaskOptions`,
  asserting the injected key lands in a data-source plugin's task options
  alongside its own options, using `mock.Anything` for the "org" plugin's
  `MapProject` call (re-registering "org" per-test to avoid cross-test mock
  state leaking via the global `plugin` registry — see that file for why).

**Other findings (fixed):**
1. `gofmt` failure — `JqlTemplateData` struct and the `JqlTemplateData{}`
   literal in `issue_collector.go` used tabs for alignment instead of
   spaces; new `ProjectName` field wasn't aligned either. Fixed with
   `gofmt -w`.
2. Zero test coverage for the new variable — `issue_collector_test.go`'s
   existing `Test_renderExtraJQL` test has a `makeData` helper whose third
   parameter (`_ string`) was unused and looked purpose-built for exactly
   this addition, but nothing wired it up. Fixed: the helper now sets
   `ProjectName` from that parameter, and new subtests exercise
   `{{.ProjectName}}` substitution (present and empty/unmapped board). This
   part of the test is still valid after the `ProjectName`-resolution
   redesign above, since it tests `renderExtraJQL` given an already-built
   `JiraTaskData`, independent of how `ProjectName` gets populated upstream.
3. Config-UI help tooltip
   (`config-ui/src/plugins/register/jira/transformation.tsx`) documented
   only `{{.BoardName}}`/`{{.BoardId}}`. Updated to also mention
   `{{.ProjectName}}` so the feature is discoverable from the UI.

## Verification performed

- `go build ./plugins/jira/...`, `./server/services/...` on both branches
  (fix branch built in an isolated git worktree at the time of review).
- `go test ./plugins/jira/api/... ./plugins/jira/tasks/... ./plugins/jira/impl/... ./server/services/...`
  on both branches. Pre-existing `plugins/jira/e2e` requires `E2E_DB_URL`
  and is expected to fail/skip outside a real DB — unrelated to these
  changes. `go build ./...` at the repo root also fails on unrelated
  pre-existing issues in this sandbox (no `libgit2` for `gitextractor`;
  `plugins/org` and other plugin packages are `main` packages meant to be
  built with `-buildmode=plugin`, not as ordinary binaries) — neither is a
  regression from these branches.
- `gofmt -l` on all touched files, before and after fixes, on both
  branches.
- Generating `backend/mocks/{core,helpers}` via `mockery` was required to
  run `server/services` tests locally (`make mock`, or the two `mockery
  --recursive ...` commands in `backend/Makefile`); that directory is
  gitignored and not checked in.

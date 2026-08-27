/*
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
*/

package e2e

import (
	"testing"
	"time"

	"github.com/apache/incubator-devlake/core/config"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/migration"
	"github.com/apache/incubator-devlake/core/models/domainlayer/code"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	coreMigration "github.com/apache/incubator-devlake/core/models/migrationscripts"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/runner"
	"github.com/apache/incubator-devlake/helpers/e2ehelper"
	"github.com/apache/incubator-devlake/impls/dalgorm"
	"github.com/apache/incubator-devlake/impls/logruslog"
	"github.com/apache/incubator-devlake/plugins/monorepo/models"
	monorepoMigration "github.com/apache/incubator-devlake/plugins/monorepo/models/migrationscripts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// backfillMigrationVersion is 20260810_backfill_sub_project_from_monorepo.go's
// Version(). The concrete script type is unexported (as is this codebase's convention
// for migration scripts), so it is looked up by version through the exported
// plugin.MigrationScript interface returned by coreMigration.All() rather than
// constructed directly.
const backfillMigrationVersion = uint64(20260810100200)

func findMigration(t *testing.T, scripts []plugin.MigrationScript, version uint64) plugin.MigrationScript {
	t.Helper()
	for _, s := range scripts {
		if s.Version() == version {
			return s
		}
	}
	t.Fatalf("migration version %d not found", version)
	return nil
}

// TestMigrationAddsSubProjectColumnsAndTable runs the REAL registered core migration
// scripts (not AutoMigrate on the runtime model, which would hide a mistake in the
// migration itself) against a fresh, isolated database, and asserts that:
//
//  1. pull_requests, pull_request_commits and project_pr_metrics all gain a sub_project
//     column (20260810_add_sub_project_to_pr_and_metrics.go).
//  2. cicd_deployment_subprojects exists (20260810_add_cicd_deployment_subprojects.go).
//
// Requires E2E_DB_URL (runs under `make e2e-test` / `make e2e-test-go-plugins`).
func TestMigrationAddsSubProjectColumnsAndTable(t *testing.T) {
	db := e2ehelper.NewIsolatedMigrationDb(t, "monorepo_migration_columns")
	d := dalgorm.NewDalgorm(db)
	basicRes := runner.CreateBasicRes(config.GetConfig(), logruslog.Global, db)

	migrator, err := migration.NewMigrator(basicRes)
	require.NoError(t, err)
	migrator.Register(coreMigration.All(), "Framework")
	require.NoError(t, migrator.Execute())

	assert.True(t, d.HasColumn(&code.PullRequest{}, "sub_project"), "pull_requests.sub_project should exist")
	assert.True(t, d.HasColumn(&code.PullRequestCommit{}, "sub_project"), "pull_request_commits.sub_project should exist")
	assert.True(t, d.HasTable("cicd_deployment_subprojects"), "cicd_deployment_subprojects table should exist")
	assert.True(t, d.HasColumn(&devops.CicdDeploymentSubproject{}, "sub_project"), "cicd_deployment_subprojects.sub_project should exist")
}

// TestBackfillSubProjectFromMonorepo runs the real core migrations AND the real monorepo
// plugin migrations against a fresh, isolated database, seeds data in the shape an
// existing monorepo-plugin install would have accumulated, re-runs just the backfill
// migration (20260810_backfill_sub_project_from_monorepo.go) a second time by calling its
// Up() directly, and asserts it produces the expected sub_project values - i.e. that
// projects already using the monorepo plugin keep their classification across the
// upgrade instead of reverting to NULL/"All". Re-invoking Up() a second time (the first
// happened, as a no-op, during migrator.Execute() before any monorepo data existed) also
// doubles as an idempotency check: the design doc requires the backfill be safe to re-run.
//
// Requires E2E_DB_URL (runs under `make e2e-test` / `make e2e-test-go-plugins`).
func TestBackfillSubProjectFromMonorepo(t *testing.T) {
	db := e2ehelper.NewIsolatedMigrationDb(t, "monorepo_migration_backfill")
	d := dalgorm.NewDalgorm(db)
	basicRes := runner.CreateBasicRes(config.GetConfig(), logruslog.Global, db)

	migrator, err := migration.NewMigrator(basicRes)
	require.NoError(t, err)
	migrator.Register(coreMigration.All(), "Framework")
	migrator.Register(monorepoMigration.All(), "monorepo")
	require.NoError(t, migrator.Execute())

	// Seed a pull_requests/pull_request_commits row via raw INSERT that omits the
	// sub_project column entirely, so it gets the database's real column default (NULL) -
	// exactly what an ALTER TABLE ADD COLUMN produces for rows that existed before this
	// migration ran, and what the backfill's `WHERE sub_project IS NULL` guard expects.
	// (Going through the GORM model here would write sub_project = '' instead, since
	// code.PullRequest.SubProject is a plain string, not a pointer - a real difference,
	// but not the scenario the backfill migration exists to handle.)
	now := time.Now()
	require.NoError(t, d.Exec(
		"INSERT INTO pull_requests (id, base_repo_id, created_date, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		"pr1", "repo1", now, now, now,
	))
	require.NoError(t, d.Exec(
		"INSERT INTO pull_request_commits (commit_sha, pull_request_id, commit_authored_date, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		"commitA1", "pr1", now, now, now,
	))

	require.NoError(t, d.Create(&models.SubProjectPrMetric{
		ProjectName:   "monorepo",
		PullRequestId: "pr1",
		SubProject:    "serviceA",
	}))
	require.NoError(t, d.Create(&models.SubProjectDeployment{
		ProjectName:      "monorepo",
		SubProject:       "serviceA",
		CicdDeploymentId: "pipeline1",
		CommitSha:        "commitA1",
	}))

	// Re-run the backfill script directly: the copy that ran inside migrator.Execute()
	// above was a no-op because it executed before this seed data existed.
	backfill := findMigration(t, coreMigration.All(), backfillMigrationVersion)
	require.NoError(t, backfill.Up(basicRes))

	var gotPr code.PullRequest
	require.NoError(t, d.First(&gotPr, dal.Where("id = ?", "pr1")))
	assert.Equal(t, "serviceA", gotPr.SubProject)

	var gotCommit code.PullRequestCommit
	require.NoError(t, d.First(&gotCommit, dal.Where("commit_sha = ?", "commitA1")))
	assert.Equal(t, "serviceA", gotCommit.SubProject)

	assert.True(t, d.HasTable("cicd_deployment_subprojects"))
	mappingCount, cErr := d.Count(dal.From("cicd_deployment_subprojects"))
	require.NoError(t, cErr)
	assert.Equal(t, int64(1), mappingCount)

	// Re-running the migration a second time must not error and must not duplicate rows
	// (idempotency, per the design doc's risk mitigation).
	require.NoError(t, backfill.Up(basicRes))
	mappingCount, cErr = d.Count(dal.From("cicd_deployment_subprojects"))
	require.NoError(t, cErr)
	assert.Equal(t, int64(1), mappingCount)
}

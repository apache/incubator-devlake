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

package migrationscripts

import (
	"time"

	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
)

var _ plugin.MigrationScript = (*backfillSubProjectFromMonorepo)(nil)

type backfillSubProjectFromMonorepo struct{}

// Up copies sub_project values that existing monorepo-plugin installs already computed
// into the new core columns/table, so projects that were using the monorepo plugin before
// this release keep their classification instead of reverting to NULL ("All").
//
// Every statement is written with a correlated subquery / NOT EXISTS guard rather than the
// MySQL-only `UPDATE ... JOIN` or Postgres-only `UPDATE ... FROM` forms, so the same SQL
// runs unmodified on both of DevLake's supported databases. Each statement only touches
// rows it has not already touched (sub_project IS NULL / NOT EXISTS), which makes the
// whole migration idempotent and safe to re-run if it is interrupted partway through - the
// batching called out as a risk in the design doc was judged unnecessary for a first
// implementation given that guard, but would be a reasonable follow-up for very large
// instances (see the design doc's risk table).
//
// This is a core migration, so it runs on every DevLake install, including ones that have
// never enabled the monorepo plugin. On those installs the plugin's own tables
// (monorepo_subproject_pr_metrics / monorepo_subproject_deployments) do not exist, so each
// half of the backfill is skipped independently when its source table is absent, rather
// than breaking the migration for every non-monorepo user.
func (script *backfillSubProjectFromMonorepo) Up(basicRes context.BasicRes) errors.Error {
	db := basicRes.GetDal()

	if db.HasTable("monorepo_subproject_pr_metrics") {
		if err := backfillPrSubProjects(db); err != nil {
			return err
		}
	}
	if db.HasTable("monorepo_subproject_deployments") {
		if err := backfillDeploymentSubProjects(db); err != nil {
			return err
		}
	}
	return nil
}

func backfillPrSubProjects(db dal.Dal) errors.Error {
	// 1. pull_requests.sub_project <- monorepo_subproject_pr_metrics.sub_project
	if err := db.Exec(`
		UPDATE pull_requests
		SET sub_project = (
			SELECT m.sub_project FROM monorepo_subproject_pr_metrics m
			WHERE m.pull_request_id = pull_requests.id
		)
		WHERE sub_project IS NULL
		  AND EXISTS (
			SELECT 1 FROM monorepo_subproject_pr_metrics m2
			WHERE m2.pull_request_id = pull_requests.id
		  )
	`); err != nil {
		return errors.Default.Wrap(err, "error backfilling pull_requests.sub_project")
	}

	// 2. pull_request_commits.sub_project <- pull_requests.sub_project
	if err := db.Exec(`
		UPDATE pull_request_commits
		SET sub_project = (
			SELECT pr.sub_project FROM pull_requests pr
			WHERE pr.id = pull_request_commits.pull_request_id
		)
		WHERE sub_project IS NULL
		  AND EXISTS (
			SELECT 1 FROM pull_requests pr2
			WHERE pr2.id = pull_request_commits.pull_request_id
			  AND pr2.sub_project IS NOT NULL
		  )
	`); err != nil {
		return errors.Default.Wrap(err, "error backfilling pull_request_commits.sub_project")
	}

	// 3. project_pr_metrics.sub_project <- pull_requests.sub_project
	if err := db.Exec(`
		UPDATE project_pr_metrics
		SET sub_project = (
			SELECT pr.sub_project FROM pull_requests pr
			WHERE pr.id = project_pr_metrics.id
		)
		WHERE sub_project IS NULL
		  AND EXISTS (
			SELECT 1 FROM pull_requests pr2
			WHERE pr2.id = project_pr_metrics.id
			  AND pr2.sub_project IS NOT NULL
		  )
	`); err != nil {
		return errors.Default.Wrap(err, "error backfilling project_pr_metrics.sub_project")
	}

	return nil
}

func backfillDeploymentSubProjects(db dal.Dal) errors.Error {
	now := time.Now()
	if err := db.Exec(`
		INSERT INTO cicd_deployment_subprojects (project_name, cicd_deployment_id, sub_project, created_at, updated_at)
		SELECT DISTINCT d.project_name, d.cicd_deployment_id, d.sub_project, ?, ?
		FROM monorepo_subproject_deployments d
		WHERE NOT EXISTS (
			SELECT 1 FROM cicd_deployment_subprojects x
			WHERE x.project_name = d.project_name
			  AND x.cicd_deployment_id = d.cicd_deployment_id
			  AND x.sub_project = d.sub_project
		)
	`, now, now); err != nil {
		return errors.Default.Wrap(err, "error backfilling cicd_deployment_subprojects")
	}
	return nil
}

func (*backfillSubProjectFromMonorepo) Version() uint64 {
	return 20260810100200
}

func (*backfillSubProjectFromMonorepo) Name() string {
	return "backfill sub_project into core tables from existing monorepo plugin data"
}

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

package tasks

import (
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/monorepo/models"
)

var UpdateProjectPrMetricsSubProjectMeta = plugin.SubTaskMeta{
	Name:       "updateProjectPrMetricsSubProject",
	EntryPoint: UpdateProjectPrMetricsSubProject,
	// This subtask reads pull_requests.sub_project (written by attributePullRequests) and
	// project_pr_metrics.deployment_commit_id (written by DORA's calculateChangeLeadTime),
	// so both must have already run. Ordering against attributePullRequests is guaranteed
	// by subtask declaration order within this plugin's single stage; ordering against DORA
	// is guaranteed by the stage-padding workaround in impl.go.
	EnabledByDefault: true,
	Description: "Tag project_pr_metrics with the monorepo sub-project of the pull request it belongs to, " +
		"cross-check it against the deployment DORA attributed the PR to, and backfill " +
		"monorepo_subproject_pr_metrics for backward compatibility",
	DomainTypes: []string{plugin.DOMAIN_TYPE_CICD, plugin.DOMAIN_TYPE_CODE_REVIEW},
}

// SubProjectMismatchRow is one PR whose label-derived sub-project disagrees with the
// deploy-job-pattern-derived sub-project of the deployment that actually shipped it.
// Exported so tests (including e2e tests in a different package) can query for mismatches
// directly rather than scraping log output.
type SubProjectMismatchRow struct {
	PrId                 string
	PrSubProject         string
	DeploymentSubProject string
}

// prMetricSubProjectRow is one project_pr_metrics row joined with the pull request's
// sub_project and (if any) the pipeline id of the deployment DORA attributed it to. It is
// the source data for the monorepo_subproject_pr_metrics backfill.
//
// RawDataOrigin is embedded because DataConverter copies that field from the input row
// onto every result; without it the conversion panics.
type prMetricSubProjectRow struct {
	common.RawDataOrigin
	PullRequestId string
	SubProject    string
	CodingTime    *int64
	PickupTime    *int64
	ReviewTime    *int64
	DeployTime    *int64
	CycleTime     *int64
	DeploymentId  string
	PrCreatedDate *time.Time
	PrMergedDate  *time.Time
	DeployedDate  *time.Time
}

// UpdateProjectPrMetricsSubProject is the third and final monorepo subtask. It:
//
//  1. Tags project_pr_metrics.sub_project from pull_requests.sub_project.
//  2. Cross-checks that tag against the sub-project(s) of the deployment that DORA's
//     commit-accurate attribution (project_pr_metrics.deployment_commit_id) says actually
//     shipped the PR, logging - not failing on - any disagreement as a configuration
//     hygiene signal (see the design doc's risk table).
//  3. Backfills monorepo_subproject_pr_metrics for one release, for backward
//     compatibility, using DORA's already-computed numbers rather than recomputing them.
//     This retires the old merge-date-nearest-deployment heuristic that used to live in
//     AttributePullRequests: deploy/cycle time sourced this way are more accurate, so
//     existing monorepo users will see those two columns change (improve) on upgrade.
func UpdateProjectPrMetricsSubProject(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	data := taskCtx.GetData().(*MonorepoTaskData)
	projectName := data.Options.ProjectName

	if err := tagProjectPrMetricsSubProject(db, projectName); err != nil {
		return err
	}
	mismatches, err := FindSubProjectMismatches(db, projectName)
	if err != nil {
		return err
	}
	for _, row := range mismatches {
		logger.Warn(nil,
			"monorepo: pull request %s is labelled for sub-project %q but was shipped by a deployment "+
				"attributed to sub-project %q - check prLabels/deployJobPattern configuration",
			row.PrId, row.PrSubProject, row.DeploymentSubProject)
	}
	if len(mismatches) > 0 {
		logger.Info("monorepo: found %d pull request(s) with a label/deployment sub-project mismatch", len(mismatches))
	}
	return backfillSubProjectPrMetrics(taskCtx, db, projectName)
}

// tagProjectPrMetricsSubProject implements Step 1. It is written as a correlated
// subquery rather than the MySQL-only `UPDATE ... JOIN` form so the same SQL runs
// unmodified on both of DevLake's supported databases.
func tagProjectPrMetricsSubProject(db dal.Dal, projectName string) errors.Error {
	if err := db.Exec(`
		UPDATE project_pr_metrics
		SET sub_project = (
			SELECT pr.sub_project FROM pull_requests pr
			WHERE pr.id = project_pr_metrics.id
		)
		WHERE project_name = ?
	`, projectName); err != nil {
		return errors.Default.Wrap(err, "error tagging project_pr_metrics.sub_project")
	}
	return nil
}

// FindSubProjectMismatches implements Step 2's cross-check query. A mismatch means the
// PR's label-derived sub-project disagrees with the deploy-job-pattern-derived
// sub-project of a deployment that shipped it - almost always a misconfigured
// deployJobPattern/prLabels entry. The caller (UpdateProjectPrMetricsSubProject) only
// logs these as warnings; it never fails the subtask or overrides
// pull_requests/project_pr_metrics.sub_project because of them - labels win, the
// deployment side is only a cross-check.
//
// Note that a PR shipped by a pipeline that deploys several sub-projects (one row per
// sub-project in cicd_deployment_subprojects) will always produce a mismatch row against
// every sub-project other than its own - that is expected noise from the many-to-many
// deployment mapping, not necessarily a misconfiguration.
func FindSubProjectMismatches(db dal.Dal, projectName string) ([]SubProjectMismatchRow, errors.Error) {
	var rows []SubProjectMismatchRow
	err := db.All(&rows,
		dal.Select("ppm.id AS pr_id, pr.sub_project AS pr_sub_project, ds.sub_project AS deployment_sub_project"),
		dal.From("project_pr_metrics ppm"),
		dal.Join("JOIN pull_requests pr ON pr.id = ppm.id"),
		dal.Join("JOIN cicd_deployment_commits dc ON dc.id = ppm.deployment_commit_id"),
		dal.Join("JOIN cicd_deployment_subprojects ds ON ds.cicd_deployment_id = dc.cicd_deployment_id AND ds.project_name = ppm.project_name"),
		dal.Where("ppm.project_name = ? AND pr.sub_project IS NOT NULL AND ds.sub_project <> pr.sub_project", projectName),
	)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error checking for sub-project label/deployment mismatches")
	}
	return rows, nil
}

// backfillSubProjectPrMetrics implements Step 3.
func backfillSubProjectPrMetrics(taskCtx plugin.SubTaskContext, db dal.Dal, projectName string) errors.Error {
	if err := db.Exec(
		"DELETE FROM monorepo_subproject_pr_metrics WHERE project_name = ?",
		projectName,
	); err != nil {
		return errors.Default.Wrap(err, "error deleting previous monorepo_subproject_pr_metrics")
	}

	clauses := []dal.Clause{
		dal.Select(`ppm.id AS pull_request_id, pr.sub_project AS sub_project,
			ppm.pr_coding_time AS coding_time, ppm.pr_pickup_time AS pickup_time,
			ppm.pr_review_time AS review_time, ppm.pr_deploy_time AS deploy_time,
			ppm.pr_cycle_time AS cycle_time, dc.cicd_deployment_id AS deployment_id,
			ppm.pr_created_date AS pr_created_date, ppm.pr_merged_date AS pr_merged_date,
			ppm.pr_deployed_date AS deployed_date`),
		dal.From("project_pr_metrics ppm"),
		dal.Join("JOIN pull_requests pr ON pr.id = ppm.id"),
		dal.Join("LEFT JOIN cicd_deployment_commits dc ON dc.id = ppm.deployment_commit_id"),
		dal.Where("ppm.project_name = ? AND pr.sub_project IS NOT NULL", projectName),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: MonorepoApiParams{
				ProjectName: projectName,
			},
			Table: "project_pr_metrics",
		},
		InputRowType: reflect.TypeOf(prMetricSubProjectRow{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			row := inputRow.(*prMetricSubProjectRow)
			return []interface{}{&models.SubProjectPrMetric{
				ProjectName:   projectName,
				PullRequestId: row.PullRequestId,
				SubProject:    row.SubProject,
				CodingTime:    row.CodingTime,
				PickupTime:    row.PickupTime,
				ReviewTime:    row.ReviewTime,
				DeployTime:    row.DeployTime,
				CycleTime:     row.CycleTime,
				DeploymentId:  row.DeploymentId,
				PrCreatedDate: row.PrCreatedDate,
				PrMergedDate:  row.PrMergedDate,
				DeployedDate:  row.DeployedDate,
			}}, nil
		},
	})
	if err != nil {
		return err
	}
	return converter.Execute()
}

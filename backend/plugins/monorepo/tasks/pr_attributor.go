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
	"math"
	"reflect"
	"sort"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/monorepo/models"
)

var AttributePullRequestsMeta = plugin.SubTaskMeta{
	Name:             "attributePullRequests",
	EntryPoint:       AttributePullRequests,
	EnabledByDefault: true,
	Description:      "Attribute merged pull requests to monorepo sub-projects by label and compute their change lead time",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD, plugin.DOMAIN_TYPE_CODE_REVIEW},
}

// RawDataOrigin is embedded because DataConverter copies that field from the input row
// onto every result; without it the conversion panics.
type pullRequestRow struct {
	common.RawDataOrigin
	Id          string
	CreatedDate time.Time
	MergedDate  *time.Time
}

type prLabelRow struct {
	PullRequestId string
	LabelName     string
}

// deployedAt is the minimum information needed to link a merged PR to a deployment.
type deployedAt struct {
	Id           string
	FinishedDate time.Time
}

func AttributePullRequests(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	data := taskCtx.GetData().(*MonorepoTaskData)
	projectName := data.Options.ProjectName

	if err := db.Exec(
		"DELETE FROM monorepo_subproject_pr_metrics WHERE project_name = ?",
		projectName,
	); err != nil {
		return errors.Default.Wrap(err, "error deleting previous monorepo_subproject_pr_metrics")
	}

	labelsByPr, err := loadPrLabels(db, projectName)
	if err != nil {
		return err
	}
	// DORA already computes coding/pickup/review time correctly for a monorepo: they
	// depend only on the pull request itself. Only the deploy leg needs recomputing.
	doraMetrics, err := loadDoraPrMetrics(db, projectName)
	if err != nil {
		return err
	}
	deploymentsBySubProject, err := loadSubProjectDeployments(db, projectName)
	if err != nil {
		return err
	}
	logger.Info("monorepo: %d labelled PRs, %d DORA metric rows, %d sub-projects with deployments",
		len(labelsByPr), len(doraMetrics), len(deploymentsBySubProject))

	clauses := []dal.Clause{
		dal.Select("pr.id, pr.created_date, pr.merged_date"),
		dal.From("pull_requests pr"),
		dal.Join("JOIN project_mapping pm ON (pm.table = 'repos' AND pm.row_id = pr.base_repo_id)"),
		dal.Where("pm.project_name = ? AND pr.merged_date IS NOT NULL", projectName),
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	unattributed := 0
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: MonorepoApiParams{
				ProjectName: projectName,
			},
			Table: "pull_requests",
		},
		InputRowType: reflect.TypeOf(pullRequestRow{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			pr := inputRow.(*pullRequestRow)
			subProject := data.Matcher.MatchPrLabels(labelsByPr[pr.Id])
			if subProject == "" {
				// No sub-project claims this PR; it is simply out of scope here.
				unattributed++
				return nil, nil
			}

			metric := &models.SubProjectPrMetric{
				ProjectName:   projectName,
				PullRequestId: pr.Id,
				SubProject:    subProject,
				PrCreatedDate: &pr.CreatedDate,
				PrMergedDate:  pr.MergedDate,
			}
			if dm := doraMetrics[pr.Id]; dm != nil {
				metric.CodingTime = dm.PrCodingTime
				metric.PickupTime = dm.PrPickupTime
				metric.ReviewTime = dm.PrReviewTime
			}

			if deployment := firstDeploymentAfter(deploymentsBySubProject[subProject], pr.MergedDate); deployment != nil {
				metric.DeploymentId = deployment.Id
				metric.DeployedDate = &deployment.FinishedDate
				metric.DeployTime = computeTimeSpan(pr.MergedDate, &deployment.FinishedDate)
			}

			// Mirrors DORA's definition: coding + (merged - created) + deploy.
			var cycleTime int64
			if metric.CodingTime != nil {
				cycleTime += *metric.CodingTime
			}
			if prDuring := computeTimeSpan(&pr.CreatedDate, pr.MergedDate); prDuring != nil {
				cycleTime += *prDuring
			}
			if metric.DeployTime != nil {
				cycleTime += *metric.DeployTime
			}
			metric.CycleTime = &cycleTime

			return []interface{}{metric}, nil
		},
	})
	if err != nil {
		return err
	}

	if err := converter.Execute(); err != nil {
		return err
	}
	if unattributed > 0 {
		logger.Info("monorepo: %d merged PRs matched no sub-project label and were skipped", unattributed)
	}
	return nil
}

// firstDeploymentAfter returns the earliest deployment that finished after mergedDate.
//
// This is an approximation: it assumes a merged change is shipped by the next successful
// production deployment of its sub-project. Hotfixes, cherry-picks, rollbacks and re-runs
// can break that assumption. Exact attribution would need each deployment's commit range
// from the refdiff plugin's commits_diffs table; this function is the seam where that
// swap would happen.
func firstDeploymentAfter(deployments []deployedAt, mergedDate *time.Time) *deployedAt {
	if mergedDate == nil || len(deployments) == 0 {
		return nil
	}
	// deployments is sorted by FinishedDate ascending.
	i := sort.Search(len(deployments), func(i int) bool {
		return deployments[i].FinishedDate.After(*mergedDate)
	})
	if i >= len(deployments) {
		return nil
	}
	return &deployments[i]
}

func loadPrLabels(db dal.Dal, projectName string) (map[string][]string, errors.Error) {
	var rows []prLabelRow
	err := db.All(&rows,
		dal.Select("prl.pull_request_id, prl.label_name"),
		dal.From("pull_request_labels prl"),
		dal.Join("JOIN pull_requests pr ON (pr.id = prl.pull_request_id)"),
		dal.Join("JOIN project_mapping pm ON (pm.table = 'repos' AND pm.row_id = pr.base_repo_id)"),
		dal.Where("pm.project_name = ?", projectName),
	)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error loading pull request labels")
	}
	byPr := make(map[string][]string)
	for _, r := range rows {
		byPr[r.PullRequestId] = append(byPr[r.PullRequestId], r.LabelName)
	}
	return byPr, nil
}

func loadDoraPrMetrics(db dal.Dal, projectName string) (map[string]*crossdomain.ProjectPrMetric, errors.Error) {
	var rows []*crossdomain.ProjectPrMetric
	err := db.All(&rows,
		dal.From(&crossdomain.ProjectPrMetric{}),
		dal.Where("project_name = ?", projectName),
	)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error loading project_pr_metrics")
	}
	byPr := make(map[string]*crossdomain.ProjectPrMetric, len(rows))
	for _, r := range rows {
		byPr[r.Id] = r
	}
	return byPr, nil
}

// loadSubProjectDeployments returns the successful production deployments of each
// sub-project, sorted by finish time so they can be searched by merge date.
func loadSubProjectDeployments(db dal.Dal, projectName string) (map[string][]deployedAt, errors.Error) {
	var rows []models.SubProjectDeployment
	err := db.All(&rows,
		dal.From(&models.SubProjectDeployment{}),
		dal.Where(
			"project_name = ? AND result = ? AND environment = ? AND finished_date IS NOT NULL",
			projectName, devops.RESULT_SUCCESS, devops.PRODUCTION,
		),
	)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error loading monorepo_subproject_deployments")
	}
	bySubProject := make(map[string][]deployedAt)
	for _, r := range rows {
		bySubProject[r.SubProject] = append(bySubProject[r.SubProject], deployedAt{
			Id:           r.CicdDeploymentId,
			FinishedDate: *r.FinishedDate,
		})
	}
	for name := range bySubProject {
		list := bySubProject[name]
		sort.Slice(list, func(i, j int) bool {
			return list[i].FinishedDate.Before(list[j].FinishedDate)
		})
		bySubProject[name] = list
	}
	return bySubProject, nil
}

// computeTimeSpan returns the whole minutes between start and end, or nil when either is
// missing or the span is negative. Mirrors the identical unexported helper in the DORA
// plugin (plugins/dora/tasks/change_lead_time_calculator.go) so the two produce the same
// numbers; it cannot be imported because it is not exported there.
func computeTimeSpan(start, end *time.Time) *int64 {
	if start == nil || end == nil {
		return nil
	}
	span := end.Sub(*start)
	minutes := int64(math.Ceil(span.Minutes()))
	if minutes < 0 {
		return nil
	}
	return &minutes
}

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
	"fmt"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

var DeploymentCommitsGeneratorMeta = plugin.SubTaskMeta{
	Name:             "generateDeploymentCommits",
	EntryPoint:       GenerateDeploymentCommits,
	EnabledByDefault: false, // it should be executed before refdiff.calculateDeploymentCommitsDiff, check https://github.com/apache/incubator-devlake/issues/4869 for detail
	Description:      "Generate deployment_commits from cicd_pipeline_commits if cicd_pipeline.type == DEPLOYMENT or any of its cicd_tasks is a deployment task",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

type pipelineCommitEx struct {
	devops.CiCDPipelineCommit
	PipelineName       string
	Result             string
	Status             string
	OriginalStatus     string
	OriginalResult     string
	DurationSec        *float64
	QueuedDurationSec  *float64
	CreatedDate        *time.Time
	FinishedDate       *time.Time
	Environment        string
	CicdScopeId        string
	HasTestingTasks    bool
	HasStagingTasks    bool
	HasProductionTasks bool
}

func GenerateDeploymentCommits(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*DoraTaskData)
	// select all cicd_pipeline_commits from all "Deployments" in the project
	// Note that failed records shall be included as well
	noneSkippedResult := []string{devops.RESULT_FAILURE, devops.RESULT_SUCCESS}
	var clauses = []dal.Clause{
		dal.Select(
			`
				pc.*,
				p.name as pipeline_name,
				p.result,
				p.status,
				p.duration_sec,
				p.queued_duration_sec,
				p.created_date,
				p.finished_date,
				p.environment,
				p.cicd_scope_id,
				p.original_status,
				p.original_result,
				EXISTS(SELECT 1 FROM cicd_tasks t WHERE t.pipeline_id = p.id AND t.environment = ? AND t.result IN ?) as has_testing_tasks,
				EXISTS(SELECT 1 FROM cicd_tasks t WHERE t.pipeline_id = p.id AND t.environment = ? AND t.result IN ?) as has_staging_tasks,
				EXISTS( SELECT 1 FROM cicd_tasks t WHERE t.pipeline_id = p.id AND t.environment = ? AND t.result IN ?) as has_production_tasks
			`,
			devops.TESTING, noneSkippedResult,
			devops.STAGING, noneSkippedResult,
			devops.PRODUCTION, noneSkippedResult,
		),
		dal.From("cicd_pipeline_commits pc"),
		dal.Join("LEFT JOIN cicd_pipelines p ON (p.id = pc.pipeline_id)"),
		dal.Where(
			`
			p.result IN ? AND (
				p.type = ? OR EXISTS(SELECT 1 FROM cicd_tasks t WHERE t.pipeline_id = p.id AND t.type = ? AND t.result IN ?)
			)
			`,
			noneSkippedResult,
			devops.DEPLOYMENT,
			devops.DEPLOYMENT,
			noneSkippedResult,
		),
	}
	if data.Options.ScopeId != nil {
		clauses = append(clauses, dal.Where(`p.cicd_scope_id = ?`, data.Options.ScopeId))
	} else {
		clauses = append(clauses,
			dal.Join("LEFT JOIN project_mapping pm ON (pm.table = 'cicd_scopes' AND pm.row_id = p.cicd_scope_id)"),
			dal.Where(`pm.project_name = ?`, data.Options.ProjectName),
		)
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	count, err := db.Count(clauses...)
	if err != nil {
		return errors.Default.Wrap(err, "error getting count of clauses")
	}
	if count == 0 {
		// Clear previous results from the project
		deleteSql := fmt.Sprintf("DELETE cdc FROM cicd_deployment_commits cdc LEFT JOIN project_mapping pm ON (pm.table = 'cicd_scopes' AND pm.row_id = cdc.cicd_scope_id) WHERE pm.project_name = '%s'", data.Options.ProjectName)
		err := db.Exec(deleteSql)
		if err != nil {
			return errors.Default.Wrap(err, "error deleting previous cicd_deployment_commits")
		}
		return nil
	}

	enricher, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: DoraApiParams{
				ProjectName: data.Options.ProjectName,
			},
			Table: "cicd_pipeline_commits",
		},
		InputRowType: reflect.TypeOf(pipelineCommitEx{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			pipelineCommit := inputRow.(*pipelineCommitEx)

			domainDeployCommit := &devops.CicdDeploymentCommit{
				DomainEntity: domainlayer.DomainEntity{
					Id: fmt.Sprintf("%s:%s", pipelineCommit.PipelineId, pipelineCommit.RepoUrl),
				},
				CicdScopeId:      pipelineCommit.CicdScopeId,
				CicdDeploymentId: pipelineCommit.PipelineId,
				Name:             pipelineCommit.PipelineName,
				Result:           pipelineCommit.Result,
				OriginalStatus:   pipelineCommit.OriginalStatus,
				OriginalResult:   pipelineCommit.OriginalResult,
				Status:           pipelineCommit.Status,
				Environment:      pipelineCommit.Environment,
				TaskDatesInfo: devops.TaskDatesInfo{
					CreatedDate:  *pipelineCommit.CreatedDate,
					FinishedDate: pipelineCommit.FinishedDate,
				},
				DurationSec:       pipelineCommit.DurationSec,
				QueuedDurationSec: pipelineCommit.QueuedDurationSec,
				CommitSha:         pipelineCommit.CommitSha,
				RefName:           pipelineCommit.Branch,
				RepoId:            pipelineCommit.RepoId,
				RepoUrl:           pipelineCommit.RepoUrl,
			}
			if pipelineCommit.FinishedDate != nil && pipelineCommit.DurationSec != nil {
				s := pipelineCommit.FinishedDate.Add(-time.Duration(*pipelineCommit.DurationSec) * time.Second)
				domainDeployCommit.StartedDate = &s
			}
			// it is tricky when Environment was declared on the cicd_tasks level
			// lets talk common sense and assume that one pipeline can only be deployed to one environment
			// so if the pipeline has both staging and production tasks, we will treat it as a production pipeline
			// and if it has staging tasks without production tasks, we will treat it as a staging pipeline
			// and then a testing pipeline
			// lastly, we will leave Environment empty if any of the above measures didn't work out

			// However, there is another catch, what if one deployed multiple TESTING(STAGING or PRODUCTION)
			// environments? e.g. testing1, testing2, etc., Does it matter?
			if pipelineCommit.Environment == "" {
				if pipelineCommit.HasProductionTasks {
					domainDeployCommit.Environment = devops.PRODUCTION
				} else if pipelineCommit.HasStagingTasks {
					domainDeployCommit.Environment = devops.STAGING
				} else if pipelineCommit.HasTestingTasks {
					domainDeployCommit.Environment = devops.TESTING
				}
			}
			return []interface{}{domainDeployCommit}, nil
		},
	})
	if err != nil {
		return err
	}

	return enricher.Execute()
}

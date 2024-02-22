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
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

var DeploymentGeneratorMeta = plugin.SubTaskMeta{
	Name:             "generateDeployments",
	EntryPoint:       GenerateDeployment,
	EnabledByDefault: true,
	Description:      "Generate cicd_deployments from cicd_pipelines if cicd_pipeline.type == DEPLOYMENT or any of its cicd_tasks is a deployment task",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

type pipelineEx struct {
	devops.CICDPipeline
	HasTestingTasks    bool
	HasStagingTasks    bool
	HasProductionTasks bool
}

func GenerateDeployment(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*DoraTaskData)
	// Note that failed records shall be included as well
	noneSkippedResult := []string{devops.RESULT_FAILURE, devops.RESULT_SUCCESS}
	var clauses = []dal.Clause{
		dal.Select(
			`
				p.*,
				EXISTS(SELECT 1 FROM cicd_tasks t WHERE t.pipeline_id = p.id AND t.environment = ? AND t.result IN ?)
				as has_testing_tasks,
				EXISTS(SELECT 1 FROM cicd_tasks t WHERE t.pipeline_id = p.id AND t.environment = ? AND t.result IN ?)
				as has_staging_tasks,
				EXISTS(SELECT 1 FROM cicd_tasks t WHERE t.pipeline_id = p.id AND t.environment = ? AND t.result IN ?)
				as has_production_tasks
			`,
			devops.TESTING, noneSkippedResult,
			devops.STAGING, noneSkippedResult,
			devops.PRODUCTION, noneSkippedResult,
		),
		dal.From("cicd_pipelines p"),
		dal.Where(`
			p.result IN ? AND (
				p.type = ? OR EXISTS(SELECT 1 FROM cicd_tasks t WHERE t.pipeline_id = p.id AND t.type = ? AND t.result IN ?)
			)`,
			noneSkippedResult,
			devops.DEPLOYMENT,
			devops.DEPLOYMENT,
			noneSkippedResult,
		),
	}
	if data.Options.ScopeId != nil {
		clauses = append(clauses,
			dal.Where("p.cicd_scope_id = ?", data.Options.ScopeId),
		)
	} else {
		clauses = append(clauses,
			dal.Join("LEFT JOIN project_mapping pm ON (pm.table = 'cicd_scopes' AND pm.row_id = p.cicd_scope_id)"),
			dal.Where("pm.project_name = ?", data.Options.ProjectName),
		)
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()

	enricher, err := api.NewDataConverter(api.DataConverterArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: DoraApiParams{
				ProjectName: data.Options.ProjectName,
			},
			Table: devops.CICDPipeline{}.TableName(),
		},
		InputRowType: reflect.TypeOf(pipelineEx{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			pipelineExInfo := inputRow.(*pipelineEx)
			domainDeployment := &devops.CICDDeployment{
				DomainEntity: domainlayer.DomainEntity{
					Id: pipelineExInfo.Id,
				},
				CicdScopeId:    pipelineExInfo.CicdScopeId,
				Name:           pipelineExInfo.Name,
				Result:         pipelineExInfo.Result,
				Status:         pipelineExInfo.Status,
				OriginalStatus: pipelineExInfo.OriginalStatus,
				OriginalResult: pipelineExInfo.OriginalResult,
				Environment:    pipelineExInfo.Environment,
				TaskDatesInfo: devops.TaskDatesInfo{
					CreatedDate:  pipelineExInfo.CreatedDate,
					QueuedDate:   pipelineExInfo.QueuedDate,
					StartedDate:  pipelineExInfo.StartedDate,
					FinishedDate: pipelineExInfo.FinishedDate,
				},
				DurationSec:       &pipelineExInfo.DurationSec,
				QueuedDurationSec: pipelineExInfo.QueuedDurationSec,
			}
			if pipelineExInfo.FinishedDate != nil && pipelineExInfo.DurationSec != 0 {
				s := pipelineExInfo.FinishedDate.Add(-time.Duration(pipelineExInfo.DurationSec) * time.Second)
				domainDeployment.StartedDate = &s
			}
			if pipelineExInfo.Environment == "" {
				if pipelineExInfo.HasProductionTasks {
					domainDeployment.Environment = devops.PRODUCTION
				} else if pipelineExInfo.HasStagingTasks {
					domainDeployment.Environment = devops.STAGING
				} else if pipelineExInfo.HasTestingTasks {
					domainDeployment.Environment = devops.TESTING
				}
			}
			return []interface{}{domainDeployment}, nil
		},
	})
	if err != nil {
		return err
	}

	return enricher.Execute()
}

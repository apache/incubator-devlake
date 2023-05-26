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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/circleci/models"
	"reflect"
)

var ConvertWorkflowsMeta = plugin.SubTaskMeta{
	Name:             "convertWorkflows",
	EntryPoint:       ConvertWorkflows,
	EnabledByDefault: true,
	Description:      "convert circleci workflows",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ConvertWorkflows(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_WORKFLOW_TABLE)
	db := taskCtx.GetDal()
	clauses := []dal.Clause{
		dal.From(&models.CircleciWorkflow{}),
		dal.Where("connection_id = ? AND project_slug = ?", data.Options.ConnectionId, data.Options.ProjectSlug),
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.CircleciWorkflow{}),
		Input:              cursor,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			userTool := inputRow.(*models.CircleciWorkflow)
			pipeline := &devops.CICDPipeline{
				DomainEntity: domainlayer.DomainEntity{
					Id: getPipelineIdGen().Generate(data.Options.ConnectionId, userTool.Id),
				},
				Name:         userTool.Name,
				Status:       userTool.Status,
				DurationSec:  userTool.DurationSec,
				CreatedDate:  userTool.CreatedAt.ToTime(),
				FinishedDate: userTool.StoppedAt.ToNullableTime(),
				Environment:  userTool.Tag,
			}
			switch userTool.Status {
			case "success":
				pipeline.Result = devops.SUCCESS
			case "failed", "error", "failing":
				pipeline.Result = devops.FAILURE
			}
			if p, err := findProjectByProjectSlug(db, data.Options.ProjectSlug); err == nil {
				pipeline.CicdScopeId = getProjectIdGen().Generate(data.Options.ConnectionId, p.Id)
			}
			result := make([]interface{}, 0, 2)
			result = append(result, pipeline)

			if p, err := findPipelineById(db, userTool.PipelineId); err == nil {
				if p.Vcs.Revision != "" {
					result = append(result, &devops.CiCDPipelineCommit{
						PipelineId: pipeline.Id,
						CommitSha:  p.Vcs.Revision,
						Branch:     p.Vcs.Branch,
						RepoId:     p.Vcs.OriginRepositoryUrl,
						RepoUrl:    p.Vcs.OriginRepositoryUrl,
					})
				}
			}
			return result, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}

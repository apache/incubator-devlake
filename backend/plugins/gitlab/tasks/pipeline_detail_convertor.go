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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ConvertDetailPipelineMeta)
}

var ConvertDetailPipelineMeta = plugin.SubTaskMeta{
	Name:             "Convert Detail Pipelines",
	EntryPoint:       ConvertDetailPipelines,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitlab_detail_pipeline into domain layer table pipeline",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	Dependencies:     []*plugin.SubTaskMeta{&ExtractApiPipelineDetailsMeta},
}

func ConvertDetailPipelines(subtaskCtx plugin.SubTaskContext) errors.Error {
	subtaskCommonArgs, data := CreateSubtaskCommonArgs(subtaskCtx, RAW_PIPELINE_DETAILS_TABLE)

	db := subtaskCtx.GetDal()

	pipelineIdGen := didgen.NewDomainIdGenerator(&models.GitlabPipeline{})
	projectIdGen := didgen.NewDomainIdGenerator(&models.GitlabProject{})

	converter, err := api.NewStatefulDataConverter(&api.StatefulDataConverterArgs[models.GitlabPipeline]{
		SubtaskCommonArgs: subtaskCommonArgs,
		Input: func(stateManager *api.SubtaskStateManager) (dal.Rows, errors.Error) {
			clauses := []dal.Clause{
				dal.From(models.GitlabPipeline{}),
				dal.Where("project_id = ? and connection_id = ?", data.Options.ProjectId, data.Options.ConnectionId),
			}
			if stateManager.IsIncremental() {
				since := stateManager.GetSince()
				if since != nil {
					clauses = append(clauses, dal.Where("updated_at >= ? ", since))
				}
			}
			return db.Cursor(clauses...)
		},
		Convert: func(gitlabPipeline *models.GitlabPipeline) ([]interface{}, errors.Error) {
			createdAt := time.Now()
			if gitlabPipeline.GitlabCreatedAt != nil {
				createdAt = *gitlabPipeline.GitlabCreatedAt
			}
			domainPipeline := &devops.CICDPipeline{
				DomainEntity: domainlayer.DomainEntity{
					Id: pipelineIdGen.Generate(data.Options.ConnectionId, gitlabPipeline.GitlabId),
				},
				Name: pipelineIdGen.Generate(data.Options.ConnectionId, gitlabPipeline.GitlabId),
				Result: devops.GetResult(&devops.ResultRule{
					Success: []string{StatusSuccess, StatusCompleted},
					Failure: []string{StatusFailed, StatusCanceled},
					Default: devops.RESULT_DEFAULT,
				}, gitlabPipeline.Status),
				Status: devops.GetStatus(&devops.StatusRule{
					Done:       []string{StatusSuccess, StatusCompleted, StatusFailed, StatusCanceled},
					InProgress: []string{StatusRunning, StatusWaitingForResource, StatusPending, StatusPreparing},
					Default:    devops.STATUS_OTHER,
				}, gitlabPipeline.Status),
				OriginalStatus: gitlabPipeline.Status,
				TaskDatesInfo: devops.TaskDatesInfo{
					CreatedDate:  createdAt,
					StartedDate:  gitlabPipeline.StartedAt,
					FinishedDate: gitlabPipeline.FinishedAt,
				},
				QueuedDurationSec: gitlabPipeline.QueuedDuration,
				CicdScopeId:       projectIdGen.Generate(data.Options.ConnectionId, gitlabPipeline.ProjectId),
				Environment:       gitlabPipeline.Environment,
				Type:              gitlabPipeline.Type,
				DurationSec:       float64(gitlabPipeline.Duration),
				// DisplayTitle:      gitlabPipeline.Ref,
				Url:     gitlabPipeline.WebUrl,
				IsChild: gitlabPipeline.GenerateIsChild(),
			}
			return []interface{}{
				domainPipeline,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}

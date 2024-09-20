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
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	gitlabModels "github.com/apache/incubator-devlake/plugins/gitlab/models"
)

func init() {
	RegisterSubtaskMeta(&ConvertJobMeta)
}

var ConvertJobMeta = plugin.SubTaskMeta{
	Name:             "Convert Job Runs",
	EntryPoint:       ConvertJobs,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitlab_job into domain layer table job",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	Dependencies:     []*plugin.SubTaskMeta{&ConvertPipelineCommitMeta},
}

func ConvertJobs(taskCtx plugin.SubTaskContext) (err errors.Error) {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GitlabTaskData)
	regexEnricher := data.RegexEnricher

	cursor, err := db.Cursor(dal.From(gitlabModels.GitlabJob{}),
		dal.Where("project_id = ? and connection_id = ?", data.Options.ProjectId, data.Options.ConnectionId))
	if err != nil {
		return err
	}
	defer cursor.Close()

	jobIdGen := didgen.NewDomainIdGenerator(&gitlabModels.GitlabJob{})
	projectIdGen := didgen.NewDomainIdGenerator(&gitlabModels.GitlabProject{})
	pipelineIdGen := didgen.NewDomainIdGenerator(&gitlabModels.GitlabPipeline{})
	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType: reflect.TypeOf(gitlabModels.GitlabJob{}),
		Input:        cursor,
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: gitlabModels.GitlabApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_JOB_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			gitlabJob := inputRow.(*gitlabModels.GitlabJob)

			createdAt := time.Now()
			if gitlabJob.GitlabCreatedAt != nil {
				createdAt = *gitlabJob.GitlabCreatedAt
			}
			domainJob := &devops.CICDTask{
				DomainEntity: domainlayer.DomainEntity{
					Id: jobIdGen.Generate(data.Options.ConnectionId, gitlabJob.GitlabId),
				},
				Name:       gitlabJob.Name,
				PipelineId: pipelineIdGen.Generate(data.Options.ConnectionId, gitlabJob.PipelineId),
				Result: devops.GetResult(&devops.ResultRule{
					Success: []string{StatusSuccess, StatusCompleted},
					Failure: []string{StatusCanceled, StatusFailed},
					Default: devops.RESULT_DEFAULT,
				}, gitlabJob.Status),
				Status: devops.GetStatus(&devops.StatusRule{
					Done:       []string{StatusSuccess, StatusCompleted, StatusFailed},
					InProgress: []string{StatusRunning, StatusWaitingForResource, StatusPreparing, StatusPending},
					Default:    devops.STATUS_OTHER,
				}, gitlabJob.Status),
				OriginalStatus:    gitlabJob.Status,
				DurationSec:       gitlabJob.Duration,
				QueuedDurationSec: &gitlabJob.QueuedDuration,
				TaskDatesInfo: devops.TaskDatesInfo{
					CreatedDate:  createdAt,
					StartedDate:  gitlabJob.StartedAt,
					FinishedDate: gitlabJob.FinishedAt,
				},
				CicdScopeId: projectIdGen.Generate(data.Options.ConnectionId, gitlabJob.ProjectId),
			}
			if data.Options.ScopeConfig.DeploymentPattern != nil || data.Options.ScopeConfig.ProductionPattern != nil {
				domainJob.Type = regexEnricher.ReturnNameIfMatched(devops.DEPLOYMENT, gitlabJob.Name)
				domainJob.Environment = regexEnricher.ReturnNameIfOmittedOrMatched(devops.PRODUCTION, gitlabJob.Name)
			}
			return []interface{}{
				domainJob,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}

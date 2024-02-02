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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/spf13/cast"
	"reflect"
	"time"
)

var _ plugin.SubTaskEntryPoint = ConvertDeployment

func init() {
	RegisterSubtaskMeta(&ConvertDeploymentMeta)
}

var ConvertDeploymentMeta = plugin.SubTaskMeta{
	Name:             "ConvertDeployment",
	EntryPoint:       ConvertDeployment,
	EnabledByDefault: true,
	Description:      "Convert gitlab deployment from tool layer to domain layer",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	Dependencies:     []*plugin.SubTaskMeta{&ExtractDeploymentMeta},
}

// ConvertDeployment should be split into two task theoretically
// But in GitLab, all deployments have commits, so there is no need to change it.
func ConvertDeployment(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_DEPLOYMENT)
	db := taskCtx.GetDal()

	repo := &models.GitlabProject{}
	err := db.First(repo, dal.Where("gitlab_id = ? and connection_id = ?", data.Options.ProjectId, data.Options.ConnectionId))
	if err != nil {
		return err
	}

	projectIdGen := didgen.NewDomainIdGenerator(&models.GitlabProject{})

	cursor, err := db.Cursor(
		dal.From(&models.GitlabDeployment{}),
		dal.Where("connection_id = ? AND gitlab_id = ?", data.Options.ConnectionId, data.Options.ProjectId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	idGen := didgen.NewDomainIdGenerator(&models.GitlabDeployment{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.GitlabDeployment{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			gitlabDeployment := inputRow.(*models.GitlabDeployment)

			var duration *float64
			if gitlabDeployment.DeployableDuration != nil {
				deployableDuration := cast.ToFloat64(*gitlabDeployment.DeployableDuration)
				duration = &deployableDuration
			}
			// Use duration field in resp. DO NOT calculate it manually.
			// GitLab Cloud and GitLab Server both have this fields in response.
			//if duration == nil || *duration == 0 {
			//	if gitlabDeployment.DeployableFinishedAt != nil && gitlabDeployment.DeployableStartedAt != nil {
			//		deployableDuration := float64(gitlabDeployment.DeployableFinishedAt.Sub(*gitlabDeployment.DeployableStartedAt).Milliseconds() / 1e3)
			//		duration = &deployableDuration
			//	}
			//}
			createdDate := time.Now()
			if gitlabDeployment.DeployableCreatedAt != nil {
				createdDate = *gitlabDeployment.DeployableCreatedAt
			}
			domainDeployCommit := &devops.CicdDeploymentCommit{
				DomainEntity: domainlayer.NewDomainEntity(idGen.Generate(data.Options.ConnectionId, data.Options.ProjectId, gitlabDeployment.DeploymentId)),
				CicdScopeId:  projectIdGen.Generate(data.Options.ConnectionId, data.Options.ProjectId),
				Name:         fmt.Sprintf("%s:%d", gitlabDeployment.Name, gitlabDeployment.DeploymentId),
				Result: devops.GetResult(&devops.ResultRule{
					Success: []string{StatusSuccess, StatusCompleted},
					Failure: []string{StatusCanceled, StatusFailed},
					Default: devops.RESULT_DEFAULT,
				}, gitlabDeployment.Status),
				Status: devops.GetStatus(&devops.StatusRule{
					Done:       []string{StatusSuccess, StatusCompleted, StatusFailed, StatusCanceled},
					InProgress: []string{StatusRunning},
					Default:    devops.STATUS_OTHER,
				}, gitlabDeployment.Status),
				OriginalStatus:      gitlabDeployment.Status,
				Environment:         gitlabDeployment.Environment,
				OriginalEnvironment: gitlabDeployment.Environment,
				TaskDatesInfo: devops.TaskDatesInfo{
					CreatedDate:  createdDate,
					StartedDate:  gitlabDeployment.DeployableStartedAt,
					FinishedDate: gitlabDeployment.DeployableFinishedAt,
				},
				DurationSec:       duration,
				QueuedDurationSec: gitlabDeployment.QueuedDuration,
				CommitSha:         gitlabDeployment.Sha,
				RefName:           gitlabDeployment.Ref,
				RepoId:            projectIdGen.Generate(data.Options.ConnectionId, data.Options.ProjectId),
				RepoUrl:           repo.WebUrl,
			}
			if data.RegexEnricher != nil {
				if data.RegexEnricher.ReturnNameIfMatched(devops.ENV_NAME_PATTERN, gitlabDeployment.Environment) != "" {
					domainDeployCommit.Environment = devops.PRODUCTION
				}
			}

			domainDeployCommit.CicdDeploymentId = domainDeployCommit.Id
			return []interface{}{
				domainDeployCommit,
				domainDeployCommit.ToDeployment(),
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}

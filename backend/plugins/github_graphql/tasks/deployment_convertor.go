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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	githubModels "github.com/apache/incubator-devlake/plugins/github/models"
)

var _ plugin.SubTaskEntryPoint = ConvertDeployment

var ConvertDeploymentMeta = plugin.SubTaskMeta{
	Name:             "ConvertDeployment",
	EntryPoint:       ConvertDeployment,
	EnabledByDefault: true,
	Description:      "Convert github deployment from tool layer to domain layer",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ConvertDeployment(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_DEPLOYMENT)
	cursor, err := db.Cursor(
		dal.From(&githubModels.GithubDeployment{}),
		dal.Where("connection_id = ? and github_id = ?", data.Options.ConnectionId, data.Options.GithubId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	jobBuildIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubDeployment{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(githubModels.GithubDeployment{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			githubDeployment := inputRow.(*githubModels.GithubDeployment)
			deploymentCommit := &devops.CicdDeploymentCommit{
				DomainEntity: domainlayer.DomainEntity{
					Id: jobBuildIdGen.Generate(githubDeployment.ConnectionId, githubDeployment.Id),
				},
				CicdScopeId: fmt.Sprintf("%d:%d", githubDeployment.ConnectionId, githubDeployment.GithubId),
				Name:        fmt.Sprintf("%s:%d", githubDeployment.RepositoryName, githubDeployment.DatabaseId), // fixme where does the deploy name field exist?
				Result: devops.GetResult(&devops.ResultRule{
					Success: []string{"SUCCESS"},
					Failed:  []string{"ERROR", "FAILURE"},
					Abort:   []string{"QUEUED", "ABANDONED", "DESTROYED", "INACTIVE"},
					Manual:  []string{"WAITING", "PENDING", "ACTIVE", "IN_PROGRESS"},
					Skipped: []string{},
					Default: githubDeployment.LatestStatusState,
				}, githubDeployment.State),
				Status: devops.GetStatus(&devops.StatusRule[string]{
					InProgress: []string{"ACTIVE", "QUEUED", "IN_PROGRESS", "ABANDONED", "DESTROYED", "FAILURE", "INACTIVE"},
					NotStarted: []string{"PENDING"},
					Done:       []string{"SUCCESS"},
					Manual:     []string{"ERROR", "WAITING"},
					Default:    githubDeployment.State,
				}, githubDeployment.State),
				Environment:  githubDeployment.Environment,
				CreatedDate:  githubDeployment.CreatedDate,
				StartedDate:  &githubDeployment.CreatedDate, // fixme there is no such field
				FinishedDate: &githubDeployment.UpdatedDate, // fixme there is no such field
				CommitSha:    githubDeployment.CommitOid,
				RefName:      githubDeployment.RefName,
				RepoId:       githubDeployment.RepositoryID,
				RepoUrl:      githubDeployment.RepositoryUrl,
			}

			durationSec := uint64(githubDeployment.UpdatedDate.Sub(githubDeployment.CreatedDate).Seconds())
			deploymentCommit.DurationSec = &durationSec

			if data.RegexEnricher != nil {
				if data.RegexEnricher.ReturnNameIfMatched(devops.ENV_NAME_PATTERN, githubDeployment.Environment) != "" {
					deploymentCommit.Environment = devops.PRODUCTION
				}
			}

			deploymentCommit.CicdDeploymentId = deploymentCommit.Id
			return []interface{}{
				deploymentCommit,
				deploymentCommit.ToDeployment(),
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return converter.Execute()
}

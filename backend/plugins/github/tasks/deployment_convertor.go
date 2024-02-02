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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/devops"
	"github.com/apache/incubator-devlake/core/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&ConvertDeploymentsMeta)
}

const (
	RAW_DEPLOYMENT_TABLE = "github_deployment"
)

var ConvertDeploymentsMeta = plugin.SubTaskMeta{
	Name:             "ConvertDeployments",
	EntryPoint:       ConvertDeployment,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_deployments into domain layer table deployment",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
	DependencyTables: []string{models.GithubDeployment{}.TableName()},
	ProductTables:    []string{devops.CicdDeploymentCommit{}.TableName(), devops.CICDDeployment{}.TableName()},
}

// ConvertDeployment should be split into two task theoretically
// But in GitHub, all deployments have commits, so there is no need to change it.
func ConvertDeployment(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_DEPLOYMENT_TABLE)
	cursor, err := db.Cursor(
		dal.From(&models.GithubDeployment{}),
		dal.Where("connection_id = ? and github_id = ?", data.Options.ConnectionId, data.Options.GithubId),
	)
	if err != nil {
		return err
	}
	defer cursor.Close()

	deploymentIdGen := didgen.NewDomainIdGenerator(&models.GithubDeployment{})
	deploymentScopeIdGen := didgen.NewDomainIdGenerator(&models.GithubRepo{})

	converter, err := api.NewDataConverter(api.DataConverterArgs{
		InputRowType:       reflect.TypeOf(models.GithubDeployment{}),
		Input:              cursor,
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Convert: func(inputRow interface{}) ([]interface{}, errors.Error) {
			githubDeployment := inputRow.(*models.GithubDeployment)
			deploymentCommit := &devops.CicdDeploymentCommit{
				DomainEntity: domainlayer.DomainEntity{
					Id: deploymentIdGen.Generate(githubDeployment.ConnectionId, githubDeployment.Id),
				},
				CicdScopeId: deploymentScopeIdGen.Generate(githubDeployment.ConnectionId, githubDeployment.GithubId),
				Name:        githubDeployment.CommitOid,
				Result: devops.GetResult(&devops.ResultRule{
					Success: []string{StatusSuccess, StatusInactive, StatusActive},
					Failure: []string{StatusError, StatusFailure},
					Default: devops.RESULT_DEFAULT,
				}, githubDeployment.State),
				Status: devops.GetStatus(&devops.StatusRule{
					Done:       []string{StatusSuccess, StatusError, StatusFailure, StatusInactive, StatusActive},
					InProgress: []string{StatusInProgress, StatusQueued, StatusWaiting, StatusPending},
					Default:    devops.STATUS_OTHER,
				}, githubDeployment.State),
				OriginalStatus:      githubDeployment.State,
				Environment:         githubDeployment.Environment,
				OriginalEnvironment: githubDeployment.Environment,
				TaskDatesInfo: devops.TaskDatesInfo{
					CreatedDate:  githubDeployment.CreatedDate,
					StartedDate:  &githubDeployment.CreatedDate,
					FinishedDate: &githubDeployment.UpdatedDate,
				},
				CommitSha: githubDeployment.CommitOid,
				RefName:   githubDeployment.RefName,
				RepoId:    deploymentScopeIdGen.Generate(githubDeployment.ConnectionId, githubDeployment.GithubId),
				RepoUrl:   githubDeployment.RepositoryUrl,
			}

			durationSec := float64(githubDeployment.UpdatedDate.Sub(githubDeployment.CreatedDate).Milliseconds() / 1e3)
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

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
	"encoding/json"
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	githubModels "github.com/apache/incubator-devlake/plugins/github/models"
	githubTasks "github.com/apache/incubator-devlake/plugins/github/tasks"
)

var _ plugin.SubTaskEntryPoint = ExtractDeployments

var ExtractDeploymentsMeta = plugin.SubTaskMeta{
	Name:             "Extract Deployments",
	EntryPoint:       ExtractDeployments,
	EnabledByDefault: true,
	Description:      "extract raw deployment data into tool layer table github_graphql_deployment",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

func ExtractDeployments(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*githubTasks.GithubTaskData)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: githubTasks.GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			Table: RAW_DEPLOYMENT,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			deployment := &GraphqlQueryDeploymentDeployment{}
			err := errors.Convert(json.Unmarshal(row.Data, deployment))
			if err != nil {
				return nil, err
			}

			var results []interface{}
			githubDeployment, err := convertGithubDeployment(deployment, data.Options.ConnectionId, data.Options.GithubId)
			if err != nil {
				return nil, errors.Convert(err)
			}
			results = append(results, githubDeployment)

			return results, nil
		},
	})
	if err != nil {
		return err
	}
	return extractor.Execute()
}

func convertGithubDeployment(deployment *GraphqlQueryDeploymentDeployment, connectionId uint64, githubId int) (*githubModels.GithubDeployment, errors.Error) {
	ret := &githubModels.GithubDeployment{
		ConnectionId:      connectionId,
		GithubId:          githubId,
		NoPKModel:         common.NewNoPKModel(),
		Id:                deployment.Id,
		DisplayTitle:      strings.Split(deployment.Commit.Message, "\n")[0],
		Url:               deployment.Repository.Url + "/deployments/" + deployment.Environment,
		DatabaseId:        deployment.DatabaseId,
		Payload:           deployment.Payload,
		Description:       deployment.Description,
		CommitOid:         deployment.CommitOid,
		Environment:       deployment.Environment,
		State:             deployment.State,
		RepositoryID:      deployment.Repository.Id,
		RepositoryName:    deployment.Repository.Name,
		RepositoryUrl:     deployment.Repository.Url,
		CreatedDate:       deployment.CreatedAt,
		UpdatedDate:       deployment.UpdatedAt,
		LatestStatusState: deployment.LatestStatus.State,
		LatestUpdatedDate: deployment.LatestStatus.UpdatedAt,
	}
	if deployment.Ref != nil {
		ret.RefName = deployment.Ref.Name
	}
	return ret, nil
}

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
	"strings"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/common"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	githubModels "github.com/apache/incubator-devlake/plugins/github/models"
	githubTasks "github.com/apache/incubator-devlake/plugins/github/tasks"
	"github.com/merico-dev/graphql"
)

var _ plugin.SubTaskEntryPoint = CollectAndExtractDeployment

const (
	RAW_DEPLOYMENT = "github_deployment"
)

var CollectAndExtractDeploymentMeta = plugin.SubTaskMeta{
	Name:             "CollectAndExtractDeployment",
	EntryPoint:       CollectAndExtractDeployment,
	EnabledByDefault: true,
	Description:      "collect and extract github deployments to raw and tool layer",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

type GraphqlQueryDeploymentWrapper struct {
	RateLimit struct {
		Cost int `graphql:"cost"`
	} `graphql:"rateLimit"`
	Repository struct {
		Deployments struct {
			TotalCount  graphql.Int                        `graphql:"totalCount"`
			PageInfo    *helper.GraphqlQueryPageInfo       `graphql:"pageInfo"`
			Deployments []GraphqlQueryDeploymentDeployment `graphql:"nodes"`
		} `graphql:"deployments(first: $pageSize, after: $skipCursor, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type GraphqlQueryDeploymentDeployment struct {
	Task        string `graphql:"task"` // is value always "deploy"? not sure.
	Id          string `graphql:"id"`
	CommitOid   string `graphql:"commitOid"`
	Environment string `graphql:"environment"`
	State       string `graphql:"state"`
	DatabaseId  uint   `graphql:"databaseId"`
	Description string `graphql:"description"`
	Payload     string `graphql:"payload"`
	Ref         *struct {
		ID     string `graphql:"id"`
		Name   string `graphql:"name"`
		Prefix string `graphql:"prefix"`
	} `graphql:"ref"`
	LatestStatus struct {
		Id        string    `graphql:"id"`
		State     string    `graphql:"state"`
		UpdatedAt time.Time `json:"updatedAt"`
	} `graphql:"latestStatus"`
	Repository struct {
		Id   string `graphql:"id"`
		Name string `graphql:"name"`
		Url  string `graphql:"url"`
	} `graphql:"repository"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CollectAndExtractDeployment will request github api via graphql and store the result into raw layer by default
// ResponseParser's return will be stored to tool layer. So it's called CollectorAndExtractor.
func CollectAndExtractDeployment(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*githubTasks.GithubTaskData)
	collectorWithState, err := helper.NewStatefulApiCollector(helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: githubTasks.GithubApiParams{
			ConnectionId: data.Options.ConnectionId,
			Name:         data.Options.Name,
		},
		Table: RAW_DEPLOYMENT,
	})
	if err != nil {
		return err
	}
	since := helper.DateTime{}
	if collectorWithState.Since != nil {
		since = helper.DateTime{Time: *collectorWithState.Since}
	}
	isSince := false

	err = collectorWithState.InitGraphQLCollector(helper.GraphqlCollectorArgs{
		GraphqlClient: data.GraphqlClient,
		PageSize:      100,
		BuildQuery: func(reqData *helper.GraphqlRequestData) (interface{}, map[string]interface{}, error) {
			if isSince {
				return nil, nil, nil
			}
			query := &GraphqlQueryDeploymentWrapper{}
			variables := make(map[string]interface{})
			if reqData == nil {
				return query, variables, nil
			}
			ownerName := strings.Split(data.Options.Name, "/")
			variables = map[string]interface{}{
				"pageSize":   graphql.Int(reqData.Pager.Size),
				"skipCursor": (*graphql.String)(reqData.Pager.SkipCursor),
				"owner":      graphql.String(ownerName[0]),
				"name":       graphql.String(ownerName[1]),
			}
			return query, variables, nil
		},
		GetPageInfo: func(iQuery interface{}, args *helper.GraphqlCollectorArgs) (*helper.GraphqlQueryPageInfo, error) {
			query := iQuery.(*GraphqlQueryDeploymentWrapper)
			return query.Repository.Deployments.PageInfo, nil
		},
		ResponseParser: func(iQuery interface{}, variables map[string]interface{}) ([]interface{}, error) {
			query := iQuery.(*GraphqlQueryDeploymentWrapper)
			deployments := query.Repository.Deployments.Deployments
			var results []interface{}
			for _, deployment := range deployments {
				//Skip deployments with createdAt earlier than 'since'
				if deployment.CreatedAt.Before(since.Time) {
					isSince = true
					continue
				}
				githubDeployment, err := convertGithubDeployment(deployment, data.Options.ConnectionId, data.Options.GithubId)
				if err != nil {
					return nil, err
				}
				results = append(results, githubDeployment)
			}

			return results, nil
		},
	})
	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}

func convertGithubDeployment(deployment GraphqlQueryDeploymentDeployment, connectionId uint64, githubId int) (*githubModels.GithubDeployment, error) {
	ret := &githubModels.GithubDeployment{
		ConnectionId:      connectionId,
		GithubId:          githubId,
		NoPKModel:         common.NewNoPKModel(),
		Id:                deployment.Id,
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

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
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	githubTasks "github.com/apache/incubator-devlake/plugins/github/tasks"
	"github.com/merico-dev/graphql"
)

var _ plugin.SubTaskEntryPoint = CollectDeployments

const (
	RAW_DEPLOYMENT = "github_graphql_deployment"
)

var CollectDeploymentsMeta = plugin.SubTaskMeta{
	Name:             "Collect Deployments",
	EntryPoint:       CollectDeployments,
	EnabledByDefault: true,
	Description:      "collect github deployments to raw and tool layer from GithubGraphql api",
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
		Id        string     `graphql:"id"`
		State     string     `graphql:"state"`
		UpdatedAt *time.Time `json:"updatedAt"`
	} `graphql:"latestStatus"`
	Repository struct {
		Id   string `graphql:"id"`
		Name string `graphql:"name"`
		Url  string `graphql:"url"`
	} `graphql:"repository"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Commit    struct {
		Oid     string `graphql:"oid"`
		Message string `graphql:"message"`
		Author  struct {
			Name  string `graphql:"name"`
			Email string `graphql:"email"`
		} `graphql:"author"`
		CommittedDate time.Time `graphql:"committedDate"`
	} `graphql:"commit"`
}

// CollectDeployments will request github api via graphql and store the result into raw layer.
func CollectDeployments(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*githubTasks.GithubTaskData)
	apiCollector, err := helper.NewStatefulApiCollector(helper.RawDataSubTaskArgs{
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

	err = apiCollector.InitGraphQLCollector(helper.GraphqlCollectorArgs{
		GraphqlClient: data.GraphqlClient,
		PageSize:      100,
		BuildQuery: func(reqData *helper.GraphqlRequestData) (interface{}, map[string]interface{}, error) {
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
		ResponseParser: func(queryWrapper any) (messages []json.RawMessage, err errors.Error) {
			query := queryWrapper.(*GraphqlQueryDeploymentWrapper)
			deployments := query.Repository.Deployments.Deployments
			for _, rawL := range deployments {
				if apiCollector.GetSince() != nil && !apiCollector.GetSince().Before(rawL.UpdatedAt) {
					return messages, helper.ErrFinishCollect
				}
				messages = append(messages, errors.Must1(json.Marshal(rawL)))
			}
			return
		},
	})
	if err != nil {
		return err
	}
	return apiCollector.Execute()
}

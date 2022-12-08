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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
	githubTasks "github.com/apache/incubator-devlake/plugins/github/tasks"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/merico-dev/graphql"
	"time"
)

const RAW_REPO_TABLE = "github_graphql_repo"

var _ core.SubTaskEntryPoint = CollectRepo

type GraphqlQueryRepo struct {
	RateLimit struct {
		Cost int
	}
	Repository struct {
		Name      string `graphql:"name"`
		GithubId  int    `graphql:"databaseId"`
		HTMLUrl   string `graphql:"url"`
		Languages struct {
			Nodes []struct {
				Name string
			}
		} `graphql:"languages(first: 1)"`
		Description string `graphql:"description"`
		Owner       GraphqlInlineAccountQuery
		CreatedDate time.Time  `graphql:"createdAt"`
		UpdatedDate *time.Time `graphql:"updatedAt"`
		Parent      *struct {
			GithubId int    `graphql:"databaseId"`
			HTMLUrl  string `graphql:"url"`
		}
	} `graphql:"repository(owner: $owner, name: $name)"`
}

var CollectRepoMeta = core.SubTaskMeta{
	Name:             "CollectRepo",
	EntryPoint:       CollectRepo,
	EnabledByDefault: false,
	Description:      "Collect Repo data from GithubGraphql api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE, core.DOMAIN_TYPE_TICKET, core.DOMAIN_TYPE_CICD, core.DOMAIN_TYPE_CODE_REVIEW, core.DOMAIN_TYPE_CROSS},
}

func CollectRepo(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*githubTasks.GithubTaskData)
	db := taskCtx.GetDal()

	collector, err := helper.NewGraphqlCollector(helper.GraphqlCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: githubTasks.GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_REPO_TABLE,
		},
		GraphqlClient: data.GraphqlClient,
		/*
			(Optional) Return query string for request, or you can plug them into UrlTemplate directly
		*/
		BuildQuery: func(reqData *helper.GraphqlRequestData) (interface{}, map[string]interface{}, error) {
			query := &GraphqlQueryRepo{}
			variables := map[string]interface{}{
				"owner": graphql.String(data.Options.Owner),
				"name":  graphql.String(data.Options.Repo),
			}
			return query, variables, nil
		},
		ResponseParser: func(iQuery interface{}, variables map[string]interface{}) ([]interface{}, error) {
			query := iQuery.(*GraphqlQueryRepo)
			repository := query.Repository
			results := make([]interface{}, 0, 1)
			language := ``
			if len(repository.Languages.Nodes) > 0 {
				language = repository.Languages.Nodes[0].Name
			}
			githubRepository := &models.GithubRepo{
				ConnectionId: data.Options.ConnectionId,
				GithubId:     repository.GithubId,
				Name:         repository.Name,
				HTMLUrl:      repository.HTMLUrl,
				Description:  repository.Description,
				OwnerId:      repository.Owner.Id,
				OwnerLogin:   repository.Owner.Login,
				Language:     language,
				CreatedDate:  repository.CreatedDate,
				UpdatedDate:  repository.UpdatedDate,
			}
			data.Repo = githubRepository

			if repository.Parent != nil {
				githubRepository.ParentGithubId = repository.Parent.GithubId
				githubRepository.ParentHTMLUrl = repository.Parent.HTMLUrl
			}
			err := db.CreateOrUpdate(githubRepository)
			if err != nil {
				return nil, err
			}
			results = append(results, githubRepository)

			githubUser, err := convertGraphqlPreAccount(repository.Owner, data.Repo.GithubId, data.Options.ConnectionId)
			if err != nil {
				return nil, err
			}
			results = append(results, githubUser)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

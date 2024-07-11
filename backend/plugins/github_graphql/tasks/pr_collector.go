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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/tasks"
	"github.com/merico-dev/graphql"
)

const RAW_PRS_TABLE = "github_graphql_prs"

type GraphqlQueryPrWrapper struct {
	RateLimit struct {
		Cost int
	}
	// now it orderBy UPDATED_AT and use cursor pagination
	// It may miss some PRs updated when collection.
	// Because these missed PRs will be collected on next, But it's not enough.
	// So Next Millstone(0.17) we should change it to filter by CREATE_AT + collect detail
	Repository struct {
		PullRequests struct {
			PageInfo   *api.GraphqlQueryPageInfo
			Prs        []GraphqlQueryPr `graphql:"nodes"`
			TotalCount graphql.Int
		} `graphql:"pullRequests(first: $pageSize, after: $skipCursor, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type GraphqlQueryPr struct {
	DatabaseId int
	Number     int
	State      string
	Title      string
	Body       string
	Url        string
	Labels     struct {
		Nodes []struct {
			Id   string
			Name string
		}
	} `graphql:"labels(first: 100)"`
	Author    *GraphqlInlineAccountQuery
	Assignees struct {
		// FIXME now domain layer just support one assignee
		Assignees []GraphqlInlineAccountQuery `graphql:"nodes"`
	} `graphql:"assignees(first: 1)"`
	ClosedAt    *time.Time
	MergedAt    *time.Time
	UpdatedAt   time.Time
	CreatedAt   time.Time
	MergeCommit *struct {
		Oid string
	}
	HeadRefName string
	HeadRefOid  string
	BaseRefName string
	BaseRefOid  string
	Commits     struct {
		PageInfo   *api.GraphqlQueryPageInfo
		Nodes      []GraphqlQueryCommit `graphql:"nodes"`
		TotalCount graphql.Int
	} `graphql:"commits(first: 100)"`
	Reviews struct {
		TotalCount graphql.Int
		Nodes      []GraphqlQueryReview `graphql:"nodes"`
	} `graphql:"reviews(first: 100)"`
	MergedBy       *GraphqlInlineAccountQuery
	ReviewRequests struct {
		Nodes []ReviewRequestNode `graphql:"nodes"`
	} `graphql:"reviewRequests(first: 10)"`
}

type ReviewRequestNode struct {
	RequestedReviewer RequestedReviewer `graphql:"requestedReviewer"`
}

type RequestedReviewer struct {
	User User `graphql:"... on User"`
	Team Team `graphql:"... on Team"`
}

type User struct {
	Id    int    `graphql:"databaseId"`
	Login string `graphql:"login"`
	Name  string `graphql:"name"`
}

type Team struct {
	Id   int    `graphql:"databaseId"`
	Name string `graphql:"name"`
	Slug string `graphql:"slug"`
}

type GraphqlQueryReview struct {
	Body       string
	Author     *GraphqlInlineAccountQuery
	State      string `json:"state"`
	DatabaseId int    `json:"databaseId"`
	Commit     struct {
		Oid string
	}
	SubmittedAt *time.Time `json:"submittedAt"`
}

type GraphqlQueryCommit struct {
	Commit struct {
		Oid     string
		Message string
		Author  struct {
			Name  string
			Email string
			Date  time.Time
			User  *GraphqlInlineAccountQuery
		}
		Committer struct {
			Date  time.Time
			Email string
			Name  string
		}
	}
	Url string
}

var CollectPrsMeta = plugin.SubTaskMeta{
	Name:             "Collect Pull Requests",
	EntryPoint:       CollectPrs,
	EnabledByDefault: true,
	Description:      "Collect Pr data from GithubGraphql api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_REVIEW},
}

var _ plugin.SubTaskEntryPoint = CollectPrs

func CollectPrs(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*tasks.GithubTaskData)
	var err errors.Error
	apiCollector, err := api.NewStatefulApiCollector(api.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: tasks.GithubApiParams{
			ConnectionId: data.Options.ConnectionId,
			Name:         data.Options.Name,
		},
		Table: RAW_PRS_TABLE,
	})
	if err != nil {
		return err
	}

	err = apiCollector.InitGraphQLCollector(api.GraphqlCollectorArgs{
		GraphqlClient: data.GraphqlClient,
		PageSize:      10,
		/*
			(Optional) Return query string for request, or you can plug them into UrlTemplate directly
		*/
		BuildQuery: func(reqData *api.GraphqlRequestData) (interface{}, map[string]interface{}, error) {
			query := &GraphqlQueryPrWrapper{}
			if reqData == nil {
				return query, map[string]interface{}{}, nil
			}
			ownerName := strings.Split(data.Options.Name, "/")
			variables := map[string]interface{}{
				"pageSize":   graphql.Int(reqData.Pager.Size),
				"skipCursor": (*graphql.String)(reqData.Pager.SkipCursor),
				"owner":      graphql.String(ownerName[0]),
				"name":       graphql.String(ownerName[1]),
			}
			return query, variables, nil
		},
		GetPageInfo: func(iQuery interface{}, args *api.GraphqlCollectorArgs) (*api.GraphqlQueryPageInfo, error) {
			query := iQuery.(*GraphqlQueryPrWrapper)
			return query.Repository.PullRequests.PageInfo, nil
		},
		ResponseParser: func(iQuery interface{}, variables map[string]interface{}) ([]interface{}, error) {
			query := iQuery.(*GraphqlQueryPrWrapper)
			prs := query.Repository.PullRequests.Prs
			for _, rawL := range prs {
				if apiCollector.GetSince() != nil && !apiCollector.GetSince().Before(rawL.CreatedAt) {
					return nil, api.ErrFinishCollect
				}
			}
			return nil, nil
		},
	})
	if err != nil {
		return err
	}

	return apiCollector.Execute()
}

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
	"github.com/apache/incubator-devlake/core/utils"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	githubTasks "github.com/apache/incubator-devlake/plugins/github/tasks"
	"github.com/merico-dev/graphql"
)

const RAW_RELEASE_TABLE = "github_graphql_release"

type GraphqlQueryReleaseWrapper struct {
	RateLimit struct {
		Cost int
	}
	Repository struct {
		Releases struct {
			TotalCount graphql.Int                  `graphql:"totalCount"`
			PageInfo   *helper.GraphqlQueryPageInfo `graphql:"pageInfo"`
			Releases   []*GraphqlQueryRelease       `graphql:"nodes"`
		} `graphql:"releases(first: $pageSize, after: $skipCursor, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type GraphqlQueryReleaseAuthor struct {
	Name *string `graphql:"name"`
	ID   string  `graphql:"id"`
}

type GraphqlQueryReleaseTagCommit struct {
	ID             string `graphql:"id"`
	Oid            string `graphql:"oid"`
	AbbreviatedOid string `graphql:"abbreviatedOid"`
}

type GraphqlQueryRelease struct {
	Author          GraphqlQueryReleaseAuthor    `graphql:"author"`
	DatabaseID      int                          `graphql:"databaseId"`
	Id              string                       `graphql:"id"`
	CreatedAt       time.Time                    `graphql:"createdAt"`
	Description     string                       `graphql:"description"`
	DescriptionHTML string                       `graphql:"descriptionHTML"`
	IsDraft         bool                         `graphql:"isDraft"`
	IsLatest        bool                         `graphql:"isLatest"`
	IsPrerelease    bool                         `graphql:"isPrerelease"`
	Name            string                       `graphql:"name"`
	PublishedAt     *time.Time                   `graphql:"publishedAt"`
	ResourcePath    string                       `graphql:"resourcePath"`
	TagName         string                       `graphql:"tagName"`
	TagCommit       GraphqlQueryReleaseTagCommit `graphql:"tagCommit"`
	UpdatedAt       time.Time                    `graphql:"updatedAt"`
	URL             string                       `graphql:"url"`
}

var CollectReleaseMeta = plugin.SubTaskMeta{
	Name:             "Collect Releases",
	EntryPoint:       CollectRelease,
	EnabledByDefault: true,
	Description:      "Collect Release data from GithubGraphql api, does not support either timeFilter or diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CICD},
}

var _ plugin.SubTaskEntryPoint = CollectRelease

func CollectRelease(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*githubTasks.GithubTaskData)
	apiCollector, err := helper.NewStatefulApiCollector(helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: githubTasks.GithubApiParams{
			ConnectionId: data.Options.ConnectionId,
			Name:         data.Options.Name,
		},
		Table: RAW_RELEASE_TABLE,
	})
	if err != nil {
		return err
	}

	since := apiCollector.GetSince()
	err = apiCollector.InitGraphQLCollector(helper.GraphqlCollectorArgs{
		GraphqlClient: data.GraphqlClient,
		PageSize:      100,
		BuildQuery: func(reqData *helper.GraphqlRequestData) (interface{}, map[string]interface{}, error) {
			query := &GraphqlQueryReleaseWrapper{}
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
			query := iQuery.(*GraphqlQueryReleaseWrapper)
			return query.Repository.Releases.PageInfo, nil
		},
		ResponseParser: func(queryWrapper any) (messages []json.RawMessage, err errors.Error) {
			query := queryWrapper.(*GraphqlQueryReleaseWrapper)
			releases := query.Repository.Releases.Releases
			for _, rawL := range releases {
				rawL.PublishedAt = utils.NilIfZeroTime(rawL.PublishedAt)
				if since != nil && since.After(rawL.UpdatedAt) {
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

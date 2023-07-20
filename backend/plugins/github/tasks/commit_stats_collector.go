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
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

func init() {
	RegisterSubtaskMeta(&CollectApiCommitStatsMeta)
}

const RAW_COMMIT_STATS_TABLE = "github_api_commit_stats"

var CollectApiCommitStatsMeta = plugin.SubTaskMeta{
	Name:             "collectApiCommitStats",
	EntryPoint:       CollectApiCommitStats,
	EnabledByDefault: false,
	Description:      "Collect commitStats data from Github api, does not support either timeFilter or diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE},
	DependencyTables: []string{
		//models.GithubRepoCommit{}.TableName(), // cursor, config will not regard as dependency
		//models.GithubCommit{}.TableName()}, // cursor
	},
	ProductTables: []string{RAW_COMMIT_STATS_TABLE},
}

func CollectApiCommitStats(taskCtx plugin.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)

	var latestUpdated models.GithubCommitStat
	err := db.First(
		&latestUpdated,
		dal.Join("left join _tool_github_repo_commits on _tool_github_commit_stats.sha = _tool_github_repo_commits.commit_sha"),
		dal.Where("_tool_github_repo_commits.repo_id = ? and _tool_github_repo_commits.connection_id = ?", data.Options.GithubId, data.Options.ConnectionId),
		dal.Orderby("committed_date DESC"),
		dal.Limit(1),
	)
	if err != nil && !db.IsErrorNotFound(err) {
		return errors.Default.Wrap(err, "failed to get latest github commit record")
	}

	cursor, err := db.Cursor(
		dal.Join("left join _tool_github_repo_commits on _tool_github_commits.sha = _tool_github_repo_commits.commit_sha"),
		dal.From(models.GithubCommit{}.TableName()),
		dal.Where("_tool_github_repo_commits.repo_id = ? and _tool_github_repo_commits.connection_id = ? and _tool_github_commits.committed_date >= ?",
			data.Options.GithubId, data.Options.ConnectionId, latestUpdated.CommittedDate.String()),
	)
	if err != nil {
		return err
	}
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(models.GithubCommit{}))
	if err != nil {
		return err
	}

	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraCommits by Board
			*/
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Name:         data.Options.Name,
			},
			/*
				Table store raw data
			*/
			Table: RAW_COMMIT_STATS_TABLE,
		},
		ApiClient: data.ApiClient,
		PageSize:  100,
		Input:     iterator,
		/*
			url may use arbitrary variables from different source in any order, we need GoTemplate to allow more
			flexible for all kinds of possibility.
			Pager contains information for a particular page, calculated by ApiCollector, and will be passed into
			GoTemplate to generate a url for that page.
			We want to do page-fetching in ApiCollector, because the logic are highly similar, by doing so, we can
			avoid duplicate logic for every tasks, and when we have a better idea like improving performance, we can
			do it in one place
		*/
		UrlTemplate: "repos/{{ .Params.Name }}/commits/{{ .Input.Sha }}",
		/*
			(Optional) Return query string for request, or you can plug them into UrlTemplate directly
		*/
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("state", "all")
			query.Set("direction", "asc")
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))

			return query, nil
		},

		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			body, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				return nil, errors.Convert(err)
			}
			return []json.RawMessage{body}, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

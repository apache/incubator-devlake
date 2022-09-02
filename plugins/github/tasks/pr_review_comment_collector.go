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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
	"net/http"
	"net/url"
)

const RAW_PR_REVIEW_COMMENTS_TABLE = "github_api_pull_request_review_comments"

// this struct should be moved to `github_api_common.go`

func CollectPrReviewComments(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)

	since := data.Since
	incremental := false
	if since == nil {
		var latestUpdatedPrReviewComt models.GithubPrComment
		err := db.All(
			&latestUpdatedPrReviewComt,
			dal.Join(`left join _tool_github_pull_requests on 
				_tool_github_pull_requests.github_id = _tool_github_pull_request_comments.pull_request_id 
				and _tool_github_pull_requests.connection_id = _tool_github_pull_request_comments.connection_id`),
			dal.Where(
				"_tool_github_pull_requests.repo_id = ? AND _tool_github_pull_requests.connection_id = ? AND _tool_github_pull_request_comments.type = ?",
				data.Repo.GithubId, data.Repo.ConnectionId, "DIFF",
			),
			dal.Orderby("github_updated_at DESC"),
			dal.Limit(1),
		)
		if err != nil {
			return errors.Default.Wrap(err, "failed to get latest github issue record")
		}
		if latestUpdatedPrReviewComt.GithubId > 0 {
			since = &latestUpdatedPrReviewComt.GithubUpdatedAt
			incremental = true
		}
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_PR_REVIEW_COMMENTS_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: incremental,

		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/pulls/comments",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			if since != nil {
				query.Set("since", since.String())
			}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("direction", "asc")
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var items []json.RawMessage
			err := helper.UnmarshalResponse(res, &items)
			if err != nil {
				return nil, err
			}
			return items, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

var CollectApiPrReviewCommentsMeta = core.SubTaskMeta{
	Name:             "CollectApiPrReviewCommentsMeta",
	EntryPoint:       CollectPrReviewComments,
	EnabledByDefault: true,
	Description:      "Collect pr review comments data from Github api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CODE},
}

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
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"net/http"
	"net/url"
	"reflect"

	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/plugins/core"
)

const RAW_PR_REVIEW_TABLE = "github_api_pull_request_reviews"

// this struct should be moved to `gitub_api_common.go`

var CollectApiPullRequestReviewsMeta = core.SubTaskMeta{
	Name:             "collectApiPullRequestReviews",
	EntryPoint:       CollectApiPullRequestReviews,
	EnabledByDefault: true,
	Description:      "Collect PullRequestReviews data from Github api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

func CollectApiPullRequestReviews(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)
	since := data.Since

	incremental := false
	if since == nil {
		var latestUpdatedPrReview models.GithubPrReview
		err := db.All(
			&latestUpdatedPrReview,
			dal.Join(`left join _tool_github_pull_requests on 
				_tool_github_pull_requests.github_id = _tool_github_pull_request_reviews.pull_request_id 
				and _tool_github_pull_requests.connection_id = _tool_github_pull_request_reviews.connection_id`),
			dal.Where(
				"_tool_github_pull_requests.repo_id = ? AND _tool_github_pull_requests.connection_id = ?", data.Repo.GithubId, data.Repo.ConnectionId,
			),
			dal.Orderby("github_updated_at DESC"),
			dal.Limit(1),
		)
		if err != nil {
			return fmt.Errorf("failed to get latest github issue record: %w", err)
		}
		if latestUpdatedPrReview.GithubId > 0 {
			since = latestUpdatedPrReview.GithubSubmitAt
			incremental = true
		}
	}
	clauses := []dal.Clause{
		dal.Select("number, github_id"),
		dal.From(models.GithubPullRequest{}.TableName()),
		dal.Where("repo_id = ? and connection_id=?", data.Repo.GithubId, data.Options.ConnectionId),
	}
	if since != nil {
		clauses = append(clauses, dal.Where("github_updated_at > ?", *since))
	}
	cursor, err := db.Cursor(
		clauses...,
	)
	if err != nil {
		return err
	}

	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimplePr{}))
	if err != nil {
		return err
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			/*
				This struct will be JSONEncoded and stored into database along with raw data itself, to identity minimal
				set of data to be process, for example, we process JiraIssues by Board
			*/
			Params: GithubApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},

			/*
				Table store raw data
			*/
			Table: RAW_PR_REVIEW_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: incremental,
		Input:       iterator,

		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/pulls/{{ .Input.Number }}/reviews",

		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("direction", "asc")
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))

			return query, nil
		},
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

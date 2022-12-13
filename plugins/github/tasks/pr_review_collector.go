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
	"net/http"
	"net/url"
	"reflect"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/github/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_PR_REVIEW_TABLE = "github_api_pull_request_reviews"

// this struct should be moved to `gitub_api_common.go`

var CollectApiPullRequestReviewsMeta = core.SubTaskMeta{
	Name:             "collectApiPullRequestReviews",
	EntryPoint:       CollectApiPullRequestReviews,
	EnabledByDefault: true,
	Description:      "Collect PullRequestReviews data from Github api",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS, core.DOMAIN_TYPE_CODE_REVIEW},
}

func CollectApiPullRequestReviews(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)

	collectorWithState, err := helper.NewApiCollectorWithState(helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: GithubApiParams{
			ConnectionId: data.Options.ConnectionId,
			Owner:        data.Options.Owner,
			Repo:         data.Options.Repo,
		},
		Table: RAW_PR_REVIEW_TABLE,
	}, data.CreatedDateAfter)
	if err != nil {
		return err
	}

	incremental := collectorWithState.CanIncrementCollect()
	clauses := []dal.Clause{
		dal.Select("number, github_id"),
		dal.From(models.GithubPullRequest{}.TableName()),
		dal.Where("repo_id = ? and connection_id=?", data.Options.GithubId, data.Options.ConnectionId),
	}
	if collectorWithState.CreatedDateAfter != nil {
		clauses = append(clauses, dal.Where("github_created_at > ?", *collectorWithState.CreatedDateAfter))
	}
	if incremental {
		clauses = append(clauses, dal.Where("github_updated_at > ?", *collectorWithState.LatestState.LatestSuccessStart))
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

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: incremental,
		Input:       iterator,

		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/pulls/{{ .Input.Number }}/reviews",

		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("direction", "asc")
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))

			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
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
	return collectorWithState.Execute()
}

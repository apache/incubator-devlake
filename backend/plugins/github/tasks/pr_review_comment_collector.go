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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func init() {
	RegisterSubtaskMeta(&CollectApiPrReviewCommentsMeta)
}

const RAW_PR_REVIEW_COMMENTS_TABLE = "github_api_pull_request_review_comments"

// this struct should be moved to `github_api_common.go`

var CollectApiPrReviewCommentsMeta = plugin.SubTaskMeta{
	Name:             "collectApiPrReviewCommentsMeta",
	EntryPoint:       CollectPrReviewComments,
	EnabledByDefault: true,
	Description:      "Collect pr review comments data from Github api, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS, plugin.DOMAIN_TYPE_CODE_REVIEW},
	DependencyTables: []string{},
	ProductTables:    []string{RAW_PR_REVIEW_COMMENTS_TABLE},
}

func CollectPrReviewComments(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*GithubTaskData)

	collectorWithState, err := helper.NewStatefulApiCollector(helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: GithubApiParams{
			ConnectionId: data.Options.ConnectionId,
			Name:         data.Options.Name,
		},
		Table: RAW_PR_REVIEW_COMMENTS_TABLE,
	}, data.TimeAfter)
	if err != nil {
		return err
	}

	incremental := collectorWithState.IsIncremental()
	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: incremental,

		UrlTemplate: "repos/{{ .Params.Name }}/pulls/comments",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			if data.TimeAfter != nil {
				// Note that `since` is for filtering records by the `updated` time
				// which is not ideal for semantic reasons and would result in slightly more records than expected.
				// But we have no choice since it is the only available field we could exploit from the API.
				query.Set("since", data.TimeAfter.String())
			}
			// if incremental == true, we overwrite it
			if incremental {
				query.Set("since", collectorWithState.LatestState.LatestSuccessStart.String())
			}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("direction", "asc")
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
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

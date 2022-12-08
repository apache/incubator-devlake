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
	"net/http"
	"net/url"
	"time"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/github/models"
)

const RAW_EVENTS_TABLE = "github_api_events"

// this struct should be moved to `gitub_api_common.go`

var CollectApiEventsMeta = core.SubTaskMeta{
	Name:             "collectApiEvents",
	EntryPoint:       CollectApiEvents,
	EnabledByDefault: true,
	Description:      "Collect Events data from Github api",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func CollectApiEvents(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)

	collectorWithState, err := helper.NewApiCollectorWithState(helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: GithubApiParams{
			ConnectionId: data.Options.ConnectionId,
			Owner:        data.Options.Owner,
			Repo:         data.Options.Repo,
		},
		Table: RAW_EVENTS_TABLE,
	}, data.CreatedDateAfter)
	if err != nil {
		return err
	}

	var latestUpdatedTime *time.Time
	incremental := collectorWithState.CanIncrementCollect()
	// user didn't specify a time range to sync, try load from database
	// actually, for github pull, since doesn't make any sense, github pull api doesn't support it
	if incremental {
		var latestUpdatedIssueEvent models.GithubIssueEvent
		err := db.All(
			&latestUpdatedIssueEvent,
			dal.Join("left join _tool_github_issues on _tool_github_issues.github_id = _tool_github_issue_events.issue_id"),
			dal.Where("_tool_github_issues.repo_id = ? and _tool_github_issues.repo_id = ?", data.Repo.GithubId, data.Repo.ConnectionId),
			dal.Orderby("github_created_at DESC"),
			dal.Limit(1),
		)
		if err != nil {
			return errors.Default.Wrap(err, "failed to get latest github issue record")
		}

		if latestUpdatedIssueEvent.GithubId > 0 {
			latestUpdatedTime = &latestUpdatedIssueEvent.GithubCreatedAt
		} else {
			incremental = false
		}
	}

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: incremental,

		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/issues/events",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("state", "all")
			if latestUpdatedTime != nil {
				query.Set("since", latestUpdatedTime.String())
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

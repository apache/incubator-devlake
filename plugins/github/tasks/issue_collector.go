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

const RAW_ISSUE_TABLE = "github_api_issues"

var CollectApiIssuesMeta = core.SubTaskMeta{
	Name:             "collectApiIssues",
	EntryPoint:       CollectApiIssues,
	EnabledByDefault: true,
	Description:      "Collect issues data from Github api",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func CollectApiIssues(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*GithubTaskData)

	collectorWithState, err := helper.NewApiCollectorWithState(helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: GithubApiParams{
			ConnectionId: data.Options.ConnectionId,
			Owner:        data.Options.Owner,
			Repo:         data.Options.Repo,
		},
		Table: RAW_ISSUE_TABLE,
	}, data.CreatedDateAfter)
	if err != nil {
		return err
	}

	var latestUpdatedTime *time.Time
	incremental := collectorWithState.CanIncrementCollect()
	if incremental {
		// try load from database
		var latestUpdated models.GithubIssue
		err = db.All(
			&latestUpdated,
			dal.Where("repo_id = ? and connection_id = ?", data.Repo.GithubId, data.Repo.ConnectionId),
			dal.Orderby("github_updated_at DESC"),
			dal.Limit(1),
		)
		if err != nil {
			return errors.Default.Wrap(err, "failed to get latest github issue record")
		}
		if latestUpdated.GithubId > 0 {
			latestUpdatedTime = &latestUpdated.GithubUpdatedAt
		} else {
			incremental = false
		}
	}

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: incremental,
		/*
			url may use arbitrary variables from different source in any order, we need GoTemplate to allow more
			flexible for all kinds of possibility.
			Pager contains information for a particular page, calculated by ApiCollector, and will be passed into
			GoTemplate to generate a url for that page.
			We want to do page-fetching in ApiCollector, because the logic are highly similar, by doing so, we can
			avoid duplicate logic for every tasks, and when we have a better idea like improving performance, we can
			do it in one place
		*/
		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/issues",
		/*
			(Optional) Return query string for request, or you can plug them into UrlTemplate directly
		*/
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("state", "all")
			// data.CreatedDateAfter need to be used to filter data, but now no params supported
			if latestUpdatedTime != nil {
				query.Set("since", latestUpdatedTime.String())
			}
			query.Set("direction", "asc")
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		/*
			Some api might do pagination by http headers
		*/
		//Header: func(pager *core.Pager) http.Header {
		//},
		/*
			Sometimes, we need to collect data based on previous collected data, like jira changelog, it requires
			issue_id as part of the url.
			We can mimic `stdin` design, to accept a `Input` function which produces a `Iterator`, collector
			should iterate all records, and do data-fetching for each on, either in parallel or sequential order
			UrlTemplate: "api/3/issue/{{ Input.ID }}/changelog"
		*/
		//Input: databaseIssuesIterator,
		/*
			For api endpoint that returns number of total pages, ApiCollector can collect pages in parallel with ease,
			or other techniques are required if this information was missing.
		*/
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

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

	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/plugins/core"
	githubTasks "github.com/apache/incubator-devlake/plugins/github/tasks"
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

func CollectApiEvents(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GithubGraphqlTaskData)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubGraphqlApiParams{
				ConnectionId: data.Options.ConnectionId,
				Owner:        data.Options.Owner,
				Repo:         data.Options.Repo,
			},
			Table: RAW_EVENTS_TABLE,
		},
		ApiClient:   data.HttpClient,
		PageSize:    100,
		Incremental: true,

		UrlTemplate: "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/issues/events",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("state", "all")
			//if since != nil {
			//	query.Set("since", since.String())
			//}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("direction", "asc")
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))

			return query, nil
		},
		GetTotalPages: githubTasks.GetTotalPagesFromResponse,
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

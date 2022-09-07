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
	"time"

	"github.com/apache/incubator-devlake/plugins/core/dal"

	"github.com/apache/incubator-devlake/plugins/helper"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitea/models"
)

const RAW_ISSUE_TABLE = "gitea_api_issues"

var CollectApiIssuesMeta = core.SubTaskMeta{
	Name:             "collectApiIssues",
	EntryPoint:       CollectApiIssues,
	EnabledByDefault: true,
	Description:      "Collect issues data from Gitea api",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func CollectApiIssues(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDal()
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_ISSUE_TABLE)

	since := data.Since
	incremental := false
	// user didn't specify a time range to sync, try load from database
	if since == nil {
		var latestUpdated models.GiteaIssue
		err := db.All(
			&latestUpdated,
			dal.Where("repo_id = ? and connection_id = ?", data.Repo.GiteaId, data.Repo.ConnectionId),
			dal.Orderby("gitea_updated_at DESC"),
			dal.Limit(1),
		)
		if err != nil {
			return fmt.Errorf("failed to get latest gitea issue record: %w", err)
		}
		if latestUpdated.GiteaId > 0 {
			since = &latestUpdated.GiteaUpdatedAt
			incremental = true
		}
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           50,
		Incremental:        incremental,
		UrlTemplate:        "repos/{{ .Params.Owner }}/{{ .Params.Repo }}/issues",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("state", "all")
			if since != nil {
				query.Set("since", since.Format(time.RFC3339))
			}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("limit", fmt.Sprintf("%v", reqData.Pager.Size))

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

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

const RAW_AUTOMATION_RUN_TABLE = "testmo_automation_runs"

var CollectAutomationRunsMeta = plugin.SubTaskMeta{
	Name:             "collectAutomationRuns",
	EntryPoint:       CollectAutomationRuns,
	EnabledByDefault: true,
	Description:      "Collect automation runs data from Testmo api",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CODE_QUALITY},
}

func CollectAutomationRuns(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*TestmoTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collecting automation runs")

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TestmoApiParams{
				ConnectionId: data.Options.ConnectionId,
				ProjectId:    data.Options.ProjectId,
			},
			Table: RAW_AUTOMATION_RUN_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    100,
		Incremental: false,
		UrlTemplate: "projects/{{ .Params.ProjectId }}/automation/runs",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("sort", "automation_runs:created_at")
			query.Set("order", "desc")

			// Add time filter if specified
			if data.Options.CreatedDateAfter != nil {
				query.Set("created_after", data.Options.CreatedDateAfter.Format("2006-01-02T15:04:05Z"))
			}

			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			return GetRawMessageFromResponse(res)
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

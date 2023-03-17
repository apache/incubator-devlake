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
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

const RAW_ISSUE_TYPE_TABLE = "jira_api_issue_types"

var _ plugin.SubTaskEntryPoint = CollectIssueTypes

var CollectIssueTypesMeta = plugin.SubTaskMeta{
	Name:             "collectIssueTypes",
	EntryPoint:       CollectIssueTypes,
	EnabledByDefault: true,
	Description:      "collect Jira issue_types, does not support either timeFilter or diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func CollectIssueTypes(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect issue_types")

	urlTemplate := "api/3/issuetype"
	if data.JiraServerInfo.DeploymentType == models.DeploymentServer {
		urlTemplate = "api/2/issuetype"
	}
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_ISSUE_TYPE_TABLE,
		},
		ApiClient:   data.ApiClient,
		Concurrency: 1,
		UrlTemplate: urlTemplate,

		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data []json.RawMessage
			err := api.UnmarshalResponse(res, &data)
			return data, err
		},
		AfterResponse: ignoreHTTPStatus400,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

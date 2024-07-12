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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"net/http"
)

const RAW_ISSUE_FIELDS_TABLE = "jira_api_issue_fields"

var _ plugin.SubTaskEntryPoint = CollectIssueField

var CollectIssueFieldsMeta = plugin.SubTaskMeta{
	Name:             "collectIssuleField",
	EntryPoint:       CollectIssueField,
	EnabledByDefault: true,
	Description:      "collect Jira issue field, does not support either timeFilter or diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func CollectIssueField(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect issue fields")
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_ISSUE_FIELDS_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    0,
		UrlTemplate: "api/2/field",
		Query:       nil,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data []json.RawMessage
			err := api.UnmarshalResponse(res, &data)
			if err != nil {
				return nil, err
			}
			return data, nil
		},
		AfterResponse: ignoreHTTPStatus400,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

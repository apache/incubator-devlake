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
	"github.com/apache/incubator-devlake/errors"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_PROJECT_TABLE = "jira_api_projects"

var _ core.SubTaskEntryPoint = CollectProjects

var CollectProjectsMeta = core.SubTaskMeta{
	Name:             "collectProjects",
	EntryPoint:       CollectProjects,
	EnabledByDefault: true,
	Description:      "collect Jira projects",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func CollectProjects(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect projects")
	jql := "ORDER BY created ASC"
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Options.ConnectionId,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_PROJECT_TABLE,
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: "api/2/project",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("jql", jql)
			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var result []json.RawMessage
			err := helper.UnmarshalResponse(res, &result)
			return result, err
		},
	})
	if err != nil {
		logger.Error(err, "collect project error")
		return err
	}
	return collector.Execute()
}

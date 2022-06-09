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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"net/http"
	"net/url"
)

const RAW_TASK_CUSTOM_FIELDS_TABLE = "tapd_api_task_custom_fields"

var _ core.SubTaskEntryPoint = CollectTaskCustomFields

func CollectTaskCustomFields(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect task_custom_fields")
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,
				WorkspaceID:  data.Options.WorkspaceID,
			},
			Table: RAW_TASK_CUSTOM_FIELDS_TABLE,
		},
		ApiClient: data.ApiClient,
		//PageSize:    100,
		UrlTemplate: "tasks/custom_fields_settings",
		Query: func(reqData *helper.RequestData, taskCtx core.SubTaskContext) (url.Values, error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceID))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				TaskCustomFields []json.RawMessage `json:"data"`
			}
			err := helper.UnmarshalResponse(res, &data)
			return data.TaskCustomFields, err
		},
	})
	if err != nil {
		logger.Error("collect task_custom_fields error:", err)
		return err
	}
	return collector.Execute()
}

var CollectTaskCustomFieldsMeta = core.SubTaskMeta{
	Name:        "collectTaskCustomFields",
	EntryPoint:  CollectTaskCustomFields,
	Required:    true,
	Description: "collect Tapd TaskCustomFields",
}

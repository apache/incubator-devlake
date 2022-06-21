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

const RAW_WORKSPACE_TABLE = "tapd_api_sub_workspaces"

var _ core.SubTaskEntryPoint = CollectWorkspaces

func CollectWorkspaces(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_WORKSPACE_TABLE)
	logger := taskCtx.GetLogger()
	logger.Info("collect workspaces")
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		//PageSize:    100,
		UrlTemplate: "workspaces/sub_workspaces",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			//query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			//query.Set("limit", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Workspaces []json.RawMessage `json:"data"`
			}
			err := helper.UnmarshalResponse(res, &data)
			return data.Workspaces, err
		},
	})
	if err != nil {
		logger.Error("collect workspace error:", err)
		return err
	}
	return collector.Execute()
}

var CollectWorkspaceMeta = core.SubTaskMeta{
	Name:             "collectWorkspaces",
	EntryPoint:       CollectWorkspaces,
	EnabledByDefault: true,
	Description:      "collect Tapd workspaces",
}

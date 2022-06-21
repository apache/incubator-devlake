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
	"fmt"
	"net/url"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_STORY_STATUS_TABLE = "tapd_api_story_status"

var _ core.SubTaskEntryPoint = CollectStoryStatus

func CollectStoryStatus(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_STORY_STATUS_TABLE)
	logger := taskCtx.GetLogger()
	logger.Info("collect bugStatus")

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		UrlTemplate:        "workflows/status_map",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			query.Set("system", "story")
			return query, nil
		},
		ResponseParser: GetRawMessageDirectFromResponse,
	})
	if err != nil {
		logger.Error("collect bugStatus error:", err)
		return err
	}
	return collector.Execute()
}

var CollectStoryStatusMeta = core.SubTaskMeta{
	Name:             "collectStoryStatus",
	EntryPoint:       CollectStoryStatus,
	EnabledByDefault: true,
	Description:      "collect Tapd bugStatus",
}

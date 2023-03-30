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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

const RAW_TASK_TABLE = "tapd_api_tasks"

var _ plugin.SubTaskEntryPoint = CollectTasks

func CollectTasks(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TASK_TABLE)
	logger := taskCtx.GetLogger()
	logger.Info("collect tasks")
	collectorWithState, err := helper.NewStatefulApiCollector(*rawDataSubTaskArgs, data.TimeAfter)
	if err != nil {
		return err
	}
	incremental := collectorWithState.IsIncremental()

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Incremental:        incremental,
		ApiClient:          data.ApiClient,
		PageSize:           int(data.Options.PageSize),
		UrlTemplate:        "tasks",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("limit", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("fields", "labels")
			query.Set("order", "created asc")
			if data.TimeAfter != nil {
				query.Set("modified",
					fmt.Sprintf(">%s",
						data.TimeAfter.In(data.Options.CstZone).Format("2006-01-02")))
			}
			if incremental {
				query.Set("modified",
					fmt.Sprintf(">%s",
						collectorWithState.LatestState.LatestSuccessStart.In(data.Options.CstZone).Format("2006-01-02")))
			}
			return query, nil
		},
		ResponseParser: GetRawMessageArrayFromResponse,
	})
	if err != nil {
		logger.Error(err, "collect task error")
		return err
	}
	return collector.Execute()
}

var CollectTaskMeta = plugin.SubTaskMeta{
	Name:             "collectTasks",
	EntryPoint:       CollectTasks,
	EnabledByDefault: true,
	Description:      "collect Tapd tasks",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

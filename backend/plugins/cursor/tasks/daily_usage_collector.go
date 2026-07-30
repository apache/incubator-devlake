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
	"net/http"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

func newDailyUsageDateRangeIterator(since *time.Time, isIncremental bool) *helper.QueueIterator {
	startMs, endMs := computeUsageTimeRangeMs(since, time.Now().UTC(), isIncremental)
	iter := helper.NewQueueIterator()
	for _, chunk := range splitDailyUsageTimeRangeMs(startMs, endMs, cursorDailyUsageMaxDays) {
		iter.Push(chunk)
	}
	return iter
}

func dailyUsagePostBody(reqData *helper.RequestData) map[string]interface{} {
	input := reqData.Input.(cursorTimeRangeInput)
	page := 1
	if state, ok := reqData.CustomData.(cursorPageState); ok && state.Page > 0 {
		page = state.Page
	} else if reqData.Pager.Page > 0 {
		page = reqData.Pager.Page
	}
	return map[string]interface{}{
		"startDate": input.StartDateMs,
		"endDate":   input.EndDateMs,
		"page":      page,
		"pageSize":  reqData.Pager.Size,
	}
}

// CollectDailyUsage collects per-user per-day adoption metrics from POST /teams/daily-usage-data.
func CollectDailyUsage(taskCtx plugin.SubTaskContext) errors.Error {
	data, ok := taskCtx.TaskContext().GetData().(*CursorTaskData)
	if !ok {
		return errors.Default.New("task data is not CursorTaskData")
	}
	apiClient, err := CreateApiClient(taskCtx.TaskContext(), data.Connection)
	if err != nil {
		return err
	}

	rawArgs := helper.RawDataSubTaskArgs{
		Ctx:     taskCtx,
		Table:   rawDailyUsageTable,
		Options: rawParamsFromTaskData(data),
	}

	collector, err := helper.NewStatefulApiCollector(rawArgs)
	if err != nil {
		return err
	}

	err = collector.InitCollector(helper.ApiCollectorArgs{
		ApiClient:             apiClient,
		Input:                 newDailyUsageDateRangeIterator(collector.GetSince(), collector.IsIncremental()),
		Method:                http.MethodPost,
		UrlTemplate:           "teams/daily-usage-data",
		PageSize:              cursorApiPageSize,
		RequestBody:           dailyUsagePostBody,
		GetNextPageCustomData: nextPostPage,
		ResponseParser:        parseDailyUsageResponse,
	})
	if err != nil {
		return err
	}

	logUsageCollectionWindow(taskCtx.GetLogger(), "teams/daily-usage-data", collector.GetSince(), collector.IsIncremental())
	return collector.Execute()
}

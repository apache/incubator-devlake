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

func newUsageEventsDateRangeIterator(since *time.Time, isIncremental bool) *helper.QueueIterator {
	startMs, endMs := computeUsageTimeRangeMs(since, time.Now().UTC(), isIncremental)
	iter := helper.NewQueueIterator()
	for _, chunk := range splitDailyUsageTimeRangeMs(startMs, endMs, cursorDailyUsageMaxDays) {
		iter.Push(chunk)
	}
	return iter
}

func nextPostPage(prevReqData *helper.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
	pagination, err := parsePagination(prevPageResponse)
	if err != nil {
		return nil, err
	}
	if !pagination.HasNextPage {
		return nil, helper.ErrFinishCollect
	}
	currentPage := 1
	if state, ok := prevReqData.CustomData.(cursorPageState); ok && state.Page > 0 {
		currentPage = state.Page
	} else if prevReqData.Pager.Page > 0 {
		currentPage = prevReqData.Pager.Page
	}
	return cursorPageState{Page: currentPage + 1}, nil
}

func postPageBody(reqData *helper.RequestData) map[string]interface{} {
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

// CollectUsageEvents collects granular usage events from POST /teams/filtered-usage-events.
func CollectUsageEvents(taskCtx plugin.SubTaskContext) errors.Error {
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
		Table:   rawUsageEventsTable,
		Options: rawParamsFromTaskData(data),
	}

	collector, err := helper.NewStatefulApiCollector(rawArgs)
	if err != nil {
		return err
	}

	err = collector.InitCollector(helper.ApiCollectorArgs{
		ApiClient:             apiClient,
		Input:                 newUsageEventsDateRangeIterator(collector.GetSince(), collector.IsIncremental()),
		Method:                http.MethodPost,
		UrlTemplate:           "teams/filtered-usage-events",
		PageSize:              cursorApiPageSize,
		RequestBody:           postPageBody,
		GetNextPageCustomData: nextPostPage,
		ResponseParser:        parseUsageEventsResponse,
	})
	if err != nil {
		return err
	}

	logUsageCollectionWindow(taskCtx.GetLogger(), "teams/filtered-usage-events", collector.GetSince(), collector.IsIncremental())
	return collector.Execute()
}

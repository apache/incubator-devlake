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
	"net/url"
	"reflect"
	"strconv"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
)

const RAW_TASK_TABLE = "clickup_tasks"

// clickUpTaskListResponse mirrors the envelope returned by GET /list/{id}/task.
// ClickUp uses zero-based page numbers and signals the end of the collection
// with `last_page: true` (and/or an empty `tasks` array).
type clickUpTaskListResponse struct {
	Tasks    []json.RawMessage `json:"tasks"`
	LastPage bool              `json:"last_page"`
}

// listInput is the per-list iterator element driving task collection: tasks are
// collected list-by-list across every list in the scoped folder.
type listInput struct {
	ListId string `gorm:"column:list_id"`
}

var CollectTaskMeta = plugin.SubTaskMeta{
	Name:             "Collect Tasks",
	EntryPoint:       CollectTasks,
	EnabledByDefault: true,
	Description:      "Collect tasks for every list in the scoped ClickUp folder (page-based pagination)",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

var _ plugin.SubTaskEntryPoint = CollectTasks

func CollectTasks(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*ClickUpTaskData)
	db := taskCtx.GetDal()

	// Iterate every list of this folder collected by CollectLists.
	cursor, err := db.Cursor(
		dal.Select("list_id"),
		dal.From(&models.ClickUpList{}),
		dal.Where("connection_id = ? AND folder_id = ?", data.Options.ConnectionId, data.Options.FolderId),
	)
	if err != nil {
		return err
	}
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(listInput{}))
	if err != nil {
		return err
	}

	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ClickUpApiParams{
				ConnectionId: data.Options.ConnectionId,
				FolderId:     data.Options.FolderId,
			},
			Table: RAW_TASK_TABLE,
		},
		ApiClient:   data.ApiClient,
		Input:       iterator,
		PageSize:    100,
		UrlTemplate: "list/{{ .Input.ListId }}/task",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("subtasks", "true")
			query.Set("include_closed", "true")
			page := "0"
			if reqData.CustomData != nil {
				if p, ok := reqData.CustomData.(string); ok && p != "" {
					page = p
				}
			}
			query.Set("page", page)
			// Incremental collection: restrict to tasks updated after the
			// configured cut-off (ClickUp expects milliseconds since epoch).
			if data.TimeAfter != nil {
				query.Set("date_updated_gt", strconv.FormatInt(data.TimeAfter.UnixMilli(), 10))
			}
			return query, nil
		},
		GetNextPageCustomData: func(prevReqData *api.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
			var resp clickUpTaskListResponse
			if err := api.UnmarshalResponse(prevPageResponse, &resp); err != nil {
				return nil, err
			}
			if resp.LastPage || len(resp.Tasks) == 0 {
				return nil, api.ErrFinishCollect
			}
			prevPage := 0
			if prevReqData.CustomData != nil {
				if p, ok := prevReqData.CustomData.(string); ok {
					prevPage, _ = strconv.Atoi(p)
				}
			}
			return strconv.Itoa(prevPage + 1), nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var resp clickUpTaskListResponse
			if err := api.UnmarshalResponse(res, &resp); err != nil {
				return nil, err
			}
			return resp.Tasks, nil
		},
	})
	if err != nil {
		return err
	}
	return collector.Execute()
}

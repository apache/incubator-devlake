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
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

const RAW_WORKLOGS_TABLE = "jira_api_worklogs"

func CollectWorklogs(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*JiraTaskData)
	since := data.Since

	logger := taskCtx.GetLogger()
	connectionId := data.Connection.ID
	boardId := data.Options.BoardId
	tx := db.Table("_tool_jira_board_issues bi").
		Select("bi.issue_id, NOW() AS update_time").
		Joins("LEFT JOIN _tool_jira_issues i ON (bi.connection_id = i.connection_id AND bi.issue_id = i.issue_id)").
		Where("bi.connection_id = ? AND bi.board_id = ? AND (i.worklog_updated IS NULL OR i.worklog_updated < i.updated)", connectionId, boardId)

	if since != nil {
		tx = tx.Where("i.updated > ?", since)
	}
	cursor, err := tx.Rows()
	if err != nil {
		return err
	}
	iterator, err := helper.NewCursorIterator(db, cursor, reflect.TypeOf(apiv2models.Input{}))
	if err != nil {
		return err
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_WORKLOGS_TABLE,
		},
		Input:         iterator,
		ApiClient:     data.ApiClient,
		UrlTemplate:   "api/2/issue/{{ .Input.IssueId }}/worklog",
		PageSize:      50,
		Incremental:   true,
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Worklogs []json.RawMessage `json:"worklogs"`
			}
			err := helper.UnmarshalResponse(res, &data)
			if err != nil {
				return nil, err
			}
			return data.Worklogs, nil
		},
		AfterResponse: ignoreHTTPStatus404,
	})
	if err != nil {
		logger.Error("collect board error:", err)
		return err
	}

	return collector.Execute()
}

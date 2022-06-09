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
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"net/http"
	"net/url"
	"reflect"
	"time"
)

const RAW_TASK_COMMIT_TABLE = "tapd_api_task_commits"

var _ core.SubTaskEntryPoint = CollectTaskCommits

type SimpleTask struct {
	Id uint64
}

func CollectTaskCommits(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	db := taskCtx.GetDb()
	logger := taskCtx.GetLogger()
	logger.Info("collect issueCommits")
	num := 0
	since := data.Since
	incremental := false
	if since == nil {
		// user didn't specify a time range to sync, try load from database
		var latestUpdated models.TapdTaskCommit
		err := db.Where("connection_id = ? and workspace_id = ?", data.Connection.ID, data.Options.WorkspaceID).Order("created DESC").Limit(1).Find(&latestUpdated).Error
		if err != nil {
			return fmt.Errorf("failed to get latest tapd changelog record: %w", err)
		}
		if latestUpdated.ID > 0 {
			since = (*time.Time)(&latestUpdated.Created)
			incremental = true
		}
	}

	tx := db.Model(&models.TapdTask{})
	tx = tx.Where("connection_id = ? and workspace_id = ?", data.Options.ConnectionId, data.Options.WorkspaceID)

	if since != nil {
		tx = tx.Where("modified > ?", since)
	}
	cursor, err := tx.Select("id").Rows()
	if err != nil {
		return err
	}
	iterator, err := helper.NewCursorIterator(db, cursor, reflect.TypeOf(SimpleTask{}))
	if err != nil {
		return err
	}
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_TASK_COMMIT_TABLE,
		},
		Incremental: incremental,
		ApiClient:   data.ApiClient,
		//PageSize:    100,
		Input:       iterator,
		UrlTemplate: "code_commit_infos",
		Query: func(reqData *helper.RequestData, options interface{}) (url.Values, error) {
			input := reqData.Input.(*SimpleTask)
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceID))
			query.Set("type", "task")
			query.Set("object_id", fmt.Sprintf("%v", input.Id))
			query.Set("order", "created asc")
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Stories []json.RawMessage `json:"data"`
			}
			if len(data.Stories) > 0 {
				fmt.Println(len(data.Stories))
				num += len(data.Stories)
				fmt.Printf("num is %d", num)

			}
			err := helper.UnmarshalResponse(res, &data)
			return data.Stories, err
		},
	})
	if err != nil {
		logger.Error("collect issueCommit error:", err)
		return err
	}
	return collector.Execute()
}

var CollectTaskCommitMeta = core.SubTaskMeta{
	Name:        "collectTaskCommits",
	EntryPoint:  CollectTaskCommits,
	Required:    true,
	Description: "collect Tapd issueCommits",
}

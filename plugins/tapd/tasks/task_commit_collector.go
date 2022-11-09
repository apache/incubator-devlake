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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"net/http"
	"net/url"
	"reflect"
)

const RAW_TASK_COMMIT_TABLE = "tapd_api_task_commits"

var _ core.SubTaskEntryPoint = CollectTaskCommits

func CollectTaskCommits(taskCtx core.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TASK_COMMIT_TABLE, false)
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	logger.Info("collect issueCommits")
	since := data.Since
	clauses := []dal.Clause{
		dal.Select("_tool_tapd_tasks.id as issue_id, modified as update_time"),
		dal.From(&models.TapdTask{}),
		dal.Join("LEFT JOIN _tool_tapd_task_commits tjbc ON (tjbc.connection_id = _tool_tapd_tasks.connection_id AND tjbc.task_id = _tool_tapd_tasks.id)"),
		dal.Where("_tool_tapd_tasks.modified > _tool_tapd_tasks.created AND tjbc.connection_id = ? and tjbc.workspace_id = ? ", data.Options.ConnectionId, data.Options.WorkspaceId),
		dal.Groupby("_tool_tapd_tasks.id, _tool_tapd_tasks.modified"),
		dal.Having("_tool_tapd_tasks.modified > max(tjbc.issue_updated) OR  max(tjbc.issue_updated) IS NULL"),
	}
	if since != nil {
		clauses = append(clauses, dal.Where("modified > ?", since))
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(models.Input{}))
	if err != nil {
		return err
	}
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Incremental:        since == nil,
		ApiClient:          data.ApiClient,
		//PageSize:    100,
		Input:       iterator,
		UrlTemplate: "code_commit_infos",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			input := reqData.Input.(*models.Input)
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			query.Set("type", "task")
			query.Set("object_id", fmt.Sprintf("%v", input.IssueId))
			query.Set("order", "created asc")
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				Stories []json.RawMessage `json:"data"`
			}
			err := helper.UnmarshalResponse(res, &data)
			return data.Stories, err
		},
	})
	if err != nil {
		logger.Error(err, "collect issueCommit error")
		return err
	}
	return collector.Execute()
}

var CollectTaskCommitMeta = core.SubTaskMeta{
	Name:             "collectTaskCommits",
	EntryPoint:       CollectTaskCommits,
	EnabledByDefault: true,
	Description:      "collect Tapd issueCommits",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

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
	"github.com/apache/incubator-devlake/errors"
	"net/http"
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

const RAW_WORKLOGS_TABLE = "jira_api_worklogs"

var CollectWorklogsMeta = core.SubTaskMeta{
	Name:             "collectWorklogs",
	EntryPoint:       CollectWorklogs,
	EnabledByDefault: true,
	Description:      "collect Jira work logs",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

func CollectWorklogs(taskCtx core.SubTaskContext) errors.Error {
	db := taskCtx.GetDal()
	data := taskCtx.GetData().(*JiraTaskData)

	logger := taskCtx.GetLogger()

	collectorWithState, err := helper.NewApiCollectorWithState(helper.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: JiraApiParams{
			ConnectionId: data.Options.ConnectionId,
			BoardId:      data.Options.BoardId,
		},
		Table: RAW_WORKLOGS_TABLE,
	}, data.StartFrom)
	if err != nil {
		return err
	}

	// filter out issue_ids that needed collection
	clauses := []dal.Clause{
		dal.Select("i.issue_id, i.updated AS update_time"),
		dal.From("_tool_jira_board_issues bi"),
		dal.Join("LEFT JOIN _tool_jira_issues i ON (bi.connection_id = i.connection_id AND bi.issue_id = i.issue_id)"),
		dal.Join("LEFT JOIN _tool_jira_worklogs wl ON (wl.connection_id = i.connection_id AND wl.issue_id = i.issue_id)"),
		dal.Where("i.updated > i.created AND bi.connection_id = ?  AND bi.board_id = ?  ", data.Options.ConnectionId, data.Options.BoardId),
		dal.Groupby("i.issue_id, i.updated"),
	}
	incremental := collectorWithState.CanIncrementCollect()
	if incremental {
		clauses = append(clauses, dal.Having("i.updated > max(wl.issue_updated) OR  (max(wl.issue_updated) IS NULL AND COUNT(wl.worklog_id) > 0)"))
	}

	// construct the input iterator
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(apiv2models.Input{}))
	if err != nil {
		return err
	}

	err = collectorWithState.InitCollector(helper.ApiCollectorArgs{
		Input:         iterator,
		ApiClient:     data.ApiClient,
		UrlTemplate:   "api/2/issue/{{ .Input.IssueId }}/worklog",
		PageSize:      50,
		Incremental:   incremental,
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
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
		logger.Error(err, "collect board error")
		return err
	}

	return collectorWithState.Execute()
}

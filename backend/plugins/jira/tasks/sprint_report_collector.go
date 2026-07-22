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
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

const RAW_SPRINT_REPORT_TABLE = "jira_api_sprint_reports"

var _ plugin.SubTaskEntryPoint = CollectSprintReport

var CollectSprintReportMeta = plugin.SubTaskMeta{
	Name:             "collectSprintReport",
	EntryPoint:       CollectSprintReport,
	EnabledByDefault: true,
	Description:      "collect Jira Sprint Report, the frozen committed/completed snapshot Jira takes at sprint close",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

// CollectSprintReport fetches the Greenhopper Sprint Report for every closed
// sprint on this board. It depends on CollectSprints/ExtractSprints having
// already populated _tool_jira_board_sprints and _tool_jira_sprints, since
// that's where the (boardId, sprintId) pairs to query come from. Only closed
// sprints are queried: the report is only a meaningful, frozen snapshot once
// a sprint has actually closed.
func CollectSprintReport(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()

	apiCollector, err := api.NewStatefulApiCollector(api.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: JiraApiParams{
			ConnectionId: data.Options.ConnectionId,
			BoardId:      data.Options.BoardId,
		},
		Table: RAW_SPRINT_REPORT_TABLE,
	})
	if err != nil {
		return err
	}

	clauses := []dal.Clause{
		dal.Select("bs.board_id AS board_id, bs.sprint_id AS sprint_id, s.complete_date AS update_time"),
		dal.From("_tool_jira_board_sprints bs"),
		dal.Join("LEFT JOIN _tool_jira_sprints s ON (bs.connection_id = s.connection_id AND bs.sprint_id = s.sprint_id)"),
		dal.Where("bs.connection_id = ? AND bs.board_id = ? AND s.state = ?", data.Options.ConnectionId, data.Options.BoardId, "closed"),
	}
	if apiCollector.IsIncremental() && apiCollector.GetSince() != nil {
		clauses = append(clauses, dal.Where("s.complete_date > ?", apiCollector.GetSince()))
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		logger.Error(err, "collect sprint report error")
		return err
	}

	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(apiv2models.SprintReportInput{}))
	if err != nil {
		return err
	}

	err = apiCollector.InitCollector(api.ApiCollectorArgs{
		ApiClient: data.ApiClient,
		Input:     iterator,
		// e.g. https://xxx.atlassian.net/rest/greenhopper/1.0/rapid/charts/sprintreport?rapidViewId=1&sprintId=2
		UrlTemplate: "greenhopper/1.0/rapid/charts/sprintreport",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			input := reqData.Input.(*apiv2models.SprintReportInput)
			query := url.Values{}
			query.Set("rapidViewId", fmt.Sprintf("%d", input.BoardId))
			query.Set("sprintId", fmt.Sprintf("%d", input.SprintId))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			blob, e := io.ReadAll(res.Body)
			if e != nil {
				return nil, errors.Convert(e)
			}
			return []json.RawMessage{blob}, nil
		},
		AfterResponse: ignoreHTTPStatus400,
	})
	if err != nil {
		return err
	}

	return apiCollector.Execute()
}

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
	"net/http"
	"net/url"
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core"
	. "github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

var _ core.SubTaskEntryPoint = CollectChangelogs

const RAW_CHANGELOG_TABLE = "jira_api_changelogs"

func CollectChangelogs(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	if data.JiraServerInfo.DeploymentType == models.DeploymentServer {
		return nil
	}
	db := taskCtx.GetDal()
	// figure out the time range
	since := data.Since

	// filter out issue_ids that needed collection
	clauses := []interface{}{
		Select("bi.issue_id, NOW() AS update_time"),
		From("_tool_jira_board_issues bi"),
		Join("LEFT JOIN _tool_jira_issues i ON (bi.connection_id = i.connection_id AND bi.issue_id = i.issue_id)"),
		Where(
			`bi.connection_id = ?
			   AND bi.board_id = ?
			   AND (i.changelog_updated IS NULL OR i.changelog_updated < i.updated)`,
			data.Options.ConnectionId,
			data.Options.BoardId,
		),
	}
	// apply time range if any
	if since != nil {
		clauses = append(clauses, Where("i.updated > ?", *since))
	}

	// construct the input iterator
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	// smaller struct can reduce memory footprint, we should try to avoid using big struct
	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(apiv2models.Input{}))
	if err != nil {
		return err
	}

	// now, let ApiCollector takes care the rest
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_CHANGELOG_TABLE,
		},
		ApiClient:     data.ApiClient,
		PageSize:      100,
		Incremental:   true,
		GetTotalPages: GetTotalPagesFromResponse,
		Input:         iterator,
		UrlTemplate:   "api/3/issue/{{ .Input.IssueId }}/changelog",
		Query: func(reqData *helper.RequestData, taskCtx core.SubTaskContext) (url.Values, error) {
			query := url.Values{}
			query.Set("startAt", fmt.Sprintf("%v", reqData.Pager.Skip))
			query.Set("maxResults", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		Concurrency: 10,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Values []json.RawMessage
			}
			err := helper.UnmarshalResponse(res, &data)
			if err != nil {
				return nil, err
			}
			return data.Values, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

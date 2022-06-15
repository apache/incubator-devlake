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
	. "github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

const RAW_REMOTELINK_TABLE = "jira_api_remotelinks"

var _ core.SubTaskEntryPoint = CollectRemotelinks

func CollectRemotelinks(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	logger.Info("collect remotelink")

	/*
		`CollectIssues` will take into account of `since` option and set the `updated` field for issues that have
		updates, So when it comes to collecting remotelinks, we only need to compare an issue's `updated` field with its
		`remotelink_updated` field. If `remotelink_updated` is older, then we'll collect remotelinks for this issue and
		set its `remotelink_updated` to `updated` at the end.
	*/
	cursor, err := db.Cursor(
		Select("i.issue_id, NOW() AS update_time"),
		From("_tool_jira_remotelinks i"),
		Join(`LEFT JOIN _tool_jira_board_issues bi ON (
			bi.connection_id = i.connection_id AND
			bi.issue_id = i.issue_id
		)`),
		Where(`
			bi.connection_id = ? AND
			bi.board_id = ? AND
			(i.remotelink_updated IS NULL OR i.remotelink_updated < i.updated)
			`,
			data.Options.ConnectionId,
			data.Options.BoardId,
		),
	)
	if err != nil {
		logger.Error("collect remotelink error:%v", err)
		return err
	}

	// smaller struct can reduce memory footprint, we should try to avoid using big struct
	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(apiv2models.Input{}))
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
			Table: RAW_REMOTELINK_TABLE,
		},
		ApiClient:   data.ApiClient,
		Input:       iterator,
		UrlTemplate: "api/2/issue/{{ .Input.IssueId }}/remotelink",
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			if res.StatusCode == http.StatusNotFound {
				return nil, nil
			}
			var result []json.RawMessage
			err := helper.UnmarshalResponse(res, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		},
	})
	if err != nil {
		return err
	}
	err = collector.Execute()
	if err != nil {
		return err
	}
	return nil
}

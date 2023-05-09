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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
)

const RAW_REMOTELINK_TABLE = "jira_api_remotelinks"

var _ plugin.SubTaskEntryPoint = CollectRemotelinks

var CollectRemotelinksMeta = plugin.SubTaskMeta{
	Name:             "collectRemotelinks",
	EntryPoint:       CollectRemotelinks,
	EnabledByDefault: true,
	Description:      "collect Jira remote links, supports both timeFilter and diffSync.",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func CollectRemotelinks(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	logger.Info("collect remotelink")

	collectorWithState, err := api.NewStatefulApiCollector(api.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: JiraApiParams{
			ConnectionId: data.Options.ConnectionId,
			BoardId:      data.Options.BoardId,
		},
		Table: RAW_REMOTELINK_TABLE,
	}, data.TimeAfter)
	if err != nil {
		return err
	}

	clauses := []dal.Clause{
		dal.Select("i.issue_id AS issue_id, i.updated AS update_time"),
		dal.From("_tool_jira_board_issues bi"),
		dal.Join("LEFT JOIN _tool_jira_issues i ON (bi.connection_id = i.connection_id AND bi.issue_id = i.issue_id)"),
		dal.Where("bi.connection_id=? and bi.board_id = ?", data.Options.ConnectionId, data.Options.BoardId),
	}
	incremental := collectorWithState.IsIncremental()
	if incremental && collectorWithState.LatestState.LatestSuccessStart != nil {
		clauses = append(
			clauses,
			dal.Where("i.updated > ?", collectorWithState.LatestState.LatestSuccessStart),
		)
	}

	cursor, err := db.Cursor(clauses...)
	if err != nil {
		logger.Error(err, "collect remotelink error")
		return err
	}

	// smaller struct can reduce memory footprint, we should try to avoid using big struct
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(apiv2models.Input{}))
	if err != nil {
		return err
	}

	err = collectorWithState.InitCollector(api.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		Input:       iterator,
		Incremental: incremental,
		UrlTemplate: "api/2/issue/{{ .Input.IssueId }}/remotelink",
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			if res.StatusCode == http.StatusNotFound {
				return nil, nil
			}
			var result []json.RawMessage
			err := api.UnmarshalResponse(res, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		},
		AfterResponse: ignoreHTTPStatus404,
	})
	if err != nil {
		return err
	}
	err = collectorWithState.Execute()
	if err != nil {
		return err
	}
	return nil
}

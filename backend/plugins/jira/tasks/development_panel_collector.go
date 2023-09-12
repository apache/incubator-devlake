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

const RAW_DEVELOPMENT_PANEL = "jira_api_development_panels"

var _ plugin.SubTaskEntryPoint = CollectDevelopmentPanel

var CollectDevelopmentPanelMeta = plugin.SubTaskMeta{
	Name:             "collectDevelopmentPanel",
	EntryPoint:       CollectDevelopmentPanel,
	EnabledByDefault: true,
	Description:      "collect Jira development panel",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET, plugin.DOMAIN_TYPE_CROSS},
}

func CollectDevelopmentPanel(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*JiraTaskData)
	scopeConfig := data.Options.ScopeConfig
	// if the condition is true, it means that the task is not enabled
	if scopeConfig == nil || scopeConfig.ApplicationType == "" {
		return nil
	}
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	collectorWithState, err := api.NewStatefulApiCollector(api.RawDataSubTaskArgs{
		Ctx: taskCtx,
		Params: JiraApiParams{
			ConnectionId: data.Options.ConnectionId,
			BoardId:      data.Options.BoardId,
		},

		Table: RAW_DEVELOPMENT_PANEL,
	})
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
		logger.Error(err, "collect development panel error")
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
		// the URL looks like:
		// https://merico.atlassian.net/rest/dev-status/1.0/issue/detail?issueId=25184&applicationType=GitLab&dataType=repository
		UrlTemplate: "dev-status/1.0/issue/detail",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("issueId", fmt.Sprintf("%d", reqData.Input.(*apiv2models.Input).IssueId))
			query.Set("applicationType", scopeConfig.ApplicationType)
			query.Set("dataType", "repository")
			return query, nil
		},

		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			blob, err := io.ReadAll(res.Body)
			if err != nil {
				return nil, errors.Convert(err)
			}
			var raw apiv2models.DevelopmentPanel
			err = json.Unmarshal(blob, &raw)
			if err != nil {
				return nil, errors.Convert(err)
			}
			if len(raw.Detail) == 0 {
				return nil, nil
			}
			return []json.RawMessage{blob}, nil
		},
		AfterResponse: ignoreHTTPStatus400,
	})
	if err != nil {
		return err
	}

	return collectorWithState.Execute()
}

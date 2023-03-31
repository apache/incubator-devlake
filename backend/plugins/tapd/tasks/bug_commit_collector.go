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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"net/http"
	"net/url"
	"reflect"
)

const RAW_BUG_COMMIT_TABLE = "tapd_api_bug_commits"

var _ plugin.SubTaskEntryPoint = CollectBugCommits

func CollectBugCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_BUG_COMMIT_TABLE)
	db := taskCtx.GetDal()
	collectorWithState, err := api.NewStatefulApiCollector(*rawDataSubTaskArgs, data.TimeAfter)
	if err != nil {
		return err
	}
	logger := taskCtx.GetLogger()
	logger.Info("collect issueCommits")
	incremental := collectorWithState.IsIncremental()
	clauses := []dal.Clause{
		dal.Select("_tool_tapd_bugs.id as issue_id, modified as update_time"),
		dal.From(&models.TapdBug{}),
		dal.Where("_tool_tapd_bugs.connection_id = ? and _tool_tapd_bugs.workspace_id = ? ", data.Options.ConnectionId, data.Options.WorkspaceId),
	}
	if collectorWithState.TimeAfter != nil {
		clauses = append(clauses, dal.Where("modified > ?", *collectorWithState.TimeAfter))
	}
	if incremental {
		clauses = append(clauses, dal.Where("modified > ?", *collectorWithState.LatestState.LatestSuccessStart))
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(models.Input{}))

	if err != nil {
		return err
	}
	err = collectorWithState.InitCollector(api.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		Incremental: incremental,
		Input:       iterator,
		UrlTemplate: "code_commit_infos",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			input := reqData.Input.(*models.Input)
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			query.Set("type", "bug")
			query.Set("object_id", fmt.Sprintf("%v", input.IssueId))
			query.Set("order", "created asc")
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				Stories []json.RawMessage `json:"data"`
			}
			err := api.UnmarshalResponse(res, &data)
			return data.Stories, err
		},
	})
	if err != nil {
		logger.Error(err, "collect issueCommit error")
		return err
	}
	return collectorWithState.Execute()
}

var CollectBugCommitMeta = plugin.SubTaskMeta{
	Name:             "collectBugCommits",
	EntryPoint:       CollectBugCommits,
	EnabledByDefault: true,
	Description:      "collect Tapd issueCommits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

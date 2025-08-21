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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

const RAW_STORY_COMMIT_TABLE = "tapd_api_story_commits"

var _ plugin.SubTaskEntryPoint = CollectStoryCommits

func CollectStoryCommits(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_STORY_COMMIT_TABLE)
	db := taskCtx.GetDal()
	apiCollector, err := api.NewStatefulApiCollector(*rawDataSubTaskArgs)
	if err != nil {
		return err
	}
	logger := taskCtx.GetLogger()
	logger.Info("collect issueCommits")
	clauses := []dal.Clause{
		dal.Select("_tool_tapd_stories.id as issue_id, modified as update_time"),
		dal.From(&models.TapdStory{}),
		dal.Where("_tool_tapd_stories.connection_id = ? and _tool_tapd_stories.workspace_id = ? ", data.Options.ConnectionId, data.Options.WorkspaceId),
	}
	if apiCollector.IsIncremental() && apiCollector.GetSince() != nil {
		clauses = append(clauses, dal.Where("modified > ?", *apiCollector.GetSince()))
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(models.Input{}))
	if err != nil {
		return err
	}
	err = apiCollector.InitCollector(api.ApiCollectorArgs{
		ApiClient:   data.ApiClient,
		Input:       iterator,
		UrlTemplate: "code_commit_infos",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			input := reqData.Input.(*models.Input)
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			query.Set("type", "story")
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
	return apiCollector.Execute()
}

var CollectStoryCommitMeta = plugin.SubTaskMeta{
	Name:             "collectStoryCommits",
	EntryPoint:       CollectStoryCommits,
	EnabledByDefault: true,
	Description:      "collect Tapd issueCommits",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

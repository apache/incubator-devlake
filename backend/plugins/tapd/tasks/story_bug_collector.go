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
	"fmt"
	"net/url"
	"reflect"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

const RAW_STORY_BUG_TABLE = "tapd_api_story_bugs"

var _ plugin.SubTaskEntryPoint = CollectStoryBugs

func CollectStoryBugs(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_STORY_BUG_TABLE)
	db := taskCtx.GetDal()
	apiCollector, err := api.NewStatefulApiCollector(*rawDataSubTaskArgs)
	if err != nil {
		return err
	}
	logger := taskCtx.GetLogger()
	logger.Info("collect storyBugs")

	clauses := []dal.Clause{
		dal.Select("id as issue_id, modified as update_time"),
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
		UrlTemplate: "stories/get_related_bugs",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			input := reqData.Input.(*models.Input)
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			query.Set("story_id", fmt.Sprintf("%v", input.IssueId))
			return query, nil
		},
		ResponseParser: GetRawMessageArrayFromResponse,
	})
	if err != nil {
		logger.Error(err, "collect storyBug error")
		return err
	}
	return apiCollector.Execute()
}

var CollectStoryBugMeta = plugin.SubTaskMeta{
	Name:             "collectStoryBugs",
	EntryPoint:       CollectStoryBugs,
	EnabledByDefault: false,
	Description:      "collect Tapd storyBugs",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

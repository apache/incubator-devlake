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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"net/url"
	"reflect"
)

const RAW_STORY_BUG_TABLE = "tapd_api_story_bugs"

var _ core.SubTaskEntryPoint = CollectStoryBugs

func CollectStoryBugs(taskCtx core.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_STORY_BUG_TABLE, false)
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	logger.Info("collect storyBugs")
	since := data.Since
	clauses := []dal.Clause{
		dal.Select("id as issue_id, modified as update_time"),
		dal.From(&models.TapdStory{}),
		dal.Join("LEFT JOIN _tool_tapd_story_bugs tjbc ON (tjbc.connection_id = _tool_tapd_stories.connection_id AND tjbc.story_id = _tool_tapd_stories.id)"),
		dal.Where("_tool_tapd_stories.modified > _tool_tapd_stories.created AND _tool_tapd_stories.connection_id = ? and _tool_tapd_stories.workspace_id = ? ", data.Options.ConnectionId, data.Options.WorkspaceId),
		dal.Groupby("_tool_tapd_stories.id, _tool_tapd_stories.modified"),
		dal.Having("_tool_tapd_stories.modified > max(tjbc.issue_updated) OR  max(tjbc.issue_updated) IS NULL"),
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
		ApiClient:          data.ApiClient,
		Incremental:        since == nil,
		Input:              iterator,
		UrlTemplate:        "stories/get_related_bugs",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
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
	return collector.Execute()
}

var CollectStoryBugMeta = core.SubTaskMeta{
	Name:             "collectStoryBugs",
	EntryPoint:       CollectStoryBugs,
	EnabledByDefault: true,
	Description:      "collect Tapd storyBugs",
	DomainTypes:      []string{core.DOMAIN_TYPE_TICKET},
}

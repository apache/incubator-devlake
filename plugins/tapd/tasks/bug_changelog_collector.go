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
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

const RAW_BUG_CHANGELOG_TABLE = "tapd_api_bug_changelogs"

var _ core.SubTaskEntryPoint = CollectBugChangelogs

func CollectBugChangelogs(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_BUG_CHANGELOG_TABLE)
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	logger.Info("collect storyChangelogs")
	since := data.Since
	incremental := false
	if since == nil {
		// user didn't specify a time range to sync, try load from database
		var latestUpdated models.TapdBugChangelog
		clauses := []dal.Clause{
			dal.Where("connection_id = ? and workspace_id = ?", data.Connection.ID, data.Options.WorkspaceId),
			dal.Orderby("created DESC"),
		}
		err := db.First(&latestUpdated, clauses...)
		if err != nil && err.Error() != "record not found" {
			return fmt.Errorf("failed to get latest tapd changelog record: %w", err)
		}
		if latestUpdated.Id > 0 {
			since = (*time.Time)(latestUpdated.Created)
			incremental = true
		}
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Incremental:        incremental,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		UrlTemplate:        "bug_changes",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("limit", fmt.Sprintf("%v", reqData.Pager.Size))
			query.Set("order", "created asc")
			if since != nil {
				query.Set("created", fmt.Sprintf(">%v", since.Format("YYYY-MM-DD")))
			}
			return query, nil
		},
		ResponseParser: GetRawMessageArrayFromResponse,
	})
	if err != nil {
		logger.Error("collect story changelog error:", err)
		return err
	}
	return collector.Execute()
}

var CollectBugChangelogMeta = core.SubTaskMeta{
	Name:             "collectBugChangelogs",
	EntryPoint:       CollectBugChangelogs,
	EnabledByDefault: true,
	Description:      "collect Tapd bugChangelogs",
}

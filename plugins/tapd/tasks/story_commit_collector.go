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
	"github.com/apache/incubator-devlake/errors"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"reflect"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
)

const RAW_STORY_COMMIT_TABLE = "tapd_api_story_commits"

var _ core.SubTaskEntryPoint = CollectStoryCommits

type SimpleStory struct {
	Id uint64
}

func CollectStoryCommits(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_STORY_COMMIT_TABLE, false)
	db := taskCtx.GetDal()
	logger := taskCtx.GetLogger()
	logger.Info("collect issueCommits")
	num := 0
	since := data.Since
	incremental := false
	if since == nil {
		// user didn't specify a time range to sync, try load from database
		var latestUpdated models.TapdStoryCommit
		clauses := []dal.Clause{
			dal.Where("connection_id = ? and workspace_id = ?", data.Options.ConnectionId, data.Options.WorkspaceId),
			dal.Orderby("created DESC"),
		}
		err := db.First(&latestUpdated, clauses...)
		if err != nil && err != gorm.ErrRecordNotFound {
			return errors.NotFound.Wrap(err, "failed to get latest tapd changelog record")
		}
		if latestUpdated.Id > 0 {
			since = (*time.Time)(latestUpdated.Created)
			incremental = true
		}
	}

	clauses := []dal.Clause{
		dal.Select("id"),
		dal.From(&models.TapdStory{}),
		dal.Where("connection_id = ? and workspace_id = ?", data.Options.ConnectionId, data.Options.WorkspaceId),
	}
	if since != nil {
		clauses = append(clauses, dal.Where("modified > ?", since))
	}
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	iterator, err := helper.NewDalCursorIterator(db, cursor, reflect.TypeOf(SimpleStory{}))
	if err != nil {
		return err
	}
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Incremental:        incremental,
		ApiClient:          data.ApiClient,
		//PageSize:    100,
		Input:       iterator,
		UrlTemplate: "code_commit_infos",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			input := reqData.Input.(*SimpleStory)
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			query.Set("type", "story")
			query.Set("object_id", fmt.Sprintf("%v", input.Id))
			query.Set("order", "created asc")
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Stories []json.RawMessage `json:"data"`
			}
			if len(data.Stories) > 0 {
				fmt.Println(len(data.Stories))
				num += len(data.Stories)
				fmt.Printf("num is %d", num)
			}
			err := helper.UnmarshalResponse(res, &data)
			return data.Stories, err
		},
	})
	if err != nil {
		logger.Error("collect issueCommit error:", err)
		return err
	}
	return collector.Execute()
}

var CollectStoryCommitMeta = core.SubTaskMeta{
	Name:             "collectStoryCommits",
	EntryPoint:       CollectStoryCommits,
	EnabledByDefault: true,
	Description:      "collect Tapd issueCommits",
	DomainTypes:      []string{core.DOMAIN_TYPE_CROSS},
}

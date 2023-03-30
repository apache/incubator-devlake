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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"net/http"
	"net/url"
)

const RAW_STORY_CATEGORY_TABLE = "tapd_api_story_categories"

var _ plugin.SubTaskEntryPoint = CollectStoryCategories

func CollectStoryCategories(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_STORY_CATEGORY_TABLE)
	logger := taskCtx.GetLogger()
	logger.Info("collect story_category")
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		UrlTemplate:        "story_categories",
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data struct {
				StoryCategories []json.RawMessage `json:"data"`
			}
			err := api.UnmarshalResponse(res, &data)
			return data.StoryCategories, err
		},
	})
	if err != nil {
		logger.Error(err, "collect story_category error")
		return err
	}
	return collector.Execute()
}

var CollectStoryCategoriesMeta = plugin.SubTaskMeta{
	Name:             "collectStoryCategories",
	EntryPoint:       CollectStoryCategories,
	EnabledByDefault: true,
	Description:      "collect Tapd StoryCategories",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

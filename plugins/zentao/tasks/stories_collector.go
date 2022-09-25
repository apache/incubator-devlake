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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"io"
	"net/http"
	"net/url"
)

const RAW_STORIES_TABLE = "zentao_stories"

var _ core.SubTaskEntryPoint = CollectExecution

func CollectStories(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*ZentaoTaskData)
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: ZentaoApiParams{
				StoriesId: data.Options.StoriesId,
			},
			Table: RAW_STORIES_TABLE,
		},
		ApiClient:   data.ApiClient,
		Incremental: false,
		PageSize:    100,
		// TODO write which api would you want request
		UrlTemplate: "/stories/{{ .Params.StoriesId }}",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("limit", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			body, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				return nil, err
			}
			return []json.RawMessage{body}, nil
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

var CollectStoriesMeta = core.SubTaskMeta{
	Name:             "CollectStories",
	EntryPoint:       CollectStories,
	EnabledByDefault: true,
	Description:      "Collect Stories data from Zentao api",
}

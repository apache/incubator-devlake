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

const RAW_TASK_TAG_TABLE = "teambition_api_task_tags"

var _ plugin.SubTaskEntryPoint = CollectTasks

var CollectTaskTagsMeta = plugin.SubTaskMeta{
	Name:             "collectTaskTags",
	EntryPoint:       CollectTaskTags,
	EnabledByDefault: true,
	Description:      "collect teambition task tags",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_TICKET},
}

func CollectTaskTags(taskCtx plugin.SubTaskContext) errors.Error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_TASK_TAG_TABLE)
	logger := taskCtx.GetLogger()
	logger.Info("collect task tags")

	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           int(data.Options.PageSize),
		UrlTemplate:        "/v3/project/{{ .Params.ProjectId }}/tag/search",
		GetNextPageCustomData: func(prevReqData *api.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
			res := TeambitionComRes[any]{}
			err := api.UnmarshalResponse(prevPageResponse, &res)
			if err != nil {
				return nil, err
			}
			if res.NextPageToken == "" {
				return nil, api.ErrFinishCollect
			}
			return res.NextPageToken, nil
		},
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			if data.Options.PageSize > 0 {
				query.Set("pageSize", fmt.Sprintf("%v", data.Options.PageSize))
			}
			if pageToken, ok := reqData.CustomData.(string); ok && pageToken != "" {
				query.Set("pageToken", reqData.CustomData.(string))
			}
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			var data = TeambitionComRes[[]json.RawMessage]{}
			err := api.UnmarshalResponse(res, &data)
			return data.Result, err
		},
	})
	if err != nil {
		logger.Error(err, "collect task error")
		return err
	}
	return collector.Execute()
}

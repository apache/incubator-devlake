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
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/feishu/apimodels"
)

const RAW_CHAT_TABLE = "feishu_chat_item"

var _ plugin.SubTaskEntryPoint = CollectChat

// CollectChat collect all chats that bot is in
func CollectChat(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*FeishuTaskData)
	pageSize := 50
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: FeishuApiParams{
				ConnectionId: data.Options.ConnectionId,
			},
			Table: RAW_CHAT_TABLE,
		},
		ApiClient:   data.ApiClient,
		Incremental: false,
		UrlTemplate: "im/v1/chats",
		PageSize:    pageSize,
		GetNextPageCustomData: func(prevReqData *api.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
			res := apimodels.FeishuImApiResult{}
			err := api.UnmarshalResponse(prevPageResponse, &res)
			if err != nil {
				return nil, err
			}
			if !res.Data.HasMore {
				return nil, api.ErrFinishCollect
			}
			return res.Data.PageToken, nil
		},
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			query.Set("page_size", strconv.Itoa(pageSize))
			if pageToken, ok := reqData.CustomData.(string); ok && pageToken != "" {
				query.Set("page_token", reqData.CustomData.(string))
			}
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			body := &apimodels.FeishuImApiResult{}
			err := api.UnmarshalResponse(res, body)
			if err != nil {
				return nil, err
			}
			return body.Data.Items, nil
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

var CollectChatMeta = plugin.SubTaskMeta{
	Name:             "collectChat",
	EntryPoint:       CollectChat,
	EnabledByDefault: true,
	Description:      "Collect chats from Feishu api",
}

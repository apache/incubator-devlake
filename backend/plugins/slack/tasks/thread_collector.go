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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/slack/apimodels"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

const RAW_THREAD_TABLE = "slack_thread"

var _ plugin.SubTaskEntryPoint = CollectThread

type ThreadInput struct {
	ChannelId string `json:"channel_id"`
	ThreadTs  string `json:"thread_ts"`
}

func CollectThread(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*SlackTaskData)
	db := taskCtx.GetDal()

	clauses := []dal.Clause{
		dal.Select("thread_ts, channel_id"),
		dal.From("_tool_slack_channel_messages"),
		dal.Where("connection_id=? AND thread_ts!='' AND subtype=''", data.Options.ConnectionId),
	}

	// construct the input iterator
	cursor, err := db.Cursor(clauses...)
	if err != nil {
		return err
	}
	// smaller struct can reduce memory footprint, we should try to avoid using big struct
	iterator, err := api.NewDalCursorIterator(db, cursor, reflect.TypeOf(ThreadInput{}))
	if err != nil {
		return err
	}

	pageSize := 50
	collector, err := api.NewApiCollector(api.ApiCollectorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: SlackApiParams{
				ConnectionId: data.Options.ConnectionId,
			},
			Table: RAW_THREAD_TABLE,
		},
		ApiClient:   data.ApiClient,
		Incremental: false,
		Input:       iterator,
		UrlTemplate: "conversations.replies",
		PageSize:    pageSize,
		GetNextPageCustomData: func(prevReqData *api.RequestData, prevPageResponse *http.Response) (interface{}, errors.Error) {
			res := apimodels.SlackThreadsApiResult{}
			err := api.UnmarshalResponse(prevPageResponse, &res)
			if err != nil {
				return nil, err
			}
			if res.ResponseMetadata.NextCursor == "" {
				return nil, api.ErrFinishCollect
			}
			return res.ResponseMetadata.NextCursor, nil
		},
		Query: func(reqData *api.RequestData) (url.Values, errors.Error) {
			input := reqData.Input.(*ThreadInput)
			query := url.Values{}
			query.Set("channel", input.ChannelId)
			query.Set("ts", input.ThreadTs)
			query.Set("offset", strconv.Itoa(reqData.Pager.Skip))
			query.Set("limit", strconv.Itoa(pageSize))
			if pageToken, ok := reqData.CustomData.(string); ok && pageToken != "" {
				query.Set("cursor", reqData.CustomData.(string))
			}
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			body := &apimodels.SlackThreadsApiResult{}
			err := api.UnmarshalResponse(res, body)
			if err != nil {
				return nil, err
			}
			return body.Threads, nil
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

var CollectThreadMeta = plugin.SubTaskMeta{
	Name:             "collectThread",
	EntryPoint:       CollectThread,
	EnabledByDefault: true,
	Description:      "Collect thread from Slack api",
}

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
	"github.com/apache/incubator-devlake/errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/feishu/apimodels"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_MEETING_TOP_USER_ITEM_TABLE = "feishu_meeting_top_user_item"

var _ core.SubTaskEntryPoint = CollectMeetingTopUserItem

func CollectMeetingTopUserItem(taskCtx core.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*FeishuTaskData)
	pageSize := 100
	NumOfDaysToCollectInt := int(data.Options.NumOfDaysToCollect)
	iterator, err := helper.NewDateIterator(NumOfDaysToCollectInt)
	if err != nil {
		return err
	}
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: FeishuApiParams{
				ConnectionId: data.Options.ConnectionId,
			},
			Table: RAW_MEETING_TOP_USER_ITEM_TABLE,
		},
		ApiClient:   data.ApiClient,
		Incremental: false,
		Input:       iterator,
		UrlTemplate: "/reports/get_top_user",
		Query: func(reqData *helper.RequestData) (url.Values, errors.Error) {
			query := url.Values{}
			input := reqData.Input.(*helper.DatePair)
			query.Set("start_time", strconv.FormatInt(input.PairStartTime.Unix(), 10))
			query.Set("end_time", strconv.FormatInt(input.PairEndTime.Unix(), 10))
			query.Set("limit", strconv.Itoa(pageSize))
			query.Set("order_by", "2")
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, errors.Error) {
			body := &apimodels.FeishuMeetingTopUserItemResult{}
			err := helper.UnmarshalResponse(res, body)
			if err != nil {
				return nil, err
			}
			return body.Data.TopUserReport, nil
		},
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

var CollectMeetingTopUserItemMeta = core.SubTaskMeta{
	Name:             "collectMeetingTopUserItem",
	EntryPoint:       CollectMeetingTopUserItem,
	EnabledByDefault: true,
	Description:      "Collect top user meeting data from Feishu api",
}

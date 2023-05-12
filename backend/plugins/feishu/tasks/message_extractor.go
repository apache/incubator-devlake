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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/feishu/apimodels"
	"github.com/apache/incubator-devlake/plugins/feishu/models"
	"strconv"
	"time"
)

var _ plugin.SubTaskEntryPoint = ExtractMessage

func ExtractMessage(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*FeishuTaskData)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: FeishuApiParams{
				ConnectionId: data.Options.ConnectionId,
			},
			Table: RAW_MESSAGE_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			body := &apimodels.FeishuMessageResultItem{}
			err := errors.Convert(json.Unmarshal(row.Data, body))
			if err != nil {
				return nil, err
			}
			message := &models.FeishuMessage{}
			message.ConnectionId = data.Options.ConnectionId
			message.MessageId = body.MessageId
			message.Content = body.Body.Content
			message.ChatId = body.ChatId
			message.MsgType = body.MsgType
			message.ParentId = body.ParentId
			message.RootId = body.RootId
			message.SenderId = body.Sender.Id
			message.SenderIdType = body.Sender.IdType
			message.SenderType = body.Sender.SenderType
			message.Deleted = body.Deleted
			createTimestamp, err := errors.Convert01(strconv.Atoi(body.CreateTime))
			if err != nil {
				return nil, err
			}
			message.CreateTime = time.UnixMilli(int64(createTimestamp))
			updateTimestamp, err := errors.Convert01(strconv.Atoi(body.UpdateTime))
			if err != nil {
				return nil, err
			}
			message.UpdateTime = time.UnixMilli(int64(updateTimestamp))
			message.Updated = body.Updated
			return []interface{}{message}, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractMessageMeta = plugin.SubTaskMeta{
	Name:             "extractChatItem",
	EntryPoint:       ExtractMessage,
	EnabledByDefault: true,
	Description:      "Extract raw messages data into tool layer table feishu_meeting_top_user_item",
}

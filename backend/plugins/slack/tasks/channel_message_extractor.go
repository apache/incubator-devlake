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
	"github.com/apache/incubator-devlake/plugins/slack/apimodels"
	"github.com/apache/incubator-devlake/plugins/slack/models"
)

var _ plugin.SubTaskEntryPoint = ExtractChannelMessage

func ExtractChannelMessage(taskCtx plugin.SubTaskContext) errors.Error {
	data := taskCtx.GetData().(*SlackTaskData)
	extractor, err := api.NewApiExtractor(api.ApiExtractorArgs{
		RawDataSubTaskArgs: api.RawDataSubTaskArgs{
			Ctx:     taskCtx,
			Options: data.Options,
			Table:   RAW_CHANNEL_MESSAGE_TABLE,
		},
		Extract: func(row *api.RawData) ([]interface{}, errors.Error) {
			channel := &ChannelInput{}
			err := errors.Convert(json.Unmarshal(row.Input, channel))
			if err != nil {
				return nil, err
			}

			body := &apimodels.SlackChannelMessageResultItem{}
			err = errors.Convert(json.Unmarshal(row.Data, body))
			if err != nil {
				return nil, err
			}
			message := &models.SlackChannelMessage{}
			message.ConnectionId = data.Options.ConnectionId
			message.ChannelId = channel.ChannelId
			message.ClientMsgId = body.ClientMsgId
			message.Type = body.Type
			message.Subtype = body.Subtype
			message.Ts = body.Ts
			message.ThreadTs = body.ThreadTs
			message.User = body.User
			message.Text = body.Text
			message.Team = body.Team
			message.ReplyCount = body.ReplyCount
			message.ReplyUsersCount = body.ReplyUsersCount
			message.LatestReply = body.LatestReply
			message.IsLocked = body.IsLocked
			message.Subscribed = body.Subscribed
			message.ParentUserId = body.ParentUserId
			return []interface{}{message}, nil
		},
	})
	if err != nil {
		return err
	}

	return extractor.Execute()
}

var ExtractChannelMessageMeta = plugin.SubTaskMeta{
	Name:             "extractChannelMessage",
	EntryPoint:       ExtractChannelMessage,
	EnabledByDefault: true,
	Description:      "Extract raw channel messages data into tool layer table",
	DomainTypes:      []string{plugin.DOMAIN_TYPE_CROSS},
}

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

package models

import (
	"github.com/apache/incubator-devlake/core/models/common"
)

type SlackChannelMessage struct {
	common.NoPKModel `json:"-"`
	ConnectionId     uint64 `gorm:"primaryKey"`
	ChannelId        string `json:"channel_id" gorm:"primaryKey"`
	Ts               string `json:"ts" gorm:"primaryKey"`
	ClientMsgId      string `json:"client_msg_id"`
	Type             string `json:"type"`
	Subtype          string `json:"subtype"`
	ThreadTs         string `json:"thread_ts"`
	User             string `json:"user"`
	Text             string `json:"text"`
	Team             string `json:"team"`
	ReplyCount       int    `json:"reply_count"`
	ReplyUsersCount  int    `json:"reply_users_count"`
	LatestReply      string `json:"latest_reply"`
	IsLocked         bool   `json:"is_locked"`
	Subscribed       bool   `json:"subscribed"`
	ParentUserId     string `json:"parent_user_id"`
}

func (SlackChannelMessage) TableName() string {
	return "_tool_slack_channel_messages"
}

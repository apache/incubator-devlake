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

package archived

import (
	"github.com/apache/incubator-devlake/core/models/migrationscripts/archived"
	"time"
)

type FeishuMessage struct {
	archived.NoPKModel `json:"-"`
	ConnectionId       uint64    `gorm:"primaryKey"`
	MessageId          string    `json:"message_id" gorm:"primaryKey"`
	Content            string    `json:"content"`
	ChatId             string    `json:"chat_id"`
	MsgType            string    `json:"msg_type"`
	ParentId           string    `json:"parent_id"`
	RootId             string    `json:"root_id"`
	SenderId           string    `json:"id"`
	SenderIdType       string    `json:"id_type"`
	SenderType         string    `json:"sender_type"`
	Deleted            bool      `json:"deleted"`
	CreateTime         time.Time `json:"create_time"`
	UpdateTime         time.Time `json:"update_time"`
	Updated            bool      `json:"updated"`
}

func (FeishuMessage) TableName() string {
	return "_tool_feishu_messages"
}

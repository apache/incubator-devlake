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

package apimodels

import "encoding/json"

type FeishuImApiResult struct {
	Code int `json:"code"`
	Data struct {
		HasMore   bool              `json:"has_more"`
		Items     []json.RawMessage `json:"items"`
		PageToken string            `json:"page_token"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type FeishuMessageResultItem struct {
	Body struct {
		Content string `json:"content"`
	} `json:"body"`
	ChatId     string `json:"chat_id"`
	CreateTime string `json:"create_time"`
	Deleted    bool   `json:"deleted"`
	Mentions   []struct {
		Id        string `json:"id"`
		IdType    string `json:"id_type"`
		Key       string `json:"key"`
		Name      string `json:"name"`
		TenantKey string `json:"tenant_key"`
	} `json:"mentions"`
	MessageId string `json:"message_id"`
	MsgType   string `json:"msg_type"`
	ParentId  string `json:"parent_id"`
	RootId    string `json:"root_id"`
	Sender    struct {
		Id         string `json:"id"`
		IdType     string `json:"id_type"`
		SenderType string `json:"sender_type"`
		TenantKey  string `json:"tenant_key"`
	} `json:"sender"`
	UpdateTime string `json:"update_time"`
	Updated    bool   `json:"updated"`
}

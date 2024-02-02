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
	"github.com/apache/incubator-devlake/core/utils"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

// SlackConn holds the essential information to connect to the Slack API
type SlackConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.AccessToken    `mapstructure:",squash"`
}

func (connection SlackConn) Sanitize() SlackConn {
	connection.Token = utils.SanitizeString(connection.Token)
	return connection
}

// SlackConnection holds SlackConn plus ID/Name for database storage
type SlackConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	SlackConn             `mapstructure:",squash"`
}

func (SlackConnection) TableName() string {
	return "_tool_slack_connections"
}

func (connection SlackConnection) Sanitize() SlackConnection {
	connection.SlackConn = connection.SlackConn.Sanitize()
	return connection
}

func (connection *SlackConnection) MergeFromRequest(target *SlackConnection, body map[string]interface{}) error {
	token := target.Token
	if err := helper.DecodeMapStruct(body, target, true); err != nil {
		return err
	}
	modifiedToken := target.Token
	if modifiedToken == "" || modifiedToken == utils.SanitizeString(token) {
		target.Token = token
	}
	return nil
}

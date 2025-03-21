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
type ArgoConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.AccessToken    `mapstructure:",squash"`
}

func (connection ArgoConn) Sanitize() ArgoConn {
	connection.Token = utils.SanitizeString(connection.Token)
	return connection
}

// ArgoConnection holds ArgoConn plus ID/Name for database storage
type ArgoConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	ArgoConn              `mapstructure:",squash"`
}

func (ArgoConnection) TableName() string {
	return "_tool_argo_connections"
}

func (connection ArgoConnection) Sanitize() ArgoConnection {
	connection.ArgoConn = connection.ArgoConn.Sanitize()
	return connection
}

func (connection *ArgoConnection) MergeFromRequest(target *ArgoConnection, body map[string]interface{}) error {
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

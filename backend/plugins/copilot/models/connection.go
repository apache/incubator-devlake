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

// CopilotConn stores GitHub Copilot connection settings.
type CopilotConn struct {
	helper.RestConnection `mapstructure:",squash"`

	Token            string `mapstructure:"token" json:"token"`
	Organization     string `mapstructure:"organization" json:"organization"`
	RateLimitPerHour int    `mapstructure:"rateLimitPerHour" json:"rateLimitPerHour"`
}

func (conn *CopilotConn) Sanitize() CopilotConn {
	if conn == nil {
		return CopilotConn{}
	}
	clone := *conn
	clone.Token = utils.SanitizeString(clone.Token)
	return clone
}

// CopilotConnection persists connection details with metadata required by DevLake.
type CopilotConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	CopilotConn           `mapstructure:",squash"`
}

func (CopilotConnection) TableName() string {
	return "_tool_copilot_connections"
}

func (connection CopilotConnection) Sanitize() CopilotConnection {
	connection.CopilotConn = connection.CopilotConn.Sanitize()
	return connection
}

func (connection *CopilotConnection) MergeFromRequest(target *CopilotConnection, body map[string]interface{}) error {
	if target == nil {
		return nil
	}
	originalToken := target.Token
	if err := helper.DecodeMapStruct(body, target, true); err != nil {
		return err
	}
	sanitizedOriginal := utils.SanitizeString(originalToken)
	if target.Token == "" || target.Token == sanitizedOriginal {
		target.Token = originalToken
	}
	return nil
}

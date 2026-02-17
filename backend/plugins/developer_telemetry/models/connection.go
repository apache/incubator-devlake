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

// DeveloperTelemetryConn holds the essential information to connect to the Developer Telemetry webhook
type DeveloperTelemetryConn struct {
	helper.RestConnection `mapstructure:",squash"`
	SecretToken           string `mapstructure:"secretToken" json:"secretToken" gorm:"serializer:encdec"`
}

func (tc *DeveloperTelemetryConn) Sanitize() DeveloperTelemetryConn {
	tc.SecretToken = utils.SanitizeString(tc.SecretToken)
	return *tc
}

// DeveloperTelemetryConnection holds DeveloperTelemetryConn plus ID/Name for database storage
type DeveloperTelemetryConnection struct {
	helper.BaseConnection  `mapstructure:",squash"`
	DeveloperTelemetryConn `mapstructure:",squash"`
}

func (DeveloperTelemetryConnection) TableName() string {
	return "_tool_developer_telemetry_connections"
}

func (connection *DeveloperTelemetryConnection) MergeFromRequest(target *DeveloperTelemetryConnection, body map[string]interface{}) error {
	secretToken := target.SecretToken

	if err := helper.DecodeMapStruct(body, target, true); err != nil {
		return err
	}

	modifiedToken := target.SecretToken

	// preserve existing token if not modified
	if modifiedToken == "" || modifiedToken == utils.SanitizeString(secretToken) {
		target.SecretToken = secretToken
	}

	return nil
}

func (connection DeveloperTelemetryConnection) Sanitize() DeveloperTelemetryConnection {
	connection.DeveloperTelemetryConn = connection.DeveloperTelemetryConn.Sanitize()
	return connection
}

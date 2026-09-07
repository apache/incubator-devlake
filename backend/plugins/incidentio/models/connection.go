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
	"fmt"
	"net/http"

	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/utils"
	helper "github.com/apache/devlake/helpers/pluginhelper/api"
)

type IncidentioAccessToken helper.AccessToken

func (at *IncidentioAccessToken) SetupAuthentication(request *http.Request) errors.Error {
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", at.Token))
	return nil
}

type IncidentioConn struct {
	helper.RestConnection `mapstructure:",squash"`
	IncidentioAccessToken `mapstructure:",squash"`
}

func (connection IncidentioConn) Sanitize() IncidentioConn {
	connection.Token = utils.SanitizeString(connection.Token)
	return connection
}

type IncidentioConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	IncidentioConn        `mapstructure:",squash"`
}

// MergeFromRequest preserves the existing token when an incoming PATCH
// body omits it or echoes the sanitized form. The config-UI sends the
// sanitized token back on every PATCH to avoid round-tripping the
// secret; this guard is what makes that pattern safe.
func (connection *IncidentioConnection) MergeFromRequest(target *IncidentioConnection, body map[string]interface{}) error {
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

func (IncidentioConnection) TableName() string {
	return "_tool_incidentio_connections"
}

func (connection IncidentioConnection) Sanitize() IncidentioConnection {
	connection.Token = utils.SanitizeString(connection.Token)
	return connection
}

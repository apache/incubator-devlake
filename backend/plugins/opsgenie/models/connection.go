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
	"github.com/apache/incubator-devlake/core/utils"
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

// AccessToken implements HTTP Token Authentication with Access Token
type OpsgenieGenieKey helper.AccessToken

// SetupAuthentication sets up the request headers for authentication
func (gk *OpsgenieGenieKey) SetupAuthentication(request *http.Request) errors.Error {
	request.Header.Set("Authorization", fmt.Sprintf("GenieKey %s", gk.Token))
	return nil
}

// OpsgenieConn holds the essential information to connect to the Opsgenie API
type OpsgenieConn struct {
	helper.RestConnection `mapstructure:",squash"`
	OpsgenieGenieKey      `mapstructure:",squash"`
}

func (connection OpsgenieConn) Sanitize() OpsgenieConn {
	connection.Token = utils.SanitizeString(connection.Token)
	return connection
}

// OpsgenieConnection holds OpsgenieConn plus ID/Name for database storage
type OpsgenieConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	OpsgenieConn          `mapstructure:",squash"`
}

func (connection *OpsgenieConnection) MergeFromRequest(target *OpsgenieConnection, body map[string]interface{}) error {
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

func (connection OpsgenieConnection) Sanitize() OpsgenieConnection {
	connection.OpsgenieConn = connection.OpsgenieConn.Sanitize()
	return connection
}

// This object conforms to what the frontend currently expects.
type OpsgenieResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	OpsgenieConnection
}

// Using User because it requires authentication.
type ApiUserResponse struct {
	Id   int
	Name string `json:"name"`
}

func (OpsgenieConnection) TableName() string {
	return "_tool_opsgenie_connections"
}

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
	"net/http"

	"github.com/apache/devlake/core/errors"
	"github.com/apache/devlake/core/utils"
	helper "github.com/apache/devlake/helpers/pluginhelper/api"
)

// ClickUpConn holds the essential information to connect to the ClickUp API.
// ClickUp authenticates with a personal API token passed verbatim in the
// `Authorization` header (NO `Bearer` prefix), so we implement our own
// SetupAuthentication instead of reusing helper.AccessToken. OAuth2 is out of
// scope for now.
type ClickUpConn struct {
	helper.RestConnection `mapstructure:",squash"`
	Token                 string `mapstructure:"token" validate:"required" json:"token" gorm:"serializer:encdec"`
}

// SetupAuthentication sets up the HTTP request authentication for the ClickUp API.
func (cc *ClickUpConn) SetupAuthentication(req *http.Request) errors.Error {
	req.Header.Set("Authorization", cc.Token)
	return nil
}

func (cc *ClickUpConn) Sanitize() ClickUpConn {
	cc.Token = utils.SanitizeString(cc.Token)
	return *cc
}

// ClickUpConnection holds ClickUpConn plus ID/Name for database storage.
type ClickUpConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	ClickUpConn           `mapstructure:",squash"`
}

func (connection ClickUpConnection) Sanitize() ClickUpConnection {
	connection.ClickUpConn = connection.ClickUpConn.Sanitize()
	return connection
}

func (connection *ClickUpConnection) MergeFromRequest(target *ClickUpConnection, body map[string]interface{}) error {
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

func (ClickUpConnection) TableName() string {
	return "_tool_clickup_connections"
}

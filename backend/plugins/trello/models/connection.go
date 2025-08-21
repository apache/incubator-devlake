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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/utils"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"net/http"
)

// TrelloConn holds the essential information to connect to the Trello API
type TrelloConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.AppKey         `mapstructure:",squash"`
}

func (tc *TrelloConn) Sanitize() TrelloConn {
	tc.SecretKey = utils.SanitizeString(tc.SecretKey)
	return *tc
}

// TrelloConnection holds TrelloConn plus ID/Name for database storage
type TrelloConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	TrelloConn            `mapstructure:",squash"`
}

func (connection *TrelloConnection) MergeFromRequest(target *TrelloConnection, body map[string]interface{}) error {
	secretKey := target.SecretKey
	if err := helper.DecodeMapStruct(body, target, true); err != nil {
		return err
	}
	modifiedSecretKey := target.SecretKey
	if modifiedSecretKey == "" || modifiedSecretKey == utils.SanitizeString(secretKey) {
		target.SecretKey = secretKey
	}
	return nil
}

func (connection TrelloConnection) Sanitize() TrelloConnection {
	connection.TrelloConn = connection.TrelloConn.Sanitize()
	return connection
}

// SetupAuthentication sets up the HTTP Request Authentication
func (tc *TrelloConn) SetupAuthentication(req *http.Request) errors.Error {
	req.Header.Set("Authorization", fmt.Sprintf("OAuth oauth_consumer_key=\"%s\", oauth_token=\"%s\"", tc.AppId, tc.SecretKey))
	return nil
}

func (TrelloConnection) TableName() string {
	return "_tool_trello_connections"
}

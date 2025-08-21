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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/utils"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"net/http"
)

// CircleciConn holds the essential information to connect to the CircleciConn API
type CircleciConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.AccessToken    `mapstructure:",squash"`
}

// CircleciConnection holds CircleciConn plus ID/Name for database storage
type CircleciConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	CircleciConn          `mapstructure:",squash"`
}

func (cc *CircleciConn) SetupAuthentication(req *http.Request) errors.Error {
	req.Header.Set("Circle-Token", cc.Token)
	return nil
}

func (CircleciConnection) TableName() string {
	return "_tool_circleci_connections"
}

func (connection CircleciConnection) Sanitize() CircleciConnection {
	connection.Token = utils.SanitizeString(connection.Token)
	return connection
}

func (connection *CircleciConnection) MergeFromRequest(target *CircleciConnection, body map[string]interface{}) error {
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

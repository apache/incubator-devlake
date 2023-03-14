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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"net/http"
)

// TrelloConn holds the essential information to connect to the Trello API
type TrelloConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.AppKey         `mapstructure:",squash"`
}

// TrelloConnection holds TrelloConn plus ID/Name for database storage
type TrelloConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	TrelloConn            `mapstructure:",squash"`
}

// SetupAuthentication sets up the HTTP Request Authentication
func (tc *TrelloConn) SetupAuthentication(req *http.Request) errors.Error {
	req.Header.Set("Authorization", fmt.Sprintf("OAuth oauth_consumer_key=\"%s\", oauth_token=\"%s\"", tc.AppId, tc.SecretKey))
	return nil
}

func (TrelloConnection) TableName() string {
	return "_tool_trello_connections"
}

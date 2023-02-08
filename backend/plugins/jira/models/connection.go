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

	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type EpicResponse struct {
	Id    int
	Title string
	Value string
}

type BoardResponse struct {
	Id    int
	Title string
	Value string
}

// JiraConn holds the essential information to connect to the Jira API
type JiraConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.MultiAuth      `mapstructure:",squash"`
	helper.BasicAuth      `mapstructure:",squash"`
	helper.AccessToken    `mapstructure:",squash"`
}

// SetupAuthentication implements the `IAuthentication` interface by delegating
// the actual logic to the `MultiAuth` struct to help us write less code
func (jc *JiraConn) SetupAuthentication(req *http.Request) errors.Error {
	return jc.MultiAuth.SetupAuthenticationForConnection(jc, req)
}

// JiraConnection holds JiraConn plus ID/Name for database storage
type JiraConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	JiraConn              `mapstructure:",squash"`
}

func (JiraConnection) TableName() string {
	return "_tool_jira_connections"
}

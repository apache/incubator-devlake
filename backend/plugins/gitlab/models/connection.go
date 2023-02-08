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

	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

// GitlabConn holds the essential information to connect to the Gitlab API
type GitlabConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.AccessToken    `mapstructure:",squash"`
	IsPrivateToken        bool `gorm:"-"`
	UserId                int  `gorm:"-"`
}

// SetupAuthentication sets up the HTTP Request Authentication
func (gak GitlabConn) SetupAuthentication(req *http.Request) errors.Error {
	if gak.IsPrivateToken {
		req.Header.Set("Private-Token", gak.Token)
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", gak.Token))
	}
	return nil
}

// This object conforms to what the frontend currently sends.
// GitlabConnection holds GitlabConn plus ID/Name for database storage
type GitlabConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	GitlabConn            `mapstructure:",squash"`
}

// This object conforms to what the frontend currently expects.
type GitlabResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	GitlabConnection
}

// Using User because it requires authentication.
type ApiUserResponse struct {
	Id   int
	Name string `json:"name"`
}

func (GitlabConnection) TableName() string {
	return "_tool_gitlab_connections"
}

type AeAppKey helper.AppKey

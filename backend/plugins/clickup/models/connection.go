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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type ClickupConnection struct {
	api.BaseConnection    `mapstructure:",squash"`
	helper.RestConnection `mapstructure:",squash"`
	helper.AccessToken    `mapstructure:",squash"`
	TeamId                string `gorm:"type:varchar(255)"`
}

type TestConnectionRequest struct {
	Endpoint           string `json:"endpoint"`
	Proxy              string `json:"proxy"`
	helper.AccessToken `mapstructure:",squash"`
}

// This object conforms to what the frontend currently expects.
type ClickupResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	ClickupConnection
}

// Using User because it requires authentication.
type ApiUserResponse struct {
	Id   int
	Name string `json:"name"`
}

func (ClickupConnection) TableName() string {
	return "_tool_clickup_connections"
}

func (c *ClickupConnection) SetupAuthentication(req *http.Request) errors.Error {
	req.Header.Add("Authorization", c.Token)
	return nil
}

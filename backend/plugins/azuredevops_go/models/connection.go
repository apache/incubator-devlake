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
	"encoding/base64"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/core/utils"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"net/http"
)

var _ plugin.ApiConnection = (*AzuredevopsConn)(nil)
var _ plugin.AccessTokenAuthenticator = (*AzuredevopsAccessToken)(nil)

// AzuredevopsAccessToken implements HTTP Bearer Authentication with Access Token
type AzuredevopsAccessToken struct {
	Token string `mapstructure:"token" validate:"required" json:"token" gorm:"serializer:encdec"`
}

// GetAccessTokenAuthenticator returns SetupAuthentication
func (at *AzuredevopsAccessToken) GetAccessTokenAuthenticator() plugin.ApiAuthenticator {
	return at
}

// SetupAuthentication sets up the HTTP Request Authentication
func (at *AzuredevopsAccessToken) SetupAuthentication(req *http.Request) errors.Error {
	h := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(":%s", at.Token)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %v", h))

	return nil
}

// AzuredevopsConn holds the essential information to connect to the Azure DevOps API
type AzuredevopsConn struct {
	//api.RestConnection `mapstructure:",squash"`
	AzuredevopsAccessToken `mapstructure:",squash"`
	Organization           string `json:"organization"`
	//Endpoint         string `mapstructure:"endpoint" json:"endpoint"`
	Proxy string `mapstructure:"proxy" json:"proxy"`
	//RateLimitPerHour int    `comment:"api request rate limit per hour" json:"rateLimitPerHour"`
}

func (conn *AzuredevopsConn) GetEndpoint() string {
	return "https://dev.azure.com"
}

func (conn *AzuredevopsConn) GetProxy() string {
	return conn.Proxy
}

func (conn *AzuredevopsConn) GetRateLimitPerHour() int {
	return 0
}

var _ plugin.ApiConnection = (*AzuredevopsConnection)(nil)

// AzuredevopsConnection holds AzuredevopsConn plus ID/Name for database storage
type AzuredevopsConnection struct {
	api.BaseConnection `mapstructure:",squash"`
	AzuredevopsConn    `mapstructure:",squash"`
}

func (c AzuredevopsConnection) GetEndpoint() string {
	return "https://dev.azure.com"
}

func (c AzuredevopsConnection) GetProxy() string {
	return c.Proxy
}

func (c AzuredevopsConnection) GetRateLimitPerHour() int {
	return 0
}

func (AzuredevopsConnection) TableName() string {
	return "_tool_azuredevops_go_connections"
}

func (conn *AzuredevopsConn) Sanitize() AzuredevopsConn {
	conn.Token = utils.SanitizeString(conn.Token)
	return *conn
}

func (c AzuredevopsConnection) Sanitize() AzuredevopsConnection {
	c.AzuredevopsConn = c.AzuredevopsConn.Sanitize()
	return c
}

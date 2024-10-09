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
	"net/http"

	"github.com/apache/incubator-devlake/core/utils"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type SonarqubeAccessToken helper.AccessToken

// SetupAuthentication sets up the HTTP Request Authentication
func (sat SonarqubeAccessToken) SetupAuthentication(req *http.Request) errors.Error {
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", sat.GetEncodedToken()))
	return nil
}

func (sat SonarqubeAccessToken) GetAccessTokenAuthenticator() plugin.ApiAuthenticator {
	return sat
}

// GetEncodedToken returns encoded bearer token for HTTP Basic Authentication
func (sat SonarqubeAccessToken) GetEncodedToken() string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:", sat.Token)))
}

// SonarqubeConn holds the essential information to connect to the sonarqube API
type SonarqubeConn struct {
	helper.RestConnection `mapstructure:",squash"`
	SonarqubeAccessToken  `mapstructure:",squash"`
	Organization          string `gorm:"serializer:json" json:"org" mapstructure:"org"`
}

func (connection SonarqubeConn) Sanitize() SonarqubeConn {
	connection.Token = utils.SanitizeString(connection.Token)
	return connection
}

// This object conforms to what the frontend currently sends.
type SonarqubeConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	SonarqubeConn         `mapstructure:",squash"`
}

// This object conforms to what the frontend currently expects.
type SonarqubeResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	SonarqubeConnection
}

func (SonarqubeConnection) TableName() string {
	return "_tool_sonarqube_connections"
}

func (connection SonarqubeConnection) Sanitize() SonarqubeConnection {
	connection.SonarqubeConn = connection.SonarqubeConn.Sanitize()
	return connection
}

func (connection *SonarqubeConnection) MergeFromRequest(target *SonarqubeConnection, body map[string]interface{}) error {
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

func (connection *SonarqubeConnection) IsCloud() bool {
	return connection.Endpoint == "https://sonarcloud.io/api/"
}

const ORG = "org"

func (connection *SonarqubeConn) PrepareApiClient(apiClient plugin.ApiClient) errors.Error {
	apiClient.SetData(ORG, connection.Organization)
	apiClient.SetBeforeFunction(func(req *http.Request) errors.Error {
		org := apiClient.GetData(ORG).(string)
		if org != "" {
			query := req.URL.Query()
			query.Add("organization", org)
			req.URL.RawQuery = query.Encode()
		}
		return nil
	})

	return nil
}

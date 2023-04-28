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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api/apihelperabstract"
)

// PrepareApiClient fetches token from Zentao API for future requests
func (connection ZentaoConn) PrepareApiClient(apiClient apihelperabstract.ApiClientAbstract) errors.Error {
	// request for access token
	tokenReqBody := &ApiAccessTokenRequest{
		Account:  connection.Username,
		Password: connection.Password,
	}
	tokenRes, err := apiClient.Post("/tokens", nil, tokenReqBody, nil)
	if err != nil {
		return err
	}

	if tokenRes.StatusCode == http.StatusUnauthorized {
		return errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error while to request access token")
	}

	tokenResBody := &ApiAccessTokenResponse{}
	err = helper.UnmarshalResponse(tokenRes, tokenResBody)
	if err != nil {
		return err
	}
	if tokenResBody.Token == "" {
		return errors.Default.New("failed to request access token")
	}
	apiClient.SetHeaders(map[string]string{
		"Token": fmt.Sprintf("%v", tokenResBody.Token),
	})
	return nil
}

// ZentaoConn holds the essential information to connect to the Gitlab API
type ZentaoConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.BasicAuth      `mapstructure:",squash"`
}

// ZentaoConnection holds ZentaoConn plus ID/Name for database storage
type ZentaoConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	ZentaoConn            `mapstructure:",squash"`
}

// This object conforms to what the frontend currently expects.
type ZentaoResponse struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
	ZentaoConnection
}

// Using User because it requires authentication.
type ApiUserResponse struct {
	Id   int64
	Name string `json:"name"`
}

func (ZentaoConnection) TableName() string {
	return "_tool_zentao_connections"
}

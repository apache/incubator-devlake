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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api/apihelperabstract"
)

// GitlabConn holds the essential information to connect to the Gitlab API
type GitlabConn struct {
	api.RestConnection `mapstructure:",squash"`
	api.AccessToken    `mapstructure:",squash"`
	IsPrivateToken     bool                `gorm:"-"`
	UserId             int                 `gorm:"-"`
	Version            *ApiVersionResponse `gorm:"-"`
}

// SetupAuthentication sets up the HTTP Request Authentication
func (conn *GitlabConn) SetupAuthentication(req *http.Request) errors.Error {
	if conn.IsPrivateToken {
		req.Header.Set("Private-Token", conn.Token)
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", conn.Token))
	}
	return nil
}

// PrepareApiClient test api and set the IsPrivateToken,version,UserId and so on.
func (conn *GitlabConn) PrepareApiClient(apiClient apihelperabstract.ApiClientAbstract) errors.Error {
	conn.IsPrivateToken = false
	// test request for access token
	userResBody := &ApiUserResponse{}
	res, err := apiClient.Get("user", nil, nil)
	if res.StatusCode != http.StatusUnauthorized {
		if err != nil {
			return errors.Convert(err)
		}
		err = api.UnmarshalResponse(res, userResBody)
		if err != nil {
			return errors.Convert(err)
		}

		if res.StatusCode != http.StatusOK {
			return errors.HttpStatus(res.StatusCode).New("unexpected status code while testing connection")
		}
	} else {
		conn.IsPrivateToken = true
		res, err = apiClient.Get("user", nil, nil)
		if err != nil {
			return errors.Convert(err)
		}
		err = api.UnmarshalResponse(res, userResBody)
		if err != nil {
			return errors.Convert(err)
		}

		if res.StatusCode != http.StatusOK {
			return errors.HttpStatus(res.StatusCode).New("unexpected status code while testing connection[PrivateToken]")
		}
	}
	conn.UserId = userResBody.Id
	// get gitlab version
	versionResBody := &ApiVersionResponse{}
	res, err = apiClient.Get("version", nil, nil)
	if err != nil {
		return errors.Convert(err)
	}

	err = api.UnmarshalResponse(res, versionResBody)
	if err != nil {
		return errors.Convert(err)
	}

	conn.Version = versionResBody

	return nil
}

// GitlabConnection holds GitlabConn plus ID/Name for database storage
type GitlabConnection struct {
	api.BaseConnection `mapstructure:",squash"`
	GitlabConn         `mapstructure:",squash"`
}

// This object conforms to what the frontend currently expects.
type GitlabResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	GitlabConnection
}

type ApiVersionResponse struct {
	Version  string `json:"version"`
	Revision string `json:"revision"`
}

// Using User because it requires authentication.
type ApiUserResponse struct {
	Id   int
	Name string `json:"name"`
}

func (GitlabConnection) TableName() string {
	return "_tool_gitlab_connections"
}

type AeAppKey api.AppKey

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
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api/apihelperabstract"
)

// GitlabConn holds the essential information to connect to the Gitlab API
type GitlabConn struct {
	api.RestConnection `mapstructure:",squash"`
	api.AccessToken    `mapstructure:",squash"`
}

const GitlabApiClientData_UserId string = "UserId"
const GitlabApiClientData_ApiVersion string = "ApiVersion"

// this function is used to rewrite the same function of AccessToken
func (conn *GitlabConn) SetupAuthentication(request *http.Request) errors.Error {
	return nil
}

// PrepareApiClient test api and set the IsPrivateToken,version,UserId and so on.
func (conn *GitlabConn) PrepareApiClient(apiClient apihelperabstract.ApiClientAbstract) errors.Error {
	header1 := http.Header{}
	header1.Set("Authorization", fmt.Sprintf("Bearer %v", conn.Token))
	// test request for access token
	userResBody := &ApiUserResponse{}
	res, err := apiClient.Get("user", nil, header1)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusUnauthorized {
		err = api.UnmarshalResponse(res, userResBody)
		if err != nil {
			return errors.Convert(err)
		}
		if res.StatusCode == http.StatusUnauthorized {
			return errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error while testing connection")
		}
		if res.StatusCode != http.StatusOK {
			return errors.HttpStatus(res.StatusCode).New("unexpected status code while testing connection")
		}
		apiClient.SetHeaders(map[string]string{
			"Authorization": fmt.Sprintf("Bearer %v", conn.Token),
		})
	} else {
		header2 := http.Header{}
		header2.Set("Private-Token", conn.Token)
		res, err = apiClient.Get("user", nil, header2)
		if err != nil {
			return errors.Convert(err)
		}
		err = api.UnmarshalResponse(res, userResBody)
		if err != nil {
			return errors.Convert(err)
		}
		if res.StatusCode == http.StatusUnauthorized {
			return errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error while testing connection[PrivateToken]")
		}
		if res.StatusCode != http.StatusOK {
			return errors.HttpStatus(res.StatusCode).New("unexpected status code while testing connection[PrivateToken]")
		}
		apiClient.SetHeaders(map[string]string{
			"Private-Token": conn.Token,
		})
	}
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

	// add v for semver compare
	if versionResBody.Version[0] != 'v' {
		versionResBody.Version = "v" + versionResBody.Version
	}

	apiClient.SetData(GitlabApiClientData_UserId, userResBody.Id)
	apiClient.SetData(GitlabApiClientData_ApiVersion, versionResBody.Version)

	return nil
}

var _ plugin.ApiConnection = (*GitlabConnection)(nil)

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

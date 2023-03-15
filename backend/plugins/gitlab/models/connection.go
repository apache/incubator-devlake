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
	"context"
	"fmt"
	context2 "github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/plugin"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
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

var _ plugin.ApiConnectionForRemote[GroupResponse, GitlabApiProject] = (*GitlabConnection)(nil)
var _ plugin.ApiGroup = (*GroupResponse)(nil)

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

func (g GitlabConnection) GetGroup(basicRes context2.BasicRes, gid string, query url.Values) ([]GroupResponse, errors.Error) {
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &g)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
	}
	var res *http.Response
	if gid == "" {
		query.Set("top_level_only", "true")
		res, err = apiClient.Get("groups", query, nil)
		if err != nil {
			return nil, err
		}
	} else {
		res, err = apiClient.Get(fmt.Sprintf("groups/%s/subgroups", gid), query, nil)
		if err != nil {
			return nil, err
		}
	}
	var resBody []GroupResponse
	err = api.UnmarshalResponse(res, &resBody)
	if err != nil {
		return nil, err
	}
	return resBody, err
}

func (g GitlabConnection) GetScope(basicRes context2.BasicRes, gid string, query url.Values) ([]GitlabApiProject, errors.Error) {
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &g)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "failed to get create apiClient")
	}
	var res *http.Response
	if gid == "" {
		res, err = apiClient.Get(fmt.Sprintf("users/%d/projects", apiClient.GetData("UserId")), query, nil)
		if err != nil {
			return nil, err
		}
	} else {
		query.Set("with_shared", "false")
		res, err = apiClient.Get(fmt.Sprintf("/groups/%s/projects", gid), query, nil)
		if err != nil {
			return nil, err
		}
	}
	var resBody []GitlabApiProject
	err = api.UnmarshalResponse(res, &resBody)
	if err != nil {
		return nil, err
	}
	return resBody, err
}

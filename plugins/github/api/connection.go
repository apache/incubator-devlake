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

package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/apache/incubator-devlake/plugins/github/models"

	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/mitchellh/mapstructure"

	"github.com/apache/incubator-devlake/plugins/core"
)

/*
POST /plugins/github/test
REQUEST BODY:
{
	"endpoint": "https://api.github.com/",
	"token": "github api access token"
}
RESPONSE BODY:
Success:
{
	"login": "xxx"
}
Failure:
{
	"success": false,
	"message": "invalid token"
}
*/
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// process input
	var params models.TestConnectionRequest
	err := mapstructure.Decode(input.Body, &params)
	if err != nil {
		return nil, err
	}
	err = vld.Struct(params)
	if err != nil {
		return nil, err
	}
	// verify multiple token in parallel
	// PLEASE NOTE: This works because GitHub API Client rotates tokens on each request
	apiClient, err := helper.NewApiClient(
		context.TODO(),
		params.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", params.Token),
		},
		3*time.Second,
		params.Proxy,
	)
	if err != nil {
		return nil, err
	}
	res, err := apiClient.Get("user", nil, nil)
	if err != nil {
		return nil, err
	}
	githubApiResponse := &models.GithubUserOfToken{}
	err = helper.UnmarshalResponse(res, githubApiResponse)
	if err != nil {
		return nil, err
	} else if githubApiResponse.Login == "" {
		return nil, fmt.Errorf("invalid token")
	}

	// output
	return &core.ApiResourceOutput{Body: githubApiResponse, Status: http.StatusOK}, nil
}

/*
POST /plugins/github/connections
*/
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GithubConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
PATCH /plugins/github/connections/:connectionId
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GithubConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
DELETE /plugins/github/connections/:connectionId
*/
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GithubConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &core.ApiResourceOutput{Body: connection}, err
}

/*
GET /plugins/github/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	var connections []models.GithubConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: connections}, nil
}

/*
GET /plugins/github/connections/:connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GithubConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection}, err
}

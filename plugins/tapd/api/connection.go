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
	"fmt"
	"net/http"
	"time"

	"github.com/apache/incubator-devlake/plugins/tapd/models"

	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/mitchellh/mapstructure"

	"github.com/apache/incubator-devlake/plugins/core"
)

/*
POST /plugins/tapd/test
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
	token := params.Auth
	apiClient, err := helper.NewApiClient(
		params.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %s", token),
		},
		3*time.Second,
		params.Proxy,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("verify token failed for %s %w", token, err)
	}
	res, err := apiClient.Get("/quickstart/testauth", nil, nil)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("verify token failed for %s", token)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	// output
	return nil, nil
}

/*
POST /plugins/tapd/connections
{
	"name": "tapd data connections name",
	"endpoint": "tapd api endpoint, i.e. https://merico.atlassian.net/rest",
	"username": "username, usually should be email address",
	"password": "jenkins api access token",
	"rateLimit": 10800,
}
*/
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// create a new connections
	connection := &models.TapdConnection{}

	// update from request and save to database
	//err := refreshAndSaveTapdConnection(tapdConnection, input.Body)
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
PATCH /plugins/tapd/connections/:connectionId
{
	"name": "tapd data connections name",
	"endpoint": "tapd api endpoint, i.e. https://merico.atlassian.net/rest",
	"username": "username, usually should be email address",
	"password": "jenkins api access token",
	"rateLimit": 10800,
}
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.TapdConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: connection}, nil
}

/*
DELETE /plugins/tapd/connections/:connectionId
*/
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.TapdConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &core.ApiResourceOutput{Body: connection}, err
}

/*
GET /plugins/tapd/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	var connections []models.TapdConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

/*
GET /plugins/tapd/connections/:connectionId


{
	"name": "tapd data connections name",
	"endpoint": "tapd api endpoint, i.e. https://merico.atlassian.net/rest",
	"username": "username, usually should be email address",
	"password": "jenkins api access token",
	"rateLimit": 10800,
}
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.TapdConnection{}
	err := connectionHelper.First(connection, input.Params)
	return &core.ApiResourceOutput{Body: connection}, err
}

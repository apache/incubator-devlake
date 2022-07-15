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

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/mitchellh/mapstructure"
)

func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// decode
	var err error
	var connection models.TestConnectionRequest
	err = mapstructure.Decode(input.Body, &connection)
	if err != nil {
		return nil, err
	}
	// validate
	err = vld.Struct(connection)
	if err != nil {
		return nil, err
	}
	// test connection
	apiClient, err := helper.NewApiClient(
		context.TODO(),
		connection.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Bearer %v", connection.Token),
		},
		3*time.Second,
		connection.Proxy,
		basicRes,
	)
	if err != nil {
		return nil, err
	}

	res, err := apiClient.Get("user", nil, nil)
	if err != nil {
		return nil, err
	}
	resBody := &models.ApiUserResponse{}
	err = helper.UnmarshalResponse(res, resBody)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil, nil
}

/*
POST /plugins/gitlab/connections
{
	"name": "gitlab data connection name",
	"endpoint": "gitlab api endpoint, i.e. https://gitlab.com/api/v4/",
	"token": "gitlab api access token"
}
*/
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// update from request and save to database
	connection := &models.GitlabConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
PATCH /plugins/gitlab/connections/:connectionId
{
	"name": "gitlab data connection name",
	"endpoint": "gitlab api endpoint, i.e. https://gitlab.com/api/v4/",
	"token": "gitlab api access token"
}
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GitlabConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection}, nil
}

/*
DELETE /plugins/gitlab/connections/:connectionId
*/
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GitlabConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &core.ApiResourceOutput{Body: connection}, err
}

/*
GET /plugins/gitlab/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	var connections []models.GitlabConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

/*
GET /plugins/gitlab/connections/:connectionId


{
	"name": "gitlab data connection name",
	"endpoint": "gitlab api endpoint, i.e. https://gitlab.com/api/v4/",
	"token": "gitlab api access token"
}
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GitlabConnection{}
	err := connectionHelper.First(connection, input.Params)
	return &core.ApiResourceOutput{Body: connection}, err
}

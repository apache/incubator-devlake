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

	"github.com/apache/incubator-devlake/plugins/jenkins/models"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/utils"
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
	encodedToken := utils.GetEncodedToken(connection.Username, connection.Password)
	apiClient, err := helper.NewApiClient(
		context.TODO(),
		connection.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", encodedToken),
		},
		3*time.Second,
		connection.Proxy,
	)
	if err != nil {
		return nil, err
	}
	res, err := apiClient.Get("", nil, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil, nil
}

/*
POST /plugins/jenkins/connections
{
	"name": "jenkins data connection name",
	"endpoint": "jenkins api endpoint, i.e. https://ci.jenkins.io/",
	"username": "username, usually should be email address",
	"password": "jenkins api access token"
}
*/
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// create a new connection
	connection := &models.JenkinsConnection{}

	// update from request and save to database
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
PATCH /plugins/jenkins/connections/connectionId
{
	"name": "jenkins data connection name",
	"endpoint": "jenkins api endpoint, i.e. https://ci.jenkins.io/",
	"username": "username, usually should be email address",
	"password": "jenkins api access token"
}
*/

func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.JenkinsConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: connection}, nil
}

/*
DELETE /plugins/jenkins/connections/connectionId
*/
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.JenkinsConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &core.ApiResourceOutput{Body: connection}, err
}

/*
GET /plugins/jenkins/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	var connections []models.JenkinsConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

/*
GET /plugins/jenkins/connections/connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.JenkinsConnection{}
	err := connectionHelper.First(connection, input.Params)
	return &core.ApiResourceOutput{Body: connection}, err
}

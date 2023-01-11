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
	"github.com/apache/incubator-devlake/errors"
	"net/http"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
)

func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	// decode
	var err errors.Error
	var connection models.TestConnectionRequest
	if err = helper.Decode(input.Body, &connection, vld); err != nil {
		return nil, err
	}
	// test connection
	apiClient, err := helper.NewApiClient(
		context.TODO(),
		connection.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", connection.GetEncodedToken()),
		},
		3*time.Second,
		connection.Proxy,
		basicRes,
	)
	if err != nil {
		return nil, err
	}

	res, err := apiClient.Get("server/version", nil, nil)
	if err != nil {
		return nil, err
	}
	resBody := &models.ApiUserResponse{}
	err = helper.UnmarshalResponse(res, resBody)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("unexpected status code: %d", res.StatusCode))
	}
	return nil, nil
}

/*
POST /plugins/Sonarqube/connections

	{
		"name": "Sonarqube data connection name",
		"endpoint": "Sonarqube api endpoint, i.e. http://host:port/api/",
		"username": "Actually, we should use sonarqube user token",
		"password": "Actually, we should set this empty"
	}
*/
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	// update from request and save to database
	connection := &models.SonarqubeConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
PATCH /plugins/Sonarqube/connections/:connectionId

	{
		"name": "Sonarqube data connection name",
		"endpoint": "Sonarqube api endpoint, i.e. http://host:port/api/",
		"username": "Actually, we should use sonarqube user token",
		"password": "Actually, we should set this empty"
	}
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connection := &models.SonarqubeConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection}, nil
}

/*
DELETE /plugins/Sonarqube/connections/:connectionId
*/
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connection := &models.SonarqubeConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &core.ApiResourceOutput{Body: connection}, err
}

/*
GET /plugins/Sonarqube/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var connections []models.SonarqubeConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

//TODO Please modify the folowing code to adapt to your plugin
/*
GET /plugins/Sonarqube/connections/:connectionId
{
	"name": "Sonarqube data connection name",
	"endpoint": "Sonarqube api endpoint, i.e. http://host:port/api/",
	"username": "Actually, we should use sonarqube user token",
	"password": "Actually, we should set this empty"
}
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connection := &models.SonarqubeConnection{}
	err := connectionHelper.First(connection, input.Params)
	return &core.ApiResourceOutput{Body: connection}, err
}

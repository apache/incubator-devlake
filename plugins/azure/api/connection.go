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

	"github.com/apache/incubator-devlake/plugins/azure/models"

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
		return nil, errors.BadInput.Wrap(err, "could not decode request parameters", errors.AsUserMessage())
	}
	// validate
	err = vld.Struct(connection)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "could not validate request parameters", errors.AsUserMessage())
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
		basicRes,
	)
	if err != nil {
		return nil, err
	}
	res, err := apiClient.Get("", nil, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New("unexpected status code while testing connection")
	}
	return nil, nil
}

/*
POST /plugins/zaure/connections

	{
		"name": "zaure data connection name",
		"endpoint": "zaure api endpoint, i.e. https://ci.zaure.io/",
		"username": "username, usually should be email address",
		"password": "zaure api access token"
	}
*/
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// create a new connection
	connection := &models.AzureConnection{}

	// update from request and save to database
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
PATCH /plugins/zaure/connections/connectionId
{
	"name": "zaure data connection name",
	"endpoint": "zaure api endpoint, i.e. https://ci.zaure.io/",
	"username": "username, usually should be email address",
	"password": "zaure api access token"
}
*/

func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.AzureConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: connection}, nil
}

/*
DELETE /plugins/zaure/connections/connectionId
*/
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.AzureConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &core.ApiResourceOutput{Body: connection}, err
}

/*
GET /plugins/zaure/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	var connections []models.AzureConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

/*
GET /plugins/zaure/connections/connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.AzureConnection{}
	err := connectionHelper.First(connection, input.Params)
	return &core.ApiResourceOutput{Body: connection}, err
}

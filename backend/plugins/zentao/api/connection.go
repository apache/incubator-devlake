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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// process input
	var params models.TestConnectionRequest
	err := mapstructure.Decode(input.Body, &params)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "could not decode request parameters")
	}
	err = vld.Struct(params)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "could not validate request parameters")
	}

	authApiClient, err := helper.NewApiClient(context.TODO(), params.Endpoint, nil, 0, params.Proxy, basicRes)
	if err != nil {
		return nil, errors.Default.WrapRaw(err)
	}

	// request for access token
	tokenReqBody := &models.ApiAccessTokenRequest{
		Account:  params.Username,
		Password: params.Password,
	}
	tokenRes, err := authApiClient.Post("/tokens", nil, tokenReqBody, nil)
	if err != nil {
		return nil, errors.Default.WrapRaw(err)
	}
	tokenResBody := &models.ApiAccessTokenResponse{}
	err = helper.UnmarshalResponse(tokenRes, tokenResBody)
	if err != nil {
		return nil, errors.Default.WrapRaw(err)
	}
	if tokenResBody.Token == "" {
		return nil, errors.Default.New("failed to request access token")
	}

	// output
	return nil, nil
}

/*
POST /plugins/Zentao/connections

	{
		"name": "Zentao data connection name",
		"endpoint": "Zentao api endpoint, i.e. https://example.com",
		"username": "username, usually should be email address",
		"password": "Zentao api access token"
	}
*/
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// update from request and save to database
	connection := &models.ZentaoConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
PATCH /plugins/Zentao/connections/:connectionId

	{
		"name": "Zentao data connection name",
		"endpoint": "Zentao api endpoint, i.e. https://example.com",
		"username": "username, usually should be email address",
		"password": "Zentao api access token"
	}
*/
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.ZentaoConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection}, nil
}

/*
DELETE /plugins/Zentao/connections/:connectionId
*/
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.ZentaoConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &plugin.ApiResourceOutput{Body: connection}, err
}

/*
GET /plugins/Zentao/connections
*/
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.ZentaoConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

/*
GET /plugins/Zentao/connections/:connectionId

	{
		"name": "Zentao data connection name",
		"endpoint": "Zentao api endpoint, i.e. https://merico.atlassian.net/rest",
		"username": "username, usually should be email address",
		"password": "Zentao api access token"
	}
*/
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.ZentaoConnection{}
	err := connectionHelper.First(connection, input.Params)
	return &plugin.ApiResourceOutput{Body: connection}, err
}

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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
	"net/http"
)

type validation struct {
	Valid bool `json:"valid"`
}

/*
@Summary test sonarqube connection
@Description Test sonarqube Connection
@Tags plugins/sonarqube
@Param body body models.SonarqubeConn true "json body"
@Success 200  {object} shared.ApiBody "Success"
@Failure 400  {string} errcode.Error "Bad Request"
@Failure 500  {string} errcode.Error "Internal Error"
@Router /plugins/sonarqube/test [POST]
*/
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// decode
	var err errors.Error
	var connection models.SonarqubeConn
	if err = api.Decode(input.Body, &connection, vld); err != nil {
		return nil, err
	}
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, connection)
	if err != nil {
		return nil, err
	}

	res, err := apiClient.Get("authentication/validate", nil, nil)
	if err != nil {
		return nil, err
	}
	switch res.StatusCode {
	case 200: // right StatusCode
		body := &validation{}
		err = api.UnmarshalResponse(res, body)
		if err != nil {
			return nil, err
		}
		if !body.Valid {
			return nil, errors.Default.New("Authentication failed, please check your access token.")
		}
		return &plugin.ApiResourceOutput{Body: true, Status: 200}, nil
	case 401: // error secretKey or nonceStr
		return &plugin.ApiResourceOutput{Body: false, Status: res.StatusCode}, nil
	default: // unknow what happen , back to user
		return &plugin.ApiResourceOutput{Body: res.Body, Status: res.StatusCode}, nil
	}
}

/*
POST /plugins/Sonarqube/connections

	{
		"name": "Sonarqube data connection name",
		"endpoint": "Sonarqube api endpoint, i.e. http://host:port/api/",
		"token": "Sonarqube user token"
	}
*/
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// update from request and save to database
	connection := &models.SonarqubeConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
PATCH /plugins/Sonarqube/connections/:connectionId

	{
		"name": "Sonarqube data connection name",
		"endpoint": "Sonarqube api endpoint, i.e. http://host:port/api/",
		"token": "Sonarqube user token"
	}
*/
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.SonarqubeConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection}, nil
}

/*
DELETE /plugins/Sonarqube/connections/:connectionId
*/
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.SonarqubeConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &plugin.ApiResourceOutput{Body: connection}, err
}

/*
GET /plugins/Sonarqube/connections
*/
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.SonarqubeConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

/*
GET /plugins/Sonarqube/connections/:connectionId

	{
		"name": "Sonarqube data connection name",
		"endpoint": "Sonarqube api endpoint, i.e. http://host:port/api/",
		"token": "Sonarqube user token"
	}
*/
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.SonarqubeConnection{}
	err := connectionHelper.First(connection, input.Params)
	return &plugin.ApiResourceOutput{Body: connection}, err
}

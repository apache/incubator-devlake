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
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/opsgenie/models"
)

// @Summary test opsgenie connection
// @Description Test Opsgenie Connection
// @Tags plugins/opsgenie
// @Param body body models.OpsgenieConn true "json body"
// @Success 200  {object} shared.ApiBody "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/opsgenie/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connection models.OpsgenieConn

	err := api.Decode(input.Body, &connection, vld)
	if err != nil {
		return nil, err
	}

	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
	if err != nil {
		return nil, err
	}

	// check API permissions
	response, err := apiClient.Get("v2/heartbeats/HeartbeatName/ping", nil, nil)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusUnauthorized {
		return nil, errors.HttpStatus(http.StatusUnauthorized).New("StatusUnauthorized error when testing api or read_api permissions")
	}

	if response.StatusCode == http.StatusUnprocessableEntity {
		return nil, errors.HttpStatus(http.StatusUnprocessableEntity).New("StatusUnprocessableEntity error when testing api")
	}

	if response.StatusCode == http.StatusForbidden {
		return nil, errors.HttpStatus(http.StatusForbidden).New("API Key need 'Read' and 'Configuration access' Access rights")
	}

	if response.StatusCode == http.StatusOK {
		return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
	}

	return &plugin.ApiResourceOutput{Body: nil, Status: response.StatusCode}, errors.HttpStatus(response.StatusCode).Wrap(err, "could not validate connection")
}

// @Summary create opsgenie connection
// @Description Create Opsgenie connection
// @Tags plugins/opsgenie
// @Param body body models.OpsgenieConnection true "json body"
// @Success 200  {object} models.OpsgenieConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/opsgenie/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.OpsgenieConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

// @Summary patch opsgenie connection
// @Description Patch Opsgenie connection
// @Tags plugins/opsgenie
// @Param body body models.OpsgenieConnection true "json body"
// @Success 200  {object} models.OpsgenieConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/opsgenie/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.OpsgenieConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

// @Summary delete opsgenie connection
// @Description Delete Opsgenie connection
// @Tags plugins/opsgenie
// @Success 200  {object} models.OpsgenieConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} services.BlueprintProjectPairs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/opsgenie/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return connectionHelper.Delete(&models.OpsgenieConnection{}, input)
}

// @Summary list opsgenie connections
// @Description List Opsgenie connections
// @Tags plugins/opsgenie
// @Success 200  {object} models.OpsgenieConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/opsgenie/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.OpsgenieConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}

	return &plugin.ApiResourceOutput{Body: connections}, nil
}

// @Summary get opsgenie connection
// @Description Get Opsgenie connection
// @Tags plugins/opsgenie
// @Success 200  {object} models.OpsgenieConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/opsgenie/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.OpsgenieConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection}, nil
}

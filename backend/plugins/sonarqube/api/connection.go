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
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
	"github.com/apache/incubator-devlake/server/api/shared"
)

type validation struct {
	Valid bool `json:"valid"`
}
type SonarqubeTestConnResponse struct {
	shared.ApiBody
	Connection *models.SonarqubeConn
}

// TestConnection test sonarqube connection options
// @Summary test sonarqube connection
// @Description Test sonarqube Connection
// @Tags plugins/sonarqube
// @Param body body models.SonarqubeConn true "json body"
// @Success 200  {object} SonarqubeTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/sonarqube/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// decode
	var err errors.Error
	var connection models.SonarqubeConn
	if err = api.Decode(input.Body, &connection, vld); err != nil {
		return nil, err
	}
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
	if err != nil {
		return nil, err
	}

	res, err := apiClient.Get("authentication/validate", nil, nil)
	if err != nil {
		return nil, err
	}
	switch res.StatusCode {
	case 200: // right StatusCode
		valid := &validation{}
		err = api.UnmarshalResponse(res, valid)
		body := SonarqubeTestConnResponse{}
		body.Success = true
		body.Message = "success"
		body.Connection = &connection
		if err != nil {
			return nil, err
		}
		if !valid.Valid {
			return nil, errors.Default.New("Authentication failed, please check your access token.")
		}
		return &plugin.ApiResourceOutput{Body: body, Status: 200}, nil
	case 401: // error secretKey or nonceStr
		return &plugin.ApiResourceOutput{Body: false, Status: http.StatusBadRequest}, nil
	default: // unknow what happen , back to user
		return &plugin.ApiResourceOutput{Body: res.Body, Status: res.StatusCode}, nil
	}
}

// PostConnections create sonarqube connection
// @Summary create sonarqube connection
// @Description Create sonarqube connection
// @Tags plugins/sonarqube
// @Param body body models.SonarqubeConnection true "json body"
// @Success 200  {object} models.SonarqubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/sonarqube/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// update from request and save to database
	connection := &models.SonarqubeConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

// PatchConnection patch sonarqube connection
// @Summary patch sonarqube connection
// @Description Patch sonarqube connection
// @Tags plugins/sonarqube
// @Param body body models.SonarqubeConnection true "json body"
// @Param connectionId path int false "connection ID"
// @Success 200  {object} models.SonarqubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.SonarqubeConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection}, nil
}

// DeleteConnection delete a sonarqube connection
// @Summary delete a sonarqube connection
// @Description Delete a sonarqube connection
// @Tags plugins/sonarqube
// @Param connectionId path int false "connection ID"
// @Success 200  {object} models.SonarqubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.SonarqubeConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &plugin.ApiResourceOutput{Body: connection}, err
}

// ListConnections get all sonarqube connections
// @Summary get all sonarqube connections
// @Description Get all sonarqube connections
// @Tags plugins/sonarqube
// @Success 200  {object} []models.SonarqubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/sonarqube/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.SonarqubeConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

// GetConnection get sonarqube connection detail
// @Summary get sonarqube connection detail
// @Description Get sonarqube connection detail
// @Tags plugins/sonarqube
// @Param connectionId path int false "connection ID"
// @Success 200  {object} models.SonarqubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.SonarqubeConnection{}
	err := connectionHelper.First(connection, input.Params)
	return &plugin.ApiResourceOutput{Body: connection}, err
}

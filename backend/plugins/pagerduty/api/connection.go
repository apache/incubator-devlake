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
	"github.com/apache/incubator-devlake/plugins/pagerduty/models"
)

func testConnection(ctx context.Context, connection models.PagerDutyConn) (*plugin.ApiResourceOutput, errors.Error) {
	if vld != nil {
		if err := vld.Struct(connection); err != nil {
			return nil, errors.Default.Wrap(err, "error validating target")
		}
	}
	apiClient, err := api.NewApiClientFromConnection(ctx, basicRes, &connection)
	if err != nil {
		return nil, err
	}
	response, err := apiClient.Get("licenses", nil, nil)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == http.StatusUnauthorized {
		return nil, errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error while testing connection")
	}
	if response.StatusCode == http.StatusOK {
		return &plugin.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
	}
	return &plugin.ApiResourceOutput{Body: nil, Status: response.StatusCode}, errors.HttpStatus(response.StatusCode).Wrap(err, "could not validate connection")
}

// TestConnection test pagerduty connection
// @Summary test pagerduty connection
// @Description Test Pagerduty Connection
// @Tags plugins/pagerduty
// @Param body body models.PagerDutyConn true "json body"
// @Success 200  {object} shared.ApiBody "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/pagerduty/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connection models.PagerDutyConn
	err := api.Decode(input.Body, &connection, vld)
	if err != nil {
		return nil, err
	}
	testConnectionResult, testConnectionErr := testConnection(context.TODO(), connection)
	if testConnectionErr != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, testConnectionErr)
	}
	return testConnectionResult, nil
}

// TestExistingConnection test pagerduty connection
// @Summary test pagerduty connection
// @Description Test Pagerduty Connection
// @Tags plugins/pagerduty
// @Success 200  {object} shared.ApiBody "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/pagerduty/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.PagerDutyConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "find connection from db")
	}
	if err := api.DecodeMapStruct(input.Body, connection, false); err != nil {
		return nil, err
	}
	testConnectionResult, testConnectionErr := testConnection(context.TODO(), connection.PagerDutyConn)
	if testConnectionErr != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, testConnectionErr)
	}
	return testConnectionResult, nil
}

// @Summary create pagerduty connection
// @Description Create Pagerduty connection
// @Tags plugins/pagerduty
// @Param body body models.PagerDutyConnection true "json body"
// @Success 200  {object} models.PagerDutyConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/pagerduty/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.PagerDutyConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection.Sanitize(), Status: http.StatusOK}, nil
}

// @Summary patch pagerduty connection
// @Description Patch Pagerduty connection
// @Tags plugins/pagerduty
// @Param body body models.PagerDutyConnection true "json body"
// @Success 200  {object} models.PagerDutyConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/pagerduty/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.PagerDutyConnection{}
	if err := connectionHelper.First(&connection, input.Params); err != nil {
		return nil, err
	}
	if err := (&models.PagerDutyConnection{}).MergeFromRequest(connection, input.Body); err != nil {
		return nil, errors.Convert(err)
	}
	if err := connectionHelper.SaveWithCreateOrUpdate(connection); err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection.Sanitize(), Status: http.StatusOK}, nil
}

// @Summary delete pagerduty connection
// @Description Delete Pagerduty connection
// @Tags plugins/pagerduty
// @Success 200  {object} models.PagerDutyConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} services.BlueprintProjectPairs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/pagerduty/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	conn := &models.PagerDutyConnection{}
	output, err := connectionHelper.Delete(conn, input)
	if err != nil {
		return output, err
	}
	output.Body = conn.Sanitize()
	return output, nil
}

// @Summary list pagerduty connections
// @Description List Pagerduty connections
// @Tags plugins/pagerduty
// @Success 200  {object} models.PagerDutyConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/pagerduty/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.PagerDutyConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	for idx, c := range connections {
		connections[idx] = c.Sanitize()
	}
	return &plugin.ApiResourceOutput{Body: connections}, nil
}

// @Summary get pagerduty connection
// @Description Get Pagerduty connection
// @Tags plugins/pagerduty
// @Success 200  {object} models.PagerDutyConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/pagerduty/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.PagerDutyConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection.Sanitize()}, nil
}

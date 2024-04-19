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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/teambition/models"
	"github.com/apache/incubator-devlake/plugins/teambition/tasks"
	"github.com/apache/incubator-devlake/server/api/shared"
)

type TeambitionTestConnResponse struct {
	shared.ApiBody
	Connection *models.TeambitionConn
}

func testConnection(ctx context.Context, connection models.TeambitionConn) (*TeambitionTestConnResponse, errors.Error) {
	// process input
	if vld != nil {
		if err := vld.Struct(connection); err != nil {
			return nil, errors.Default.Wrap(err, "error validating target")
		}
	}
	// test connection
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
	if err != nil {
		return nil, err
	}

	res, err := apiClient.Get("/org/info?orgId="+connection.TenantId, nil, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return nil, errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error while testing connection")
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("unexpected status code: %d", res.StatusCode))
	}
	resBody := tasks.TeambitionComRes[any]{}
	err = api.UnmarshalResponse(res, &resBody)
	if err != nil {
		return nil, err
	}
	if resBody.Code != http.StatusOK {
		return nil, errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized on body while testing connection")
	}
	if resBody.Code != http.StatusOK {
		return nil, errors.HttpStatus(resBody.Code).New(fmt.Sprintf("unexpected body status code: %d", resBody.Code))
	}

	connection = connection.Sanitize()
	body := TeambitionTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = &connection
	// output
	return &body, nil
}

// TestConnection @Summary test teambition connection
// @Description Test teambition Connection
// @Tags plugins/teambition
// @Param body body models.TeambitionConn true "json body"
// @Success 200  {object} TeambitionTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/teambition/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// process input
	var connection models.TeambitionConn
	err := api.Decode(input.Body, &connection, vld)
	if err != nil {
		return nil, err
	}

	// test connection
	result, err := testConnection(context.TODO(), connection)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// TestExistingConnection test teambition connection options
// @Summary test teambition connection
// @Description Test teambition Connection
// @Tags plugins/teambition
// @Param connectionId path int true "connection ID"
// @Success 200  {object} TeambitionTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/teambition/connections/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection, err := dsHelper.ConnApi.GetMergedConnection(input)
	if err != nil {
		return nil, err
	}
	testConnectionResult, testConnectionErr := testConnection(context.TODO(), connection.TeambitionConn)
	if testConnectionErr != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, testConnectionErr)
	}
	return &plugin.ApiResourceOutput{Body: testConnectionResult, Status: http.StatusOK}, nil
}

// PostConnections @Summary create teambition connection
// @Description Create teambition connection
// @Tags plugins/teambition
// @Param body body models.TeambitionConnection true "json body"
// @Success 200  {object} models.TeambitionConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/teambition/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Post(input)
}

// PatchConnection @Summary patch teambition connection
// @Description Patch teambition connection
// @Tags plugins/teambition
// @Param body body models.TeambitionConnection true "json body"
// @Success 200  {object} models.TeambitionConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/teambition/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Patch(input)
}

// DeleteConnection @Summary delete a teambition connection
// @Description Delete a teambition connection
// @Tags plugins/teambition
// @Success 200  {object} models.TeambitionConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} services.BlueprintProjectPairs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/teambition/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Delete(input)
}

// ListConnections @Summary get all teambition connections
// @Description Get all teambition connections
// @Tags plugins/teambition
// @Success 200  {object} []models.TeambitionConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/teambition/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetAll(input)
}

// GetConnection @Summary get teambition connection detail
// @Description Get teambition connection detail
// @Tags plugins/teambition
// @Success 200  {object} models.TeambitionConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/teambition/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetDetail(input)
}

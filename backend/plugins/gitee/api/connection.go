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

	"github.com/apache/incubator-devlake/server/api/shared"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitee/models"
)

type GiteeTestConnResponse struct {
	shared.ApiBody
	Connection *models.GiteeConn
}

func testConnection(ctx context.Context, connection models.GiteeConn) (*GiteeTestConnResponse, errors.Error) {
	// validate
	if vld != nil {
		if err := vld.Struct(connection); err != nil {
			return nil, errors.Default.Wrap(err, "error validating target")
		}
	}
	apiClient, err := helper.NewApiClientFromConnection(ctx, basicRes, &connection)
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

	if res.StatusCode == http.StatusUnauthorized {
		return nil, errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error when testing connection")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New("unexpected status code when testing connection")
	}
	connection = connection.Sanitize()
	body := GiteeTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = &connection
	// output
	return &body, nil
}

// TestConnection test gitee connection
// @Summary test gitee connection
// @Description Test gitee Connection. endpoint: https://gitee.com/api/v5/
// @Tags plugins/gitee
// @Param body body models.GiteeConn true "json body"
// @Success 200  {object} GiteeTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitee/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var err errors.Error
	var connection models.GiteeConn
	if err = helper.Decode(input.Body, &connection, vld); err != nil {
		return nil, errors.BadInput.Wrap(err, "could not decode request parameters")
	}
	// test connection
	result, err := testConnection(context.TODO(), connection)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// TestExistingConnection test gitee connection
// @Summary test gitee connection
// @Description Test gitee Connection. endpoint: https://gitee.com/api/v5/
// @Tags plugins/gitee
// @Param connectionId path int true "connection ID"
// @Success 200  {object} GiteeTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitee/connections/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.GiteeConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "find connection from db")
	}
	if err := helper.DecodeMapStruct(input.Body, connection, false); err != nil {
		return nil, err
	}
	// test connection
	result, err := testConnection(context.TODO(), connection.GiteeConn)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// @Summary create gitee connection
// @Description Create gitee connection
// @Tags plugins/gitee
// @Param body body models.GiteeConnection true "json body"
// @Success 200  {object} models.GiteeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitee/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.GiteeConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection.Sanitize(), Status: http.StatusOK}, nil
}

// @Summary patch gitee connection
// @Description Patch gitee connection
// @Tags plugins/gitee
// @Param body body models.GiteeConnection true "json body"
// @Success 200  {object} models.GiteeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitee/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.GiteeConnection{}
	if err := connectionHelper.First(&connection, input.Params); err != nil {
		return nil, err
	}
	if err := (&models.GiteeConnection{}).MergeFromRequest(connection, input.Body); err != nil {
		return nil, errors.Convert(err)
	}
	if err := connectionHelper.SaveWithCreateOrUpdate(connection); err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection.Sanitize(), Status: http.StatusOK}, nil
}

// @Summary delete a gitee connection
// @Description Delete a gitee connection
// @Tags plugins/gitee
// @Success 200  {object} models.GiteeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} services.BlueprintProjectPairs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitee/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	conn := &models.GiteeConnection{}
	output, err := connectionHelper.Delete(conn, input)
	if err != nil {
		return output, err
	}
	output.Body = conn.Sanitize()
	return output, nil

}

// @Summary get all gitee connections
// @Description Get all gitee connections
// @Tags plugins/gitee
// @Success 200  {object} models.GiteeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitee/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.GiteeConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	for idx, c := range connections {
		connections[idx] = c.Sanitize()
	}
	return &plugin.ApiResourceOutput{Body: connections}, nil
}

// @Summary get gitee connection detail
// @Description Get gitee connection detail
// @Tags plugins/gitee
// @Success 200  {object} models.GiteeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitee/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.GiteeConnection{}
	err := connectionHelper.First(connection, input.Params)
	return &plugin.ApiResourceOutput{Body: connection.Sanitize()}, err
}

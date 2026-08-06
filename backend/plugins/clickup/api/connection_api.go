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
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/clickup/models"
	"github.com/apache/incubator-devlake/plugins/clickup/tasks"
	"github.com/apache/incubator-devlake/server/api/shared"
)

type ClickUpTestConnResponse struct {
	shared.ApiBody
	Connection *models.ClickUpConn
}

func testConnection(ctx context.Context, connection models.ClickUpConn) (*ClickUpTestConnResponse, errors.Error) {
	if vld != nil {
		if err := vld.Struct(connection); err != nil {
			return nil, errors.Default.Wrap(err, "error validating target")
		}
	}
	if connection.Endpoint == "" {
		connection.Endpoint = tasks.DefaultEndpoint
	}
	apiClient, err := helper.NewApiClientFromConnection(ctx, basicRes, &connection)
	if err != nil {
		return nil, err
	}
	// GET /team lists the authenticated user's workspaces; it verifies the token.
	res, err := apiClient.Get("team", nil, nil)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "verify token failed")
	}
	if res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusForbidden {
		return nil, errors.HttpStatus(http.StatusBadRequest).New("authentication failed, please check your API token")
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New("unexpected status code while testing connection")
	}
	connection = connection.Sanitize()
	body := ClickUpTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = &connection
	return &body, nil
}

// TestConnection test clickup connection
// @Summary test clickup connection
// @Description Test clickup Connection
// @Tags plugins/clickup
// @Param body body models.ClickUpConn true "json body"
// @Success 200  {object} ClickUpTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/clickup/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connection models.ClickUpConn
	if err := helper.Decode(input.Body, &connection, vld); err != nil {
		return nil, err
	}
	result, err := testConnection(context.TODO(), connection)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// TestExistingConnection test clickup connection by ID
// @Summary test clickup connection
// @Description Test clickup Connection
// @Tags plugins/clickup
// @Param connectionId path int true "connection ID"
// @Success 200  {object} ClickUpTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/clickup/connections/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection, err := dsHelper.ConnApi.GetMergedConnection(input)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "find connection from db")
	}
	if err := helper.DecodeMapStruct(input.Body, connection, false); err != nil {
		return nil, err
	}
	result, testErr := testConnection(context.TODO(), connection.ClickUpConn)
	if testErr != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, testErr)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// PostConnections create clickup connection
// @Summary create clickup connection
// @Description Create clickup connection
// @Tags plugins/clickup
// @Param body body models.ClickUpConnection true "json body"
// @Success 200  {object} models.ClickUpConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/clickup/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Post(input)
}

// PatchConnection patch clickup connection
// @Summary patch clickup connection
// @Description Patch clickup connection
// @Tags plugins/clickup
// @Param body body models.ClickUpConnection true "json body"
// @Success 200  {object} models.ClickUpConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/clickup/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Patch(input)
}

// DeleteConnection delete a clickup connection
// @Summary delete a clickup connection
// @Description Delete a clickup connection
// @Tags plugins/clickup
// @Success 200  {object} models.ClickUpConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} services.BlueprintProjectPairs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/clickup/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Delete(input)
}

// ListConnections get all clickup connections
// @Summary get all clickup connections
// @Description Get all clickup connections
// @Tags plugins/clickup
// @Success 200  {object} []models.ClickUpConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/clickup/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetAll(input)
}

// GetConnection get clickup connection detail
// @Summary get clickup connection detail
// @Description Get clickup connection detail
// @Tags plugins/clickup
// @Success 200  {object} models.ClickUpConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/clickup/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetDetail(input)
}

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
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/slack/models"
)

type SlackTestConnResponse struct {
	shared.ApiBody
	Connection *models.SlackConn
}

func testConnection(ctx context.Context, connection models.SlackConn) (*SlackTestConnResponse, errors.Error) {
	// validate
	if vld != nil {
		if err := vld.Struct(connection); err != nil {
			return nil, errors.Default.Wrap(err, "error validating target")
		}
	}
	// test connection
	_, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
	if err != nil {
		return nil, err
	}
	connection = connection.Sanitize()
	body := SlackTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = &connection
	return &body, nil
}

// TestConnection test slack connection
// @Summary test slack connection
// @Description Test slack Connection. endpoint: https://open.slack.cn/open-apis/
// @Tags plugins/slack
// @Param body body models.SlackConn true "json body"
// @Success 200  {object} SlackTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/slack/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// process input
	var connection models.SlackConn
	if err := api.Decode(input.Body, &connection, vld); err != nil {
		return nil, errors.BadInput.Wrap(err, "could not decode request parameters")
	}
	// test connection
	result, err := testConnection(context.TODO(), connection)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// TestExistingConnection test slack connection
// @Summary test slack connection
// @Description Test slack Connection. endpoint: https://open.slack.cn/open-apis/
// @Tags plugins/slack
// @Param connectionId path int true "connection ID"
// @Success 200  {object} SlackTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/slack/connections/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.SlackConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "find connection from db")
	}
	if err := api.DecodeMapStruct(input.Body, connection, false); err != nil {
		return nil, err
	}
	// test connection
	result, err := testConnection(context.TODO(), connection.SlackConn)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// @Summary create slack connection
// @Description Create slack connection
// @Tags plugins/slack
// @Param body body models.SlackConnection true "json body"
// @Success 200  {object} models.SlackConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/slack/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.SlackConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection.Sanitize(), Status: http.StatusOK}, nil
}

// @Summary patch slack connection
// @Description Patch slack connection
// @Tags plugins/slack
// @Param body body models.SlackConnection true "json body"
// @Success 200  {object} models.SlackConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/slack/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.SlackConnection{}
	if err := connectionHelper.First(&connection, input.Params); err != nil {
		return nil, err
	}
	if err := (&models.SlackConnection{}).MergeFromRequest(connection, input.Body); err != nil {
		return nil, errors.Convert(err)
	}
	if err := connectionHelper.SaveWithCreateOrUpdate(connection); err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection.Sanitize(), Status: http.StatusOK}, nil
}

// @Summary delete a slack connection
// @Description Delete a slack connection
// @Tags plugins/slack
// @Success 200  {object} models.SlackConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} services.BlueprintProjectPairs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/slack/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	conn := &models.SlackConnection{}
	output, err := connectionHelper.Delete(conn, input)
	if err != nil {
		return output, err
	}
	output.Body = conn.Sanitize()
	return output, nil
}

// @Summary get all slack connections
// @Description Get all slack connections
// @Tags plugins/slack
// @Success 200  {object} models.SlackConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/slack/connections [GET]
func ListConnections(_ *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.SlackConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	for idx, c := range connections {
		connections[idx] = c.Sanitize()
	}
	return &plugin.ApiResourceOutput{Body: connections}, nil
}

// @Summary get slack connection detail
// @Description Get slack connection detail
// @Tags plugins/slack
// @Success 200  {object} models.SlackConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/slack/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.SlackConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection.Sanitize()}, err
}

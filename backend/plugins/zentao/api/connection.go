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
	"github.com/apache/incubator-devlake/core/runner"
	"net/http"

	"github.com/apache/incubator-devlake/server/api/shared"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/zentao/models"
)

type ZentaoTestConnResponse struct {
	shared.ApiBody
	Connection *models.ZentaoConn
}

// @Summary test zentao connection
// @Description Test zentao Connection
// @Tags plugins/zentao
// @Param body body models.ZentaoConn true "json body"
// @Success 200  {object} ZentaoTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/zentao/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// process input
	var connection models.ZentaoConn
	err := helper.Decode(input.Body, &connection, vld)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "failed to decode input to be zentao connection")
	}

	// try to create apiClient
	client, err := helper.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
	if err != nil {
		return nil, err
	}
	resp, err := client.Get("/user", nil, nil)
	if err != nil {
		return nil, err
	}
	var body ZentaoTestConnResponse
	if resp.StatusCode != http.StatusOK {
		body.Success = false
		body.Message = err.Error()
		return &plugin.ApiResourceOutput{Body: body, Status: http.StatusBadRequest}, nil
	}
	if connection.DbUrl != "" {
		err = runner.CheckDbConnection(connection.DbUrl)
		if err != nil {
			body.Success = false
			body.Message = "invalid DbUrl"
			return &plugin.ApiResourceOutput{Body: body, Status: http.StatusBadRequest}, nil
		}
	}
	body.Success = true
	body.Message = "success"
	body.Connection = &connection
	// output
	return &plugin.ApiResourceOutput{Body: body, Status: 200}, nil
}

// @Summary create zentao connection
// @Description Create zentao connection
// @Tags plugins/zentao
// @Param body body models.ZentaoConnection true "json body"
// @Success 200  {object} models.ZentaoConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/zentao/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// update from request and save to database
	connection := &models.ZentaoConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

// @Summary patch zentao connection
// @Description Patch zentao connection
// @Tags plugins/zentao
// @Param body body models.ZentaoConnection true "json body"
// @Success 200  {object} models.ZentaoConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/zentao/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.ZentaoConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection}, nil
}

// @Summary delete a zentao connection
// @Description Delete a zentao connection
// @Tags plugins/zentao
// @Success 200  {object} models.ZentaoConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} services.BlueprintProjectPairs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/zentao/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return connectionHelper.Delete(&models.ZentaoConnection{}, input)
}

// @Summary get all zentao connections
// @Description Get all zentao connections
// @Tags plugins/zentao
// @Success 200  {object} []models.ZentaoConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/zentao/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.ZentaoConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

// @Summary get zentao connection detail
// @Description Get zentao connection detail
// @Tags plugins/zentao
// @Success 200  {object} models.ZentaoConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/zentao/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.ZentaoConnection{}
	err := connectionHelper.First(connection, input.Params)
	return &plugin.ApiResourceOutput{Body: connection}, err
}

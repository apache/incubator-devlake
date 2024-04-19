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
	"time"

	"github.com/apache/incubator-devlake/core/runner"

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

func testConnection(ctx context.Context, connection models.ZentaoConn) (*ZentaoTestConnResponse, errors.Error) {
	// process input
	if vld != nil {
		if err := vld.Struct(connection); err != nil {
			return nil, errors.Default.Wrap(err, "error validating target")
		}
	}
	// try to create apiClient
	client, err := helper.NewApiClientFromConnection(ctx, basicRes, &connection)
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
		return &body, nil
	}
	if connection.DbUrl != "" {
		err = runner.CheckDbConnection(connection.DbUrl, 5*time.Second)
		if err != nil {
			body.Success = false
			body.Message = "invalid DbUrl"
			return &body, errors.Default.New("invalid DbUrl")
		}
	}
	body.Success = true
	body.Message = "success"
	connection = connection.Sanitize()
	body.Connection = &connection
	return &body, nil
}

// TestConnection test zentao connection
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
	// test connection
	result, err := testConnection(context.TODO(), connection)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// TestExistingConnection test zentao connection options
// @Summary test zentao connection
// @Description Test zentao Connection
// @Tags plugins/zentao
// @Param connectionId path int true "connection ID"
// @Success 200  {object} ZentaoTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/zentao/connections/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection, err := dsHelper.ConnApi.GetMergedConnection(input)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "find connection from db")
	}
	if err := helper.DecodeMapStruct(input.Body, connection, false); err != nil {
		return nil, err
	}
	testConnectionResult, testConnectionErr := testConnection(context.TODO(), connection.ZentaoConn)
	if testConnectionErr != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, testConnectionErr)
	}
	return &plugin.ApiResourceOutput{Body: testConnectionResult, Status: http.StatusOK}, nil
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
	return dsHelper.ConnApi.Post(input)
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
	return dsHelper.ConnApi.Patch(input)
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
	return dsHelper.ConnApi.Delete(input)
}

// @Summary get all zentao connections
// @Description Get all zentao connections
// @Tags plugins/zentao
// @Success 200  {object} []models.ZentaoConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/zentao/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetAll(input)
}

// @Summary get zentao connection detail
// @Description Get zentao connection detail
// @Tags plugins/zentao
// @Success 200  {object} models.ZentaoConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/zentao/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetDetail(input)
}

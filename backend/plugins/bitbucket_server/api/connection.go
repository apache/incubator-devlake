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
	plugin "github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bitbucket_server/models"
)

type BitBucketServerTestConnResponse struct {
	shared.ApiBody
	Connection *models.BitbucketServerConn
}

func testConnection(ctx context.Context, connection models.BitbucketServerConn) (*BitBucketServerTestConnResponse, errors.Error) {
	// test connection
	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
	if err != nil {
		return nil, err
	}
	res, err := apiClient.Get("rest/api/1.0/projects", nil, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return nil, errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error when testing connection")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New("unexpected status code when testing connection")
	}
	body := BitBucketServerTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = &connection
	// output
	return &body, nil
}

// @Summary test bitbucket connection
// @Description Test bitbucket Connection
// @Tags plugins/bitbucket_server
// @Param body body models.BitbucketServerConn true "json body"
// @Success 200  {object} BitBucketServerTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/bitbucket_server/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// decode
	var err errors.Error
	var connection models.BitbucketServerConn
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

// TestExistingConnection test bitbucket connection
// @Summary test bitbucket connection
// @Description Test bitbucket Connection
// @Tags plugins/bitbucket
// @Param connectionId path int true "connection ID"
// @Success 200  {object} BitBucketTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/bitbucket_server/connections/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection, err := dsHelper.ConnApi.FindByPk(input)
	if err != nil {
		return nil, err
	}
	// test connection
	result, err := testConnection(context.TODO(), connection.BitbucketServerConn)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// @Summary create bitbucket connection
// @Description Create bitbucket connection
// @Tags plugins/bitbucket_server
// @Param body body models.BitbucketServerConnection true "json body"
// @Success 200  {object} models.BitbucketServerConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/bitbucket_server/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Post(input)
}

// @Summary patch bitbucket connection
// @Description Patch bitbucket connection
// @Tags plugins/bitbucket_server
// @Param body body models.BitbucketServerConnection true "json body"
// @Success 200  {object} models.BitbucketServerConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/bitbucket_server/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Patch(input)
}

// @Summary delete a bitbucket connection
// @Description Delete a bitbucket connection
// @Tags plugins/bitbucket_server
// @Success 200  {object} models.BitbucketServerConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} services.BlueprintProjectPairs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/bitbucket_server/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Delete(input)
}

// @Summary get all bitbucket connections
// @Description Get all bitbucket connections
// @Tags plugins/bitbucket_server
// @Success 200  {object} []models.BitbucketServerConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/bitbucket_server/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetAll(input)
}

// @Summary get bitbucket connection detail
// @Description Get bitbucket connection detail
// @Tags plugins/bitbucket_server
// @Success 200  {object} models.BitbucketServerConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/bitbucket_server/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetDetail(input)
}

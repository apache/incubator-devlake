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
	"github.com/apache/incubator-devlake/plugins/bitbucket/models"
)

type BitBucketTestConnResponse struct {
	shared.ApiBody
	Connection *models.BitbucketConn
}

func testConnection(ctx context.Context, connection models.BitbucketConn) (*BitBucketTestConnResponse, errors.Error) {
	// validate
	if vld != nil {
		if err := vld.Struct(connection); err != nil {
			return nil, errors.Default.Wrap(err, "error validating target")
		}
	}
	// test connection
	apiClient, err := api.NewApiClientFromConnection(ctx, basicRes, &connection)
	if err != nil {
		return nil, err
	}
	res, err := apiClient.Get("user", nil, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return nil, errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error when testing connection. Please check your credentials.")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New("unexpected status code when testing connection")
	}

	// Log deprecation warning if using App Password (not API token)
	if !connection.UsesApiToken {
		basicRes.GetLogger().Warn(nil, "Bitbucket App passwords are deprecated and will be deactivated on June 9, 2026. Please migrate to API tokens.")
	}

	connection = connection.Sanitize()
	body := BitBucketTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = &connection
	// output
	return &body, nil
}

// TestConnection test bitbucket connection
// @Summary test bitbucket connection
// @Description Test bitbucket Connection
// @Tags plugins/bitbucket
// @Param body body models.BitbucketConn true "json body"
// @Success 200  {object} BitBucketTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/bitbucket/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// decode
	var err errors.Error
	var connection models.BitbucketConn
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
// @Router /plugins/bitbucket/connections/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection, err := dsHelper.ConnApi.GetMergedConnection(input)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "find connection from db")
	}
	if err := api.DecodeMapStruct(input.Body, connection, false); err != nil {
		return nil, err
	}
	// test connection
	result, err := testConnection(context.TODO(), connection.BitbucketConn)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// @Summary create bitbucket connection
// @Description Create bitbucket connection
// @Tags plugins/bitbucket
// @Param body body models.BitbucketConnection true "json body"
// @Success 200  {object} models.BitbucketConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/bitbucket/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Post(input)
}

// @Summary patch bitbucket connection
// @Description Patch bitbucket connection
// @Tags plugins/bitbucket
// @Param body body models.BitbucketConnection true "json body"
// @Success 200  {object} models.BitbucketConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/bitbucket/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Patch(input)
}

// @Summary delete a bitbucket connection
// @Description Delete a bitbucket connection
// @Tags plugins/bitbucket
// @Success 200  {object} models.BitbucketConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} services.BlueprintProjectPairs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/bitbucket/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Delete(input)
}

// @Summary get all bitbucket connections
// @Description Get all bitbucket connections
// @Tags plugins/bitbucket
// @Success 200  {object} []models.BitbucketConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/bitbucket/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetAll(input)
}

// @Summary get bitbucket connection detail
// @Description Get bitbucket connection detail
// @Tags plugins/bitbucket
// @Success 200  {object} models.BitbucketConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/bitbucket/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetDetail(input)
}

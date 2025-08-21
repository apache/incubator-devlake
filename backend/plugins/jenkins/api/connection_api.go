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
	"strings"

	"github.com/apache/incubator-devlake/server/api/shared"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
)

type JenkinsTestConnResponse struct {
	shared.ApiBody
	Connection *models.JenkinsConn
}

func testConnection(ctx context.Context, connection models.JenkinsConn) (*JenkinsTestConnResponse, errors.Error) {
	// validate
	if vld != nil {
		if err := vld.Struct(connection); err != nil {
			return nil, errors.Default.Wrap(err, "error validating target")
		}
	}
	// Check if the URL contains "/api"
	if strings.Contains(connection.Endpoint, "/api") {
		return nil, errors.HttpStatus(http.StatusBadRequest).New("Invalid URL. Please use the base URL without /api")
	}
	// test connection
	apiClient, err := api.NewApiClientFromConnection(ctx, basicRes, &connection)
	if err != nil {
		return nil, err
	}
	res, err := apiClient.Get("", nil, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return nil, errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error while testing connection")
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New("unexpected status code when testing connection")
	}
	connection = connection.Sanitize()
	body := JenkinsTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = &connection
	// output
	return &body, nil
}

// TestConnection test jenkins connection
// @Summary test jenkins connection
// @Description Test Jenkins Connection
// @Tags plugins/jenkins
// @Param body body models.JenkinsConn true "json body"
// @Success 200  {object} JenkinsTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/jenkins/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// decode
	var err errors.Error
	var connection models.JenkinsConn
	err = api.Decode(input.Body, &connection, vld)
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

// TestExistingConnection test jenkins connection
// @Summary test jenkins connection
// @Description Test Jenkins Connection
// @Tags plugins/jenkins
// @Param connectionId path int true "connection ID"
// @Success 200  {object} JenkinsTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/jenkins/connections/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection, err := dsHelper.ConnApi.GetMergedConnection(input)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "get merged connection")
	}
	// test connection
	if result, err := testConnection(context.TODO(), connection.JenkinsConn); err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	} else {
		return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
	}
}

// @Summary create jenkins connection
// @Description Create Jenkins connection
// @Tags plugins/jenkins
// @Param body body models.JenkinsConnection true "json body"
// @Success 200  {object} models.JenkinsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/jenkins/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Post(input)
}

// @Summary patch jenkins connection
// @Description Patch Jenkins connection
// @Tags plugins/jenkins
// @Param body body models.JenkinsConnection true "json body"
// @Success 200  {object} models.JenkinsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/jenkins/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Patch(input)
}

// @Summary delete a jenkins connection
// @Description Delete a Jenkins connection
// @Tags plugins/jenkins
// @Success 200  {object} models.JenkinsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} srvhelper.DsRefs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/jenkins/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Delete(input)
}

// @Summary get all jenkins connections
// @Description Get all Jenkins connections
// @Tags plugins/jenkins
// @Success 200  {object} []models.JenkinsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/jenkins/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetAll(input)
}

// @Summary get jenkins connection detail
// @Description Get Jenkins connection detail
// @Tags plugins/jenkins
// @Success 200  {object} models.JenkinsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/jenkins/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetDetail(input)
}

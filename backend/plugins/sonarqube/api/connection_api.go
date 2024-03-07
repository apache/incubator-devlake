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
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
	"github.com/apache/incubator-devlake/server/api/shared"
)

type validation struct {
	Valid bool `json:"valid"`
}
type SonarqubeTestConnResponse struct {
	shared.ApiBody
	Connection *models.SonarqubeConn
}

func testConnection(ctx context.Context, connection models.SonarqubeConn) (*plugin.ApiResourceOutput, errors.Error) {
	// validate
	if vld != nil {
		if err := vld.Struct(connection); err != nil {
			return nil, errors.Default.Wrap(err, "error validating target")
		}
	}
	apiClient, err := api.NewApiClientFromConnection(ctx, basicRes, &connection)
	if err != nil {
		return nil, err
	}

	res, err := apiClient.Get("authentication/validate", nil, nil)
	if err != nil {
		return nil, err
	}
	switch res.StatusCode {
	case 200: // right StatusCode
		valid := &validation{}
		err = api.UnmarshalResponse(res, valid)
		if err != nil {
			return nil, err
		}
		body := SonarqubeTestConnResponse{}
		body.Success = true
		body.Message = "success"
		connection = connection.Sanitize()
		body.Connection = &connection
		if !valid.Valid {
			return nil, errors.Default.New("Authentication failed, please check your access token.")
		}
		return &plugin.ApiResourceOutput{Body: body, Status: 200}, nil
	case 401: // error secretKey or nonceStr
		return &plugin.ApiResourceOutput{Body: false, Status: http.StatusBadRequest}, nil
	default: // unknown what happen , back to user
		return &plugin.ApiResourceOutput{Body: res.Body, Status: res.StatusCode}, nil
	}
}

// TestConnection test sonarqube connection options
// @Summary test sonarqube connection
// @Description Test sonarqube Connection
// @Tags plugins/sonarqube
// @Param body body models.SonarqubeConn true "json body"
// @Success 200  {object} SonarqubeTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/sonarqube/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// decode
	var err errors.Error
	var connection models.SonarqubeConn
	if err = api.Decode(input.Body, &connection, vld); err != nil {
		return nil, err
	}
	testConnectionResult, testConnectionErr := testConnection(context.TODO(), connection)
	if testConnectionErr != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, testConnectionErr)
	}
	if testConnectionResult.Status != http.StatusOK {
		errMsg := fmt.Sprintf("Test connection fail, unexpected status code: %d", testConnectionResult.Status)
		return nil, plugin.WrapTestConnectionErrResp(basicRes, errors.Default.New(errMsg))
	}
	return testConnectionResult, nil
}

// TestExistingConnection test sonarqube connection options
// @Summary test sonarqube connection
// @Description Test sonarqube Connection
// @Tags plugins/sonarqube
// @Success 200  {object} SonarqubeTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/sonarqube/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection, err := dsHelper.ConnApi.GetMergedConnection(input)
	if err != nil {
		return nil, errors.Convert(err)
	}
	// test connection
	testConnectionResult, testConnectionErr := testConnection(context.TODO(), connection.SonarqubeConn)
	if testConnectionErr != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, testConnectionErr)
	}
	if testConnectionResult.Status != http.StatusOK {
		errMsg := fmt.Sprintf("Test connection fail, unexpected status code: %d", testConnectionResult.Status)
		return nil, plugin.WrapTestConnectionErrResp(basicRes, errors.Default.New(errMsg))
	}
	return testConnectionResult, nil
}

// PostConnections create sonarqube connection
// @Summary create sonarqube connection
// @Description Create sonarqube connection
// @Tags plugins/sonarqube
// @Param body body models.SonarqubeConnection true "json body"
// @Success 200  {object} models.SonarqubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/sonarqube/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Post(input)
}

// PatchConnection patch sonarqube connection
// @Summary patch sonarqube connection
// @Description Patch sonarqube connection
// @Tags plugins/sonarqube
// @Param body body models.SonarqubeConnection true "json body"
// @Param connectionId path int false "connection ID"
// @Success 200  {object} models.SonarqubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Patch(input)
}

// DeleteConnection delete a sonarqube connection
// @Summary delete a sonarqube connection
// @Description Delete a sonarqube connection
// @Tags plugins/sonarqube
// @Param connectionId path int false "connection ID"
// @Success 200  {object} models.SonarqubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} srvhelper.DsRefs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Delete(input)
}

// ListConnections get all sonarqube connections
// @Summary get all sonarqube connections
// @Description Get all sonarqube connections
// @Tags plugins/sonarqube
// @Success 200  {object} []models.SonarqubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/sonarqube/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetAll(input)
}

// GetConnection get sonarqube connection detail
// @Summary get sonarqube connection detail
// @Description Get sonarqube connection detail
// @Tags plugins/sonarqube
// @Param connectionId path int false "connection ID"
// @Success 200  {object} models.SonarqubeConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetDetail(input)
}

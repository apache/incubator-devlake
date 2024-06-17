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
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/api/azuredevops"
	"github.com/apache/incubator-devlake/plugins/azuredevops_go/models"
	"github.com/apache/incubator-devlake/server/api/shared"
	"net/http"
)

type AzuredevopsTestConnResponse struct {
	shared.ApiBody
}

// TestConnection tests Azure DevOps connection
// @Summary test Azure DevOps connection
// @Description Test Azure DevOps Connection
// @Tags plugins/azuredevops
// @Param body body models.AzuredevopsConn true "json body"
// @Success 200  {object} AzuredevopsTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/azuredevops/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var conn models.AzuredevopsConn

	if err := api.Decode(input.Body, &conn, vld); err != nil {
		return nil, err
	}
	connection := models.AzuredevopsConnection{
		BaseConnection: api.BaseConnection{
			Name: "conn-test",
		},
		AzuredevopsConn: conn,
	}
	body, err := testConnection(context.TODO(), connection)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: body, Status: http.StatusOK}, nil
}

// TestExistingConnection test Azure DevOps connection
// @Summary test Azure DevOps connection
// @Description Test Azure DevOps Connection
// @Tags plugins/azuredevops
// @Success 200  {object} AzuredevopsTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/azuredevops/connections/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection, err := dsHelper.ConnApi.GetMergedConnection(input)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "can't read connection from database")
	}

	body, err := testConnection(context.TODO(), *connection)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: body, Status: http.StatusOK}, nil
}

// PostConnections create Azure DevOps connection
// @Summary create Azure DevOps connection
// @Description Create Azure DevOps connection
// @Tags plugins/azuredevops
// @Param body body models.AzuredevopsConnection true "json body"
// @Success 200  {object} models.AzuredevopsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/azuredevops/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Post(input)
}

// PatchConnection patch Azure DevOps connection
// @Summary patch Azure DevOps connection
// @Description Patch Azure DevOps connection
// @Tags plugins/azuredevops
// @Param body body models.AzuredevopsConnection true "json body"
// @Success 200  {object} models.AzuredevopsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/azuredevops/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Patch(input)
}

// DeleteConnection delete an Azure DevOps connection
// @Summary delete an Azure DevOps connection
// @Description Delete an Azure DevOps connection
// @Tags plugins/azuredevops
// @Success 200  {object} models.AzuredevopsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} services.BlueprintProjectPairs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/azuredevops/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Delete(input)
}

// ListConnections get all Azure DevOps connections
// @Summary get all Azure DevOps connections
// @Description Get all Azure DevOps connections
// @Tags plugins/azuredevops
// @Success 200  {object} []models.AzuredevopsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/azuredevops/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetAll(input)
}

// GetConnection get Azure DevOps connection detail
// @Summary get Azure DevOps connection detail
// @Description Get Azure DevOps connection detail
// @Tags plugins/azuredevops
// @Success 200  {object} models.AzuredevopsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/azuredevops/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetDetail(input)
}

func testConnection(ctx context.Context, connection models.AzuredevopsConnection) (*AzuredevopsTestConnResponse, errors.Error) {
	// validate
	if vld != nil {
		if err := vld.Struct(connection); err != nil {
			return nil, errors.Default.Wrap(err, "error validating connection")
		}
	}
	apiClient, err := api.NewApiClientFromConnection(ctx, basicRes, &connection)
	if err != nil {
		return nil, err
	}

	vsc := azuredevops.NewClient(&connection, apiClient, "https://app.vssps.visualstudio.com/")
	org := connection.Organization

	if org == "" {
		_, err = vsc.GetUserProfile()
	} else {
		args := azuredevops.GetProjectsArgs{
			OrgId: org,
		}
		_, err = vsc.GetProjects(args)
	}
	if err != nil {
		return nil, err
	}

	connection = connection.Sanitize()
	body := AzuredevopsTestConnResponse{}
	body.Success = true
	body.Message = "success"

	return &body, nil
}

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
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/argocd/models"
	"github.com/apache/incubator-devlake/server/api/shared"
)

type ArgocdTestConnResponse struct {
	shared.ApiBody
	Connection *models.ArgocdConn
}

func testConnection(ctx context.Context, connection models.ArgocdConn) (*ArgocdTestConnResponse, errors.Error) {
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

	// Test ArgoCD API by listing applications
	query := url.Values{}
	res, err := apiClient.Get("applications", query, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return nil, errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error - check your token")
	}

	if res.StatusCode == http.StatusForbidden {
		return nil, errors.BadInput.New("token lacks required permissions")
	}

	connection = connection.Sanitize()
	body := ArgocdTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = &connection

	return &body, nil
}

// TestConnection test argocd connection
// @Summary test argocd connection
// @Description Test ArgoCD Connection
// @Tags plugins/argocd
// @Param body body models.ArgocdConn true "json body"
// @Success 200  {object} ArgocdTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/argocd/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var err errors.Error
	var connection models.ArgocdConn
	if err = api.Decode(input.Body, &connection, vld); err != nil {
		return nil, err
	}
	result, err := testConnection(context.TODO(), connection)
	if err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	}
	return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
}

// TestExistingConnection test argocd connection
// @Summary test argocd connection
// @Description Test ArgoCD Connection
// @Tags plugins/argocd
// @Param connectionId path int true "connection ID"
// @Success 200  {object} ArgocdTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/argocd/connections/{connectionId}/test [POST]
func TestExistingConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection, err := dsHelper.ConnApi.GetMergedConnection(input)
	if err != nil {
		return nil, errors.Convert(err)
	}
	if result, err := testConnection(context.TODO(), connection.ArgocdConn); err != nil {
		return nil, plugin.WrapTestConnectionErrResp(basicRes, err)
	} else {
		return &plugin.ApiResourceOutput{Body: result, Status: http.StatusOK}, nil
	}
}

// @Summary create argocd connection
// @Description Create ArgoCD connection
// @Tags plugins/argocd
// @Param body body models.ArgocdConnection true "json body"
// @Success 200  {object} models.ArgocdConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/argocd/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	out, err := dsHelper.ConnApi.Post(input)
	if err != nil {
		return nil, err
	}
	connId := out.Body.(*models.ArgocdConnection).ID
	_, _ = CreateDefaultScopeConfig(connId)
	return out, nil
}

// @Summary patch argocd connection
// @Description Patch ArgoCD connection
// @Tags plugins/argocd
// @Param body body models.ArgocdConnection true "json body"
// @Success 200  {object} models.ArgocdConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/argocd/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Patch(input)
}

// @Summary delete an argocd connection
// @Description Delete an ArgoCD connection
// @Tags plugins/argocd
// @Success 200  {object} models.ArgocdConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} srvhelper.DsRefs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/argocd/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.Delete(input)
}

// @Summary get all argocd connections
// @Description Get all ArgoCD connections
// @Tags plugins/argocd
// @Success 200  {object} []models.ArgocdConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/argocd/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetAll(input)
}

// @Summary get argocd connection detail
// @Description Get ArgoCD connection detail
// @Tags plugins/argocd

// CreateDefaultScopeConfig ensures a default scope config exists
func CreateDefaultScopeConfig(connectionId uint64) (*models.ArgocdScopeConfig, errors.Error) {
	if dsHelper.ScopeConfigSrv == nil {
		return nil, nil
	}
	existing, _ := dsHelper.ScopeConfigSrv.GetAllByConnectionId(connectionId)
	if len(existing) > 0 {
		return existing[0], nil
	}
	cfg := &models.ArgocdScopeConfig{
		ScopeConfig:       models.ArgocdScopeConfig{}.ScopeConfig, // zero
		DeploymentPattern: ".*",
		ProductionPattern: "(?i)(prod|production)",
		EnvNamePattern:    "(?i)prod(.*)",
	}
	cfg.ConnectionId = connectionId
	cfg.Name = "default"
	cfg.Entities = []string{"CICD"}
	err := dsHelper.ScopeConfigSrv.Create(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// @Success 200  {object} models.ArgocdConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"

// @Router /plugins/argocd/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return dsHelper.ConnApi.GetDetail(input)
}

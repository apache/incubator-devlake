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
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	coreModels "github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/webhook/models"
)

// PostConnections
// @Summary create webhook connection
// @Description Create webhook connection, example: {"name":"Webhook data connection name"}
// @Tags plugins/webhook
// @Param body body WebhookConnectionResponse true "json body"
// @Success 200  {object} WebhookConnectionResponse
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// update from request and save to database
	connection := &models.WebhookConnection{}
	tx := basicRes.GetDal().Begin()
	err := connectionHelper.CreateWithTx(tx, connection, input)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logger.Error(err, "transaction Rollback")
		}
		if strings.Contains(err.Error(), "the connection name already exists (400)") {
			return nil, errors.BadInput.New(fmt.Sprintf("A webhook with name %s already exists.", connection.Name))
		}
		return nil, err
	}
	logger.Info("connection: %+v", connection)
	name := apiKeyHelper.GenApiKeyNameForPlugin(pluginName, connection.ID)
	allowedPath := fmt.Sprintf("/plugins/%s/connections/%d/.*", pluginName, connection.ID)
	extra := fmt.Sprintf("connectionId:%d", connection.ID)
	apiKeyRecord, err := apiKeyHelper.CreateForPlugin(tx, input.User, name, pluginName, allowedPath, extra)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logger.Error(err, "transaction Rollback")
		}
		logger.Error(err, "CreateForPlugin")
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		logger.Info("transaction commit: %s", err)
	}

	webhookConnectionResponse, err := formatConnection(connection, false)
	if err != nil {
		return nil, err
	}
	webhookConnectionResponse.ApiKey = apiKeyRecord
	logger.Info("api output connection: %+v", webhookConnectionResponse)

	return &plugin.ApiResourceOutput{Body: webhookConnectionResponse, Status: http.StatusOK}, nil
}

// PatchConnection
// @Summary patch webhook connection
// @Description Patch webhook connection
// @Tags plugins/webhook
// @Param body body models.WebhookConnection true "json body"
// @Success 200  {object} models.WebhookConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection}, nil
}

// PatchConnectionByName
// @Summary patch webhook connection by name
// @Description Patch webhook connection
// @Tags plugins/webhook
// @Param body body models.WebhookConnection true "json body"
// @Success 200  {object} models.WebhookConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/connections/by-name/{connectionName} [PATCH]
func PatchConnectionByName(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.PatchByName(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection}, nil
}

// DeleteConnection
// @Summary delete a webhook connection
// @Description Delete a webhook connection
// @Tags plugins/webhook
// @Success 200  {object} models.WebhookConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} services.BlueprintProjectPairs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, e := strconv.ParseUint(input.Params["connectionId"], 10, 64)
	return deleteConnection(e, connectionId)
}

// DeleteConnectionByName
// @Summary delete a webhook connection by name
// @Description Delete a webhook connection
// @Tags plugins/webhook
// @Success 200  {object} models.WebhookConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 409  {object} services.BlueprintProjectPairs "References exist to this connection"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/connections/by-name/{connectionName} [DELETE]
func DeleteConnectionByName(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.FirstByName(connection, input.Params)

	if err != nil {
		logger.Error(err, "query connection")
		return nil, err
	}

	return deleteConnection(nil, connection.ConnectionId())
}

func deleteConnection(e error, connectionId uint64) (*plugin.ApiResourceOutput, errors.Error) {
	if e != nil {
		return nil, errors.BadInput.WrapRaw(e)
	}
	var connection models.WebhookConnection
	tx := basicRes.GetDal().Begin()
	err := tx.Delete(&connection, dal.Where("id = ?", connectionId))
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logger.Error(err, "transaction Rollback")
		}
		logger.Error(err, "delete connection: %d", connectionId)
		return nil, err
	}
	extra := fmt.Sprintf("connectionId:%d", connectionId)
	err = apiKeyHelper.DeleteForPlugin(tx, pluginName, extra)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			logger.Error(err, "transaction Rollback")
		}
		logger.Error(err, "delete connection extra: %d, name: %s", extra, pluginName)
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		logger.Info("transaction commit: %s", err)
	}

	return &plugin.ApiResourceOutput{Status: http.StatusOK}, nil
}

type WebhookConnectionResponse struct {
	models.WebhookConnection
	PostIssuesEndpoint             string             `json:"postIssuesEndpoint"`
	CloseIssuesEndpoint            string             `json:"closeIssuesEndpoint"`
	PostPullRequestsEndpoint       string             `json:"postPullRequestsEndpoint"`
	PostPipelineTaskEndpoint       string             `json:"postPipelineTaskEndpoint"`
	PostPipelineDeployTaskEndpoint string             `json:"postPipelineDeployTaskEndpoint"`
	ClosePipelineEndpoint          string             `json:"closePipelineEndpoint"`
	ApiKey                         *coreModels.ApiKey `json:"apiKey,omitempty"`
}

// ListConnections
// @Summary get all webhook connections
// @Description Get all webhook connections
// @Tags plugins/webhook
// @Success 200  {object} []WebhookConnectionResponse
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.WebhookConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	responseList := []*WebhookConnectionResponse{}
	for _, connection := range connections {
		webhookConnectionResponse, err := formatConnection(&connection, true)
		if err != nil {
			return nil, err
		}
		responseList = append(responseList, webhookConnectionResponse)
	}
	return &plugin.ApiResourceOutput{Body: responseList, Status: http.StatusOK}, nil
}

// GetConnection
// @Summary get webhook connection detail
// @Description Get webhook connection detail
// @Tags plugins/webhook
// @Success 200  {object} WebhookConnectionResponse
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.First(connection, input.Params)
	return getConnection(err, connection)
}

// GetConnectionByName
// @Summary get webhook connection detail by name
// @Description Get webhook connection detail
// @Tags plugins/webhook
// @Success 200  {object} WebhookConnectionResponse
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/connections/by-name/{connectionName} [GET]
func GetConnectionByName(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.FirstByName(connection, input.Params)
	return getConnection(err, connection)
}

func getConnection(err errors.Error, connection *models.WebhookConnection) (*plugin.ApiResourceOutput, errors.Error) {
	if err != nil {
		logger.Error(err, "query connection")
		return nil, err
	}
	response, err := formatConnection(connection, true)
	return &plugin.ApiResourceOutput{Body: response}, err
}

func formatConnection(connection *models.WebhookConnection, withApiKeyInfo bool) (*WebhookConnectionResponse, errors.Error) {
	response := &WebhookConnectionResponse{WebhookConnection: *connection}
	response.PostIssuesEndpoint = fmt.Sprintf(`/rest/plugins/webhook/connections/%d/issues`, connection.ID)
	response.CloseIssuesEndpoint = fmt.Sprintf(`/rest/plugins/webhook/connections/%d/issue/:issueKey/close`, connection.ID)
	response.PostPullRequestsEndpoint = fmt.Sprintf(`/rest/plugins/webhook/connections/%d/pull_requests`, connection.ID)
	response.PostPipelineTaskEndpoint = fmt.Sprintf(`/rest/plugins/webhook/connections/%d/cicd_tasks`, connection.ID)
	response.PostPipelineDeployTaskEndpoint = fmt.Sprintf(`/rest/plugins/webhook/connections/%d/deployments`, connection.ID)
	response.ClosePipelineEndpoint = fmt.Sprintf(`/rest/plugins/webhook/connections/%d/cicd_pipeline/:pipelineName/finish`, connection.ID)
	if withApiKeyInfo {
		db := basicRes.GetDal()
		apiKeyName := apiKeyHelper.GenApiKeyNameForPlugin(pluginName, connection.ID)
		apiKey, err := apiKeyHelper.GetApiKey(db, dal.Where("name = ?", apiKeyName))
		if err != nil {
			if db.IsErrorNotFound(err) {
				logger.Info("api key with name: %s not found in db", apiKeyName)
			} else {
				logger.Error(err, "query api key from db, name: %s", apiKeyName)
				return nil, err
			}
		} else {
			response.ApiKey = apiKey
			response.ApiKey.RemoveHashedApiKey() // delete the hashed api key to reduce the attack surface.
		}
	}
	return response, nil
}

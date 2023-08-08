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
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/plugins/webhook/models"
	"net/http"
	"strconv"
)

// PostConnections
// @Summary create webhook connection
// @Description Create webhook connection, example: {"name":"Webhook data connection name"}
// @Tags plugins/webhook
// @Param body body models.WebhookConnection true "json body"
// @Success 200  {object} models.WebhookConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// update from request and save to database
	connection := &models.WebhookConnection{}
	tx := basicRes.GetDal().Begin()
	err := connectionHelper.CreateWithTx(tx, connection, input)
	if err != nil {
		return nil, err
	}
	logger.Info("connection: %+v", connection)
	name := fmt.Sprintf("%s-%d", pluginName, connection.ID)
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

	apiOutputConnection := models.ApiOutputWebhookConnection{
		WebhookConnection: *connection,
		ApiKey:            apiKeyRecord,
	}
	logger.Info("api output connection: %+v", apiOutputConnection)

	return &plugin.ApiResourceOutput{Body: apiOutputConnection, Status: http.StatusOK}, nil
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
	connectionId, e := strconv.ParseInt(input.Params["connectionId"], 10, 64)
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
	PostIssuesEndpoint             string `json:"postIssuesEndpoint"`
	CloseIssuesEndpoint            string `json:"closeIssuesEndpoint"`
	PostPipelineTaskEndpoint       string `json:"postPipelineTaskEndpoint"`
	PostPipelineDeployTaskEndpoint string `json:"postPipelineDeployTaskEndpoint"`
	ClosePipelineEndpoint          string `json:"closePipelineEndpoint"`
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
	responseList := []WebhookConnectionResponse{}
	for _, connection := range connections {
		responseList = append(responseList, *formatConnection(&connection))
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
	response := formatConnection(connection)
	return &plugin.ApiResourceOutput{Body: response}, err
}

func formatConnection(connection *models.WebhookConnection) *WebhookConnectionResponse {
	response := &WebhookConnectionResponse{WebhookConnection: *connection}
	response.PostIssuesEndpoint = fmt.Sprintf(`/plugins/webhook/connections/%d/issues`, connection.ID)
	response.CloseIssuesEndpoint = fmt.Sprintf(`/plugins/webhook/connections/%d/issue/:issueKey/close`, connection.ID)
	response.PostPipelineTaskEndpoint = fmt.Sprintf(`/plugins/webhook/connections/%d/cicd_tasks`, connection.ID)
	response.PostPipelineDeployTaskEndpoint = fmt.Sprintf(`/plugins/webhook/connections/%d/deployments`, connection.ID)
	response.ClosePipelineEndpoint = fmt.Sprintf(`/plugins/webhook/connections/%d/cicd_pipeline/:pipelineName/finish`, connection.ID)
	return response
}

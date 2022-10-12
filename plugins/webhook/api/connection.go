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

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/webhook/models"
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
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	// update from request and save to database
	connection := &models.WebhookConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
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
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection}, nil
}

// DeleteConnection
// @Summary delete a webhook connection
// @Description Delete a webhook connection
// @Tags plugins/webhook
// @Success 200  {object} models.WebhookConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/connections/{connectionId} [DELETE]
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &core.ApiResourceOutput{Body: connection}, err
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
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var connections []models.WebhookConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	responseList := []WebhookConnectionResponse{}
	for _, connection := range connections {
		responseList = append(responseList, *formatConnection(&connection))
	}
	return &core.ApiResourceOutput{Body: responseList, Status: http.StatusOK}, nil
}

// GetConnection
// @Summary get webhook connection detail
// @Description Get webhook connection detail
// @Tags plugins/webhook
// @Success 200  {object} WebhookConnectionResponse
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/webhook/connections/{connectionId} [GET]
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.First(connection, input.Params)
	response := formatConnection(connection)
	return &core.ApiResourceOutput{Body: response}, err
}

func formatConnection(connection *models.WebhookConnection) *WebhookConnectionResponse {
	response := &WebhookConnectionResponse{WebhookConnection: *connection}
	response.PostIssuesEndpoint = fmt.Sprintf(`/plugins/webhook/%d/issues`, connection.ID)
	response.CloseIssuesEndpoint = fmt.Sprintf(`/plugins/webhook/%d/issue/:boardKey/:issueKey/close`, connection.ID)
	response.PostPipelineTaskEndpoint = fmt.Sprintf(`/plugins/webhook/%d/cicd_tasks`, connection.ID)
	response.PostPipelineDeployTaskEndpoint = fmt.Sprintf(`/plugins/webhook/%d/deployments`, connection.ID)
	response.ClosePipelineEndpoint = fmt.Sprintf(`/plugins/webhook/%d/cicd_pipeline/:pipelineName/finish`, connection.ID)
	return response
}

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
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/utils"
	"github.com/apache/incubator-devlake/impls/logruslog"
	"net/http"
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
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
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// update from request and save to database
	connection := &models.WebhookConnection{}
	tx := basicRes.GetDal().Begin()
	err := connectionHelper.CreateWithTx(tx, connection, input)
	if err != nil {
		return nil, err
	}

	logruslog.Global.Info("connection: %+v", connection)
	apiKey, hashedApiKey, err := utils.GenerateApiKey(context.Background())
	if err != nil {
		tx.Rollback()
		logruslog.Global.Error(err, "GenerateApiKey")
		return nil, err
	}
	extra, jsonMarshalErr := json.Marshal(map[string]interface{}{
		"connectionId": connection.ID,
	})
	if jsonMarshalErr != nil {
		tx.Rollback()
		return nil, errors.Default.Wrap(err, "marshal webhook api key extra")
	}

	user, email, _ := GetUserInfo(input.Request)
	apiKeyRecord := &models.ApiKey{
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Creator:      user,
		CreatorEmail: email,
		Updater:      user,
		UpdaterEmail: email,
		Name:         fmt.Sprintf("webhook-%d", connection.ID),
		ApiKey:       hashedApiKey,
		ExpiredAt:    nil,
		AllowedPath:  fmt.Sprintf("/plugins/webhook/%d/.*", connection.ID),
		Type:         "plugin:webhook",
		Extra:        extra,
	}
	if err := tx.Create(apiKeyRecord); err != nil {
		tx.Rollback()
		logruslog.Global.Error(err, "Create api key record")
		return nil, err
	}

	apiKeyRecord.ApiKey = apiKey
	apiOutputConnection := models.ApiOutputWebhookConnection{
		WebhookConnection: *connection,
		ApiKey:            apiKeyRecord,
	}
	logruslog.Global.Info("api output connection: %+v", apiOutputConnection)
	tx.Commit()
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
		tx.Rollback()
		logruslog.Global.Error(err, "delete connection: %d", connectionId)
		return nil, err
	}

	// delete api key generated by webhook
	var apiKey models.ApiKey
	clauses := []dal.Clause{
		dal.JSONQueryEqual("extra", "connectionId", connectionId),
	}
	if err := tx.First(&apiKey, clauses...); err != nil {
		// if api key doesn't exist, just commit
		if tx.IsErrorNotFound(err.Unwrap()) {
			tx.Commit()
			return &plugin.ApiResourceOutput{Status: http.StatusOK}, nil
		} else {
			tx.Rollback()
			logruslog.Global.Error(err, "find api key by connection id: %d", connectionId)
			return nil, err
		}
	}
	logruslog.Global.Info("api key: %+v", apiKey)
	if err := tx.Delete(apiKey); err != nil {
		tx.Rollback()
		logruslog.Global.Error(err, "delete api key id: %d", apiKey.ID)
		return nil, err
	}
	tx.Commit()
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
	response.PostIssuesEndpoint = fmt.Sprintf(`/plugins/webhook/%d/issues`, connection.ID)
	response.CloseIssuesEndpoint = fmt.Sprintf(`/plugins/webhook/%d/issue/:issueKey/close`, connection.ID)
	response.PostPipelineTaskEndpoint = fmt.Sprintf(`/plugins/webhook/%d/cicd_tasks`, connection.ID)
	response.PostPipelineDeployTaskEndpoint = fmt.Sprintf(`/plugins/webhook/%d/deployments`, connection.ID)
	response.ClosePipelineEndpoint = fmt.Sprintf(`/plugins/webhook/%d/cicd_pipeline/:pipelineName/finish`, connection.ID)
	return response
}

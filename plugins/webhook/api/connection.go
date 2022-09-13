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
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/webhook/models"
	"net/http"
)

/*
POST /plugins/webhook/connections
{
	"name": "Webhook data connection name"
}
*/
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// update from request and save to database
	connection := &models.WebhookConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
PATCH /plugins/webhook/connections/:connectionId
{
	"name": "Webhook data connection name"
}
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection}, nil
}

/*
DELETE /plugins/webhook/connections/:connectionId
*/
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &core.ApiResourceOutput{Body: connection}, err
}

/*
GET /plugins/webhook/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	var connections []models.WebhookConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

type WebhookConnectionResponse struct {
	models.WebhookConnection
	IssuesEndpoint       string
	CicdPipelineEndpoint string
}

/*
GET /plugins/webhook/connections/:connectionId
{
	"name": "Webhook data connection name",
	"issues_endpoint": "/plugins/webhook/:connectionId/issues",
	"cicd_pipeline_endpoint": "/plugins/webhook/:connectionId/cicd_pipeline"
}
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.WebhookConnection{}
	err := connectionHelper.First(connection, input.Params)
	response := &WebhookConnectionResponse{WebhookConnection: *connection}
	response.IssuesEndpoint = fmt.Sprintf(`/plugins/webhook/%d/issues`, connection.ID)
	response.CicdPipelineEndpoint = fmt.Sprintf(`/plugins/webhook/%d/cicd_pipeline`, connection.ID)
	return &core.ApiResourceOutput{Body: response}, err
}

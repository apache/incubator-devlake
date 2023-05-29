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
	"time"

	"github.com/apache/incubator-devlake/core/errors"

	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/kube_deployment/models"
	"github.com/apache/incubator-devlake/server/api/shared"
)

type KubeDeploymentTestConnResponse struct {
	shared.ApiBody
	Connection *models.KubeDeploymentConn
}

// TODO Please modify the following code to fit your needs
// @Summary test kube_deployment connection
// @Description Test kube_deployment Connection. endpoint: "https://dev.kube_deployment.com/{organization}/
// @Tags plugins/kube_deployment
// @Param body body models.KubeDeploymentConn true "json body"
// @Success 200  {object} KubeDeploymentTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/kube_deployment/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// decode
	var err errors.Error
	var connection models.KubeDeploymentConn
	if err = helper.Decode(input.Body, &connection, vld); err != nil {
		return nil, err
	}
	fmt.Printf("connection endpoint: %v\n", connection.Endpoint)
	// test connection
	apiClient, err := helper.NewApiClient(
		context.TODO(),
		connection.Endpoint,
		map[string]string{
			// "Authorization": fmt.Sprintf("Bearer %v", connection.Token),
		},
		3*time.Second,
		connection.Proxy,
		basicRes,
	)
	if err != nil {
		return nil, err
	}

	res, err := apiClient.Get("", nil, nil)
	if err != nil {
		return nil, err
	}
	// resBody := &models.ApiUserResponse{}
	// err = helper.UnmarshalResponse(res, resBody)
	// if err != nil {
	// 	return nil, err
	// }

	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("unexpected status code: %d", res.StatusCode))
	}
	body := KubeDeploymentTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = &connection
	// output
	return &plugin.ApiResourceOutput{Body: body, Status: 200}, nil
}

// TODO Please modify the folowing code to adapt to your plugin
// @Summary create kube_deployment connection
// @Description Create kube_deployment connection
// @Tags plugins/kube_deployment
// @Param body body models.KubeDeploymentConnection true "json body"
// @Success 200  {object} models.KubeDeploymentConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/kube_deployment/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// update from request and save to database
	connection := &models.KubeDeploymentConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

// TODO Please modify the folowing code to adapt to your plugin
// @Summary patch kube_deployment connection
// @Description Patch kube_deployment connection
// @Tags plugins/kube_deployment
// @Param body body models.KubeDeploymentConnection true "json body"
// @Success 200  {object} models.KubeDeploymentConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/kube_deployment/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.KubeDeploymentConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection}, nil
}

// @Summary delete a kube_deployment connection
// @Description Delete a kube_deployment connection
// @Tags plugins/kube_deployment
// @Success 200  {object} models.KubeDeploymentConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/kube_deployment/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.KubeDeploymentConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &plugin.ApiResourceOutput{Body: connection}, err
}

// @Summary get all kube_deployment connections
// @Description Get all kube_deployment connections
// @Tags plugins/kube_deployment
// @Success 200  {object} []models.KubeDeploymentConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/kube_deployment/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.KubeDeploymentConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

// TODO Please modify the folowing code to adapt to your plugin
// @Summary get kube_deployment connection detail
// @Description Get kube_deployment connection detail
// @Tags plugins/kube_deployment
// @Success 200  {object} models.KubeDeploymentConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/kube_deployment/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.KubeDeploymentConnection{}
	err := connectionHelper.First(connection, input.Params)
	return &plugin.ApiResourceOutput{Body: connection}, err
}

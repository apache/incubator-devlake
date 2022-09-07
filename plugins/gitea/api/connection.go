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
	"net/url"
	"time"

	"github.com/apache/incubator-devlake/plugins/gitea/models"

	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/mitchellh/mapstructure"

	"github.com/apache/incubator-devlake/plugins/core"
)

// @Summary test gitea connection
// @Description Test gitea Connection
// @Tags plugins/gitea
// @Param body body models.TestConnectionRequest true "json body"
// @Success 200  {object} shared.ApiBody "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/gitea/test [POST]
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	var connection models.TestConnectionRequest
	err := mapstructure.Decode(input.Body, &connection)
	if err != nil {
		return nil, err
	}
	err = vld.Struct(connection)
	if err != nil {
		return nil, err
	}
	// test connection
	apiClient, err := helper.NewApiClient(
		context.TODO(),
		connection.Endpoint,
		nil,
		3*time.Second,
		connection.Proxy,
		basicRes,
	)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("token", connection.Token)

	res, err := apiClient.Get("user", query, nil)
	if err != nil {
		return nil, err
	}
	resBody := &models.ApiUserResponse{}
	err = helper.UnmarshalResponse(res, resBody)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil, nil
}

// @Summary create gitea connection
// @Description Create gitea connection
// @Tags plugins/gitea
// @Param body body models.GithubConnection true "json body"
// @Success 200  {object} models.GiteaConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/gitea/connections [POST]
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GiteaConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

// @Summary patch gitea connection
// @Description Patch gitea connection
// @Tags plugins/gitea
// @Param body body models.GithubConnection true "json body"
// @Success 200  {object} models.GiteaConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/gitea/connections/{connectionId} [PATCH]
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GiteaConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

// @Summary delete a gitea connection
// @Description Delete a gitea connection
// @Tags plugins/gitea
// @Success 200  {object} models.GiteaConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/gitea/connections/{connectionId} [DELETE]
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GiteaConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &core.ApiResourceOutput{Body: connection}, err
}

// @Summary get all gitea connections
// @Description Get all gitea connections
// @Tags plugins/gitea
// @Success 200  {object} models.GiteaConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/gitea/connections [GET]
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	var connections []models.GiteaConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: connections}, nil
}

// @Summary get gitea connection detail
// @Description Get gitea connection detail
// @Tags plugins/gitea
// @Success 200  {object} models.GiteaConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/gitea/connections/{connectionId} [GET]
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.GiteaConnection{}
	err := connectionHelper.First(connection, input.Params)
	return &core.ApiResourceOutput{Body: connection}, err
}

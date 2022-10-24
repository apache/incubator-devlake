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
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/pagerduty/models"
	"net/http"
	"time"
)

// @Summary test pagerduty connection
// @Description Test Pagerduty Connection
// @Tags plugins/pagerduty
// @Param body body models.TestConnectionRequest true "json body"
// @Success 200  {object} shared.ApiBody "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/pagerduty/test [POST]
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var params models.TestConnectionRequest
	err := helper.Decode(input.Body, &params, vld)
	if err != nil {
		return nil, err
	}
	apiClient, err := helper.NewApiClient(
		context.TODO(),
		params.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Token token=%s", params.Token),
		},
		3*time.Second,
		params.Proxy,
		basicRes,
	)
	if err != nil {
		return nil, err
	}
	response, err := apiClient.Get("users/me", nil, nil)
	if err != nil {
		return nil, err
	}
	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusBadRequest {
		// error 400 can happen with a valid but non-user token, i.e. it's OK
		return &core.ApiResourceOutput{Body: nil, Status: http.StatusOK}, nil
	}
	return &core.ApiResourceOutput{Body: nil, Status: response.StatusCode}, errors.HttpStatus(response.StatusCode).Wrap(err, "could not validate connection")
}

// @Summary create pagerduty connection
// @Description Create Pagerduty connection
// @Tags plugins/pagerduty
// @Param body body models.PagerDutyConnection true "json body"
// @Success 200  {object} models.PagerDutyConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/pagerduty/connections [POST]
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connection := &models.PagerDutyConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

// @Summary patch pagerduty connection
// @Description Patch Pagerduty connection
// @Tags plugins/pagerduty
// @Param body body models.PagerDutyConnection true "json body"
// @Success 200  {object} models.PagerDutyConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/pagerduty/connections/{connectionId} [PATCH]
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connection := &models.PagerDutyConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

// @Summary delete pagerduty connection
// @Description Delete Pagerduty connection
// @Tags plugins/pagerduty
// @Success 200  {object} models.PagerDutyConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/pagerduty/connections/{connectionId} [DELETE]
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connection := &models.PagerDutyConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &core.ApiResourceOutput{Body: connection}, err
}

// @Summary list pagerduty connections
// @Description List Pagerduty connections
// @Tags plugins/pagerduty
// @Success 200  {object} models.PagerDutyConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/pagerduty/connections [GET]
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var connections []models.PagerDutyConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: connections}, nil
}

// @Summary get pagerduty connection
// @Description Get Pagerduty connection
// @Tags plugins/pagerduty
// @Success 200  {object} models.PagerDutyConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/pagerduty/connections/{connectionId} [GET]
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connection := &models.PagerDutyConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection}, nil
}

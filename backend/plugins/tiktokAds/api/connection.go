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
	"github.com/apache/incubator-devlake/server/api/shared"
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/tiktokAds/models"
)

type TiktokAdsTestConnResponse struct {
	shared.ApiBody
	Connection *models.TiktokAdsConn
}

// @Summary test tiktokAds connection
// @Description Test tiktokAds Connection. endpoint: https://open.tiktokAds.cn/open-apis/
// @Tags plugins/tiktokAds
// @Param body body models.TiktokAdsConn true "json body"
// @Success 200  {object} TiktokAdsTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/tiktokAds/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// process input
	var connection models.TiktokAdsConn
	if err := api.Decode(input.Body, &connection, vld); err != nil {
		return nil, errors.BadInput.Wrap(err, "could not decode request parameters")
	}

	// test connection
	_, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)

	body := TiktokAdsTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = &connection
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: body, Status: 200}, nil
}

// @Summary create tiktokAds connection
// @Description Create tiktokAds connection
// @Tags plugins/tiktokAds
// @Param body body models.TiktokAdsConnection true "json body"
// @Success 200  {object} models.TiktokAdsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/tiktokAds/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.TiktokAdsConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

// @Summary patch tiktokAds connection
// @Description Patch tiktokAds connection
// @Tags plugins/tiktokAds
// @Param body body models.TiktokAdsConnection true "json body"
// @Success 200  {object} models.TiktokAdsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/tiktokAds/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.TiktokAdsConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

// @Summary delete a tiktokAds connection
// @Description Delete a tiktokAds connection
// @Tags plugins/tiktokAds
// @Success 200  {object} models.TiktokAdsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/tiktokAds/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.TiktokAdsConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &plugin.ApiResourceOutput{Body: connection}, err
}

// @Summary get all tiktokAds connections
// @Description Get all tiktokAds connections
// @Tags plugins/tiktokAds
// @Success 200  {object} models.TiktokAdsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/tiktokAds/connections [GET]
func ListConnections(_ *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.TiktokAdsConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}

	return &plugin.ApiResourceOutput{Body: connections}, nil
}

// @Summary get tiktokAds connection detail
// @Description Get tiktokAds connection detail
// @Tags plugins/tiktokAds
// @Success 200  {object} models.TiktokAdsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/tiktokAds/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.TiktokAdsConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection}, err
}

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
	"time"

	"github.com/apache/incubator-devlake/plugins/ae/models"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ApiMeResponse struct {
	Name string `json:"name"`
}

var vld *validator.Validate
var connectionHelper *helper.ConnectionApiHelper

func Init(config *viper.Viper, logger core.Logger, database *gorm.DB) {
	basicRes := helper.NewDefaultBasicRes(config, logger, database)
	vld = validator.New()
	connectionHelper = helper.NewConnectionHelper(
		basicRes,
		vld,
	)
}

/*
GET /plugins/ae/test/
*/
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// decode
	var err error
	var connection models.TestConnectionRequest
	err = mapstructure.Decode(input.Body, &connection)
	if err != nil {
		return nil, err
	}
	// validate
	err = vld.Struct(connection)
	if err != nil {
		return nil, err
	}

	// load and process cconfiguration
	endpoint := connection.Endpoint
	appId := connection.AppId
	secretKey := connection.SecretKey
	proxy := connection.Proxy

	apiClient, err := helper.NewApiClient(endpoint, nil, 3*time.Second, proxy, nil)
	if err != nil {
		return nil, err
	}
	apiClient.SetBeforeFunction(func(req *http.Request) error {
		nonceStr := core.RandLetterBytes(8)
		timestamp := fmt.Sprintf("%v", time.Now().Unix())
		sign := models.GetSign(req.URL.Query(), appId, secretKey, nonceStr, timestamp)
		req.Header.Set("x-ae-app-id", appId)
		req.Header.Set("x-ae-timestamp", timestamp)
		req.Header.Set("x-ae-nonce-str", nonceStr)
		req.Header.Set("x-ae-sign", sign)
		return nil
	})
	res, err := apiClient.Get("projects", nil, nil)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case 200: // right StatusCode
		return &core.ApiResourceOutput{Body: true, Status: 200}, nil
	case 401: // error secretKey or nonceStr
		return &core.ApiResourceOutput{Body: false, Status: res.StatusCode}, nil
	default: // unknow what happen , back to user
		return &core.ApiResourceOutput{Body: res.Body, Status: res.StatusCode}, nil
	}
}

/*
POST /plugins/ae/connections
*/
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.AeConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
GET /plugins/ae/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	var connections []models.AeConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

/*
GET /plugins/ae/connections/:connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.AeConnection{}
	err := connectionHelper.First(connection, input.Params)
	return &core.ApiResourceOutput{Body: connection}, err
}

/*
PATCH /plugins/ae/connections/:connectionId
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.AeConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
DELETE /plugins/ae/connections/:connectionId
*/
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.AeConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &core.ApiResourceOutput{Body: connection}, err
}

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
	"github.com/apache/incubator-devlake/plugins/feishu/apimodels"
	"github.com/apache/incubator-devlake/plugins/feishu/models"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/mitchellh/mapstructure"
	"net/http"

	"github.com/apache/incubator-devlake/plugins/core"
)

/*
POST /plugins/feishu/test
*/
func TestConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// process input
	var params models.TestConnectionRequest
	err := mapstructure.Decode(input.Body, &params)
	if err != nil {
		return nil, err
	}
	err = vld.Struct(params)
	if err != nil {
		return nil, err
	}

	authApiClient, err := helper.NewApiClient(context.TODO(), params.Endpoint, nil, 0, params.Proxy)
	if err != nil {
		return nil, err
	}

	// request for access token
	tokenReqBody := &apimodels.ApiAccessTokenRequest{
		AppId:     params.AppId,
		AppSecret: params.SecretKey,
	}
	tokenRes, err := authApiClient.Post("open-apis/auth/v3/tenant_access_token/internal", nil, tokenReqBody, nil)
	if err != nil {
		return nil, err
	}
	tokenResBody := &apimodels.ApiAccessTokenResponse{}
	err = helper.UnmarshalResponse(tokenRes, tokenResBody)
	if err != nil {
		return nil, err
	}
	if tokenResBody.AppAccessToken == "" && tokenResBody.TenantAccessToken == "" {
		return nil, fmt.Errorf("failed to request access token")
	}

	// output
	return nil, nil
}

/*
POST /plugins/feishu/connections
*/
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.FeishuConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
PATCH /plugins/feishu/connections/:connectionId
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.FeishuConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

/*
DELETE /plugins/feishu/connections/:connectionId
*/
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.FeishuConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &core.ApiResourceOutput{Body: connection}, err
}

/*
GET /plugins/feishu/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	var connections []models.FeishuConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: connections}, nil
}

/*
GET /plugins/feishu/connections/:connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.FeishuConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection}, err
}

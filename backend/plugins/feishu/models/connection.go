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

package models

import (
	"fmt"
	"github.com/apache/incubator-devlake/core/utils"
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/feishu/apimodels"
)

// FeishuConn holds the essential information to connect to the Feishu API
type FeishuConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.AppKey         `mapstructure:",squash"`
}

func (conn *FeishuConn) Sanitize() FeishuConn {
	conn.SecretKey = utils.SanitizeString(conn.SecretKey)
	return *conn
}

func (conn *FeishuConn) PrepareApiClient(apiClient plugin.ApiClient) errors.Error {
	// request for access token
	tokenReqBody := &apimodels.ApiAccessTokenRequest{
		AppId:     conn.AppId,
		AppSecret: conn.SecretKey,
	}
	tokenRes, err := apiClient.Post("auth/v3/tenant_access_token/internal", nil, tokenReqBody, nil)
	if err != nil {
		return err
	}

	if tokenRes.StatusCode == http.StatusUnauthorized {
		return errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error when get tenant_access_token")
	}

	tokenResBody := &apimodels.ApiAccessTokenResponse{}
	err = helper.UnmarshalResponse(tokenRes, tokenResBody)
	if err != nil {
		return err
	}
	if tokenResBody.AppAccessToken == "" && tokenResBody.TenantAccessToken == "" {
		return errors.Default.New("failed to request access token")
	}
	apiClient.SetHeaders(map[string]string{
		"Authorization": fmt.Sprintf("Bearer %v", tokenResBody.TenantAccessToken),
	})
	return nil
}

// FeishuConnection holds FeishuConn plus ID/Name for database storage
type FeishuConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	FeishuConn            `mapstructure:",squash"`
}

func (FeishuConnection) TableName() string {
	return "_tool_feishu_connections"
}

func (connection FeishuConnection) Sanitize() FeishuConnection {
	connection.FeishuConn = connection.FeishuConn.Sanitize()
	return connection
}

func (connection *FeishuConnection) MergeFromRequest(target *FeishuConnection, body map[string]interface{}) error {
	secretKey := target.SecretKey
	if err := helper.DecodeMapStruct(body, target, true); err != nil {
		return err
	}
	modifiedSecretKey := target.SecretKey
	if modifiedSecretKey == "" || modifiedSecretKey == utils.SanitizeString(secretKey) {
		target.SecretKey = secretKey
	}
	return nil
}

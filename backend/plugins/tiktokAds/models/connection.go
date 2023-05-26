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
	"net/http"

	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api/apihelperabstract"
)

// TiktokAdsConn holds the essential information to connect to the TiktokAds API
type TiktokAdsConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.AppKey         `mapstructure:",squash"`
}

func (conn *TiktokAdsConn) PrepareApiClient(apiClient apihelperabstract.ApiClientAbstract) errors.Error {
	// request for access token
	tokenReqBody := &ApiAccessTokenRequest{
		AppId:    conn.AppId,
		Secret:   conn.SecretKey,
		AuthCode: conn.AuthCode,
	}
	tokenRes, err := apiClient.Post("v1.3/oauth2/access_token/", nil, tokenReqBody, nil)
	if err != nil {
		return err
	}

	if tokenRes.StatusCode == http.StatusUnauthorized {
		return errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error when get tenant_access_token")
	}

	tokenResBody := &ApiAccessTokenResponse{}
	err = helper.UnmarshalResponse(tokenRes, tokenResBody)
	if err != nil {
		return err
	}
	if tokenResBody.AccessToken == "" {
		return errors.Default.New("failed to request access token")
	}
	apiClient.SetHeaders(map[string]string{
		"Authorization": fmt.Sprintf("Access-Token %v", tokenResBody.Data.AccessToken),
	})
	return nil
}

// TiktokAdsConnection holds TiktokAdsConn plus ID/Name for database storage
type TiktokAdsConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	TiktokAdsConn         `mapstructure:",squash"`
}

func (TiktokAdsConnection) TableName() string {
	return "_tool_tiktokAds_connections"
}

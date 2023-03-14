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
	"github.com/apache/incubator-devlake/core/errors"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// TeambitionConn holds the essential information to connect to the TeambitionConn API
type TeambitionConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.AppKey         `mapstructure:",squash"`
	TenantId              string `mapstructure:"tenantId" validate:"required" json:"tenantId"`
	TenantType            string `mapstructure:"tenantType" validate:"required" json:"tenantType"`
}

// TeambitionConnection holds TeambitionConn plus ID/Name for database storage
type TeambitionConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	TeambitionConn        `mapstructure:",squash"`
}

func (tc *TeambitionConn) SetupAuthentication(req *http.Request) errors.Error {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["_appId"] = tc.AppId
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(tc.SecretKey))
	if err != err {
		return errors.Convert(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", tokenString))
	req.Header.Set("X-Tenant-Id", tc.TenantId)
	req.Header.Set("X-Tenant-Type", tc.TenantType)
	return nil
}

func (TeambitionConnection) TableName() string {
	return "_tool_teambition_connections"
}

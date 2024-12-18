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
	"net/url"
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

// PrepareApiClient fetches token from Zentao API for future requests
func (connection ZentaoConn) PrepareApiClient(apiClient plugin.ApiClient) errors.Error {
	// request for access token
	tokenReqBody := &ApiAccessTokenRequest{
		Account:  connection.Username,
		Password: connection.Password,
	}
	tokenRes, err := apiClient.Post("/tokens", nil, tokenReqBody, nil)
	if err != nil {
		return err
	}

	if tokenRes.StatusCode == http.StatusUnauthorized {
		return errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error while to request access token")
	}

	tokenResBody := &ApiAccessTokenResponse{}
	err = helper.UnmarshalResponse(tokenRes, tokenResBody)
	if err != nil {
		return errors.HttpStatus(http.StatusBadRequest).Wrap(err, "failed UnmarshalResponse for tokenResBody")
	}
	if tokenResBody.Token == "" {
		msg := "failed to request access token"
		if tokenResBody.Error != "" {
			msg = tokenResBody.Error
		}
		return errors.HttpStatus(http.StatusBadRequest).New(msg)
	}
	apiClient.SetHeaders(map[string]string{
		"Token": fmt.Sprintf("%v", tokenResBody.Token),
	})
	return nil
}

// ZentaoConn holds the essential information to connect to the Gitlab API
type ZentaoConn struct {
	helper.RestConnection `mapstructure:",squash"`
	helper.BasicAuth      `mapstructure:",squash"`

	DbUrl          string `mapstructure:"dbUrl"  json:"dbUrl" gorm:"serializer:encdec"`
	DbIdleConns    int    `json:"dbIdleConns" mapstructure:"dbIdleConns"`
	DbLoggingLevel string `json:"dbLoggingLevel" mapstructure:"dbLoggingLevel"`
	DbMaxConns     int    `json:"dbMaxConns" mapstructure:"dbMaxConns"`
}

func (connection ZentaoConn) GetHash() string {
	// zentao's token will expire after about 24min, so api client cannot be cached.
	return ""
}

func (connection ZentaoConn) Sanitize() ZentaoConn {
	connection.Password = ""
	if connection.DbUrl != "" {
		connection.DbUrl = connection.SanitizeDbUrl()
	}
	connection.Password = ""
	return connection
}

// ZentaoConnection holds ZentaoConn plus ID/Name for database storage
type ZentaoConnection struct {
	helper.BaseConnection `mapstructure:",squash"`
	ZentaoConn            `mapstructure:",squash"`
}

func (connection ZentaoConn) SanitizeDbUrl() string {
	if connection.DbUrl == "" {
		return connection.DbUrl
	}
	dbUrl := connection.DbUrl
	u, _ := url.Parse(dbUrl)
	if u != nil && u.User != nil {
		password, ok := u.User.Password()
		if ok {
			dbUrl = strings.Replace(dbUrl, password, strings.Repeat("*", len(password)), -1)
		}
	}
	if dbUrl == connection.DbUrl {
		dbUrl = ""
	}
	return dbUrl
}

func (connection ZentaoConnection) GetHash() string {
	return connection.ZentaoConn.GetHash()
}

func (connection ZentaoConnection) Sanitize() ZentaoConnection {
	connection.ZentaoConn = connection.ZentaoConn.Sanitize()
	return connection
}

func (connection *ZentaoConnection) MergeFromRequest(target *ZentaoConnection, body map[string]interface{}) error {
	password := target.Password
	existedDBUrl := target.DbUrl
	existedSanitizedConnectionDBUrl := target.Sanitize().DbUrl
	if err := helper.DecodeMapStruct(body, target, true); err != nil {
		return err
	}

	modifiedPassword := target.Password
	if modifiedPassword == "" {
		target.Password = password
	}

	if existedDBUrl != "" && target.DbUrl != "" && existedSanitizedConnectionDBUrl == target.DbUrl {
		target.DbUrl = existedDBUrl
	}
	return nil
}

// Merge works with the new connection helper.
func (connection ZentaoConnection) Merge(existed, modified *ZentaoConnection) error {
	existedDBUrl := existed.DbUrl
	if existedDBUrl != "" && modified.DbUrl != "" {
		existedSanitizedConnection := existed.Sanitize()
		if existedSanitizedConnection.DbUrl != modified.DbUrl {
			// db url is updated
			existed.DbUrl = modified.DbUrl
		} else {
			// there is no change with db url field.
			// existedDBUrl = origin(modified.DbUrl)
			return nil
		}
	} else {
		existed.DbUrl = modified.DbUrl
	}
	return nil
}

// This object conforms to what the frontend currently expects.
type ZentaoResponse struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
	ZentaoConnection
}

// Using User because it requires authentication.
type ApiUserResponse struct {
	Id   int64
	Name string `json:"name"`
}

func (ZentaoConnection) TableName() string {
	return "_tool_zentao_connections"
}

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
	"time"

	"github.com/apache/incubator-devlake/core/plugin"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api/apihelperabstract"
)

var _ plugin.ApiConnection = (*BambooConnection)(nil)

type BambooConnection struct {
	api.BaseConnection `mapstructure:",squash"`
	BambooConn         `mapstructure:",squash"`
}

// TODO Please modify the following code to fit your needs
// This object conforms to what the frontend currently sends.
type BambooConn struct {
	api.RestConnection `mapstructure:",squash"`
	//TODO you may need to use helper.BasicAuth instead of helper.AccessToken
	api.BasicAuth `mapstructure:",squash"`
}

// PrepareApiClient test api and set the IsPrivateToken,version,UserId and so on.
func (conn *BambooConn) PrepareApiClient(apiClient apihelperabstract.ApiClientAbstract) errors.Error {
	header := http.Header{}
	header.Set("Authorization", fmt.Sprintf("Basic %v", conn.GetEncodedToken()))

	res, err := apiClient.Get("info.json", nil, header)
	if err != nil {
		return errors.HttpStatus(400).New(fmt.Sprintf("Get failed %s", err.Error()))
	}
	repo := &ApiBambooServerInfo{}

	if res.StatusCode == http.StatusUnauthorized {
		return errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error")
	}

	if res.StatusCode != http.StatusOK {
		return errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("unexpected status code: %d", res.StatusCode))
	}

	err = api.UnmarshalResponse(res, repo)

	if err != nil {
		return errors.BadInput.New(fmt.Sprintf("UnmarshalResponse repository failed %s", err.Error()))
	}

	return nil
}

// This object conforms to what the frontend currently expects.
type BambooResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
	BambooConnection
}

type ApiBambooServerInfo struct {
	Version     string     `json:"version"`
	Edition     string     `json:"edition"`
	BuildDate   *time.Time `json:"buildDate"`
	BuildNumber string     `json:"buildNumber"`
	State       string     `json:"state"`
}

type ApiRepository struct {
	Size          int         `json:"size"`
	SearchResults interface{} `json:"searchResults"`
	StartIndex    int         `json:"start-index"`
	MaxResult     int         `json:"max-result"`
}

func (BambooConnection) TableName() string {
	return "_tool_bamboo_connections"
}

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
	"github.com/apache/incubator-devlake/config"
	"github.com/apache/incubator-devlake/plugins/jenkins/models"
	"net/http"
	"time"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/utils"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

var vld = validator.New()

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
	// test connection
	encodedToken := utils.GetEncodedToken(connection.Username, connection.Password)
	apiClient, err := helper.NewApiClient(
		connection.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", encodedToken),
		},
		3*time.Second,
		connection.Proxy,
		nil,
	)
	if err != nil {
		return nil, err
	}
	res, err := apiClient.Get("", nil, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil, nil
}

/*
PATCH /plugins/jenkins/connections/:connectionId
{
	"Endpoint": "your-endpoint",
	"Username": "your-username",
	"Password": "your-password",
}
*/
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	v := config.GetConfig()
	connection := &models.JenkinsConnection{}
	err := helper.EncodeStruct(v, connection, "env")
	if err != nil {
		return nil, err
	}
	// update from request and save to .env
	err = helper.DecodeStruct(v, connection, input.Body, "env")
	if err != nil {
		return nil, err
	}
	err = config.WriteConfig(v)
	if err != nil {
		return nil, err
	}
	response := models.JenkinsResponse{
		JenkinsConnection: *connection,
		Name:              "Jenkins",
		ID:                1,
	}
	return &core.ApiResourceOutput{Body: response, Status: http.StatusOK}, nil
}

/*
GET /plugins/jenkins/connections
*/
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// RETURN ONLY 1 SOURCE (FROM ENV) until multi-connection is developed.
	v := config.GetConfig()
	connection := &models.JenkinsConnection{}

	err := helper.EncodeStruct(v, connection, "env")
	if err != nil {
		return nil, err
	}
	response := models.JenkinsResponse{
		JenkinsConnection: *connection,
		Name:              "Jenkins",
		ID:                1,
	}

	return &core.ApiResourceOutput{Body: []models.JenkinsResponse{response}}, nil
}

/*
GET /plugins/jenkins/connections/:connectionId
*/
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	//  RETURN ONLY 1 SOURCE FROM ENV (Ignore ID until multi-connection is developed.)
	v := config.GetConfig()
	connection := &models.JenkinsConnection{}
	err := helper.EncodeStruct(v, connection, "env")
	if err != nil {
		return nil, err
	}
	response := &models.JenkinsResponse{
		JenkinsConnection: *connection,
		Name:              "Jenkins",
		ID:                1,
	}
	return &core.ApiResourceOutput{Body: response}, nil
}

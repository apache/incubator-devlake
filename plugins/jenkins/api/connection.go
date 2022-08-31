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
	"net/http"
	"time"

	"github.com/apache/incubator-devlake/plugins/jenkins/models"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/utils"
	"github.com/mitchellh/mapstructure"
)

// @Summary test jenkins connection
// @Description Test Jenkins Connection
// @Tags plugins/jenkins
// @Param body body models.TestConnectionRequest true "json body"
// @Success 200  {object} shared.ApiBody "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/jenkins/test [POST]
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
		context.TODO(),
		connection.Endpoint,
		map[string]string{
			"Authorization": fmt.Sprintf("Basic %v", encodedToken),
		},
		3*time.Second,
		connection.Proxy,
		basicRes,
	)
	if err != nil {
		return nil, err
	}
	res, err := apiClient.Get("", nil, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New("unexpected status code when testing connection")
	}
	return nil, nil
}

// @Summary create jenkins connection
// @Description Create Jenkins connection
// @Tags plugins/jenkins
// @Param body body models.JenkinsConnection true "json body"
// @Success 200  {object} models.JenkinsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/jenkins/connections [POST]
func PostConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	// create a new connection
	connection := &models.JenkinsConnection{}

	// update from request and save to database
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

// @Summary patch jenkins connection
// @Description Patch Jenkins connection
// @Tags plugins/jenkins
// @Param body body models.JenkinsConnection true "json body"
// @Success 200  {object} models.JenkinsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/jenkins/connections/{connectionId} [PATCH]
func PatchConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.JenkinsConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: connection}, nil
}

// @Summary delete a jenkins connection
// @Description Delete a Jenkins connection
// @Tags plugins/jenkins
// @Success 200  {object} models.JenkinsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/jenkins/connections/{connectionId} [DELETE]
func DeleteConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.JenkinsConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &core.ApiResourceOutput{Body: connection}, err
}

// @Summary get all jenkins connections
// @Description Get all Jenkins connections
// @Tags plugins/jenkins
// @Success 200  {object} []models.JenkinsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/jenkins/connections [GET]
func ListConnections(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	var connections []models.JenkinsConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}

	return &core.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

// @Summary get jenkins connection detail
// @Description Get Jenkins connection detail
// @Tags plugins/jenkins
// @Success 200  {object} models.JenkinsConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /plugins/jenkins/connections/{connectionId} [GET]
func GetConnection(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	connection := &models.JenkinsConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: connection}, err
}

// @Summary blueprints setting for jenkins
// @Description blueprint setting for jenkins
// @Tags plugins/jenkins
// @Accept application/json
// @Param blueprint body JenkinsBlueprintSetting true "json"
// @Router /blueprints/jenkins/blueprint-setting [post]
func PostJenkinsBluePrint(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	blueprint := &JenkinsBlueprintSetting{}
	return &core.ApiResourceOutput{Body: blueprint, Status: http.StatusOK}, nil
}

type JenkinsBlueprintSetting []struct {
	Version     string `json:"version"`
	Connections []struct {
		Plugin       string `json:"plugin"`
		ConnectionID int    `json:"connectionId"`
		Scope        []struct {
			Options struct {
			} `json:"options"`
			Entities []string `json:"entities"`
		} `json:"scope"`
	} `json:"connections"`
}

// @Summary pipelines plan for jenkins
// @Description pipelines plan for jenkins
// @Tags plugins/jenkins
// @Accept application/json
// @Param blueprint body JenkinsPipelinePlan true "json"
// @Router /pipelines/jenkins/pipeline-plan [post]
func PostJenkinsPipeline(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
	blueprint := &JenkinsPipelinePlan{}
	return &core.ApiResourceOutput{Body: blueprint, Status: http.StatusOK}, nil
}

type JenkinsPipelinePlan [][]struct {
	Plugin   string   `json:"plugin"`
	Subtasks []string `json:"subtasks"`
	Options  struct {
		ConnectionID int `json:"connectionId"`
	} `json:"options"`
}

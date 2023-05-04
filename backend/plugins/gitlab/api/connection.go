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
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/server/api/shared"
)

type GitlabTestConnResponse struct {
	shared.ApiBody
	Connection *models.GitlabConn
}

// @Summary test gitlab connection
// @Description Test gitlab Connection
// @Tags plugins/gitlab
// @Param body body models.GitlabConn true "json body"
// @Success 200  {object} GitlabTestConnResponse "Success"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitlab/test [POST]
func TestConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// decode
	var err errors.Error
	var connection models.GitlabConn
	if err = api.Decode(input.Body, &connection, vld); err != nil {
		return nil, err
	}

	apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, &connection)
	if err != nil {
		return nil, err
	}

	// check API/read_api permissions
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%v", 1))
	query.Set("per_page", fmt.Sprintf("%v", 1))
	res, err := apiClient.Get("projects", query, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return nil, errors.HttpStatus(http.StatusBadRequest).New("StatusUnauthorized error when testing api or read_api permissions")
	}

	if res.StatusCode == http.StatusForbidden {
		return nil, errors.BadInput.New("token need api or read_api permissions scope")
	}

	body := GitlabTestConnResponse{}
	body.Success = true
	body.Message = "success"
	body.Connection = &connection

	return &plugin.ApiResourceOutput{Body: body, Status: http.StatusOK}, nil
}

// @Summary create gitlab connection
// @Description Create gitlab connection
// @Tags plugins/gitlab
// @Param body body models.GitlabConnection true "json body"
// @Success 200  {object} models.GitlabConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitlab/connections [POST]
func PostConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	// update from request and save to database
	connection := &models.GitlabConnection{}
	err := connectionHelper.Create(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection, Status: http.StatusOK}, nil
}

// @Summary patch gitlab connection
// @Description Patch gitlab connection
// @Tags plugins/gitlab
// @Param body body models.GitlabConnection true "json body"
// @Success 200  {object} models.GitlabConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId} [PATCH]
func PatchConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.GitlabConnection{}
	err := connectionHelper.Patch(connection, input)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection}, nil
}

// @Summary delete a gitlab connection
// @Description Delete a gitlab connection
// @Tags plugins/gitlab
// @Success 200  {object} models.GitlabConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId} [DELETE]
func DeleteConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.GitlabConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	err = connectionHelper.Delete(connection)
	return &plugin.ApiResourceOutput{Body: connection}, err
}

// @Summary get all gitlab connections
// @Description Get all gitlab connections
// @Tags plugins/gitlab
// @Success 200  {object} []models.GitlabConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitlab/connections [GET]
func ListConnections(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var connections []models.GitlabConnection
	err := connectionHelper.List(&connections)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connections, Status: http.StatusOK}, nil
}

// @Summary get gitlab connection detail
// @Description Get gitlab connection detail
// @Tags plugins/gitlab
// @Success 200  {object} models.GitlabConnection
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /plugins/gitlab/connections/{connectionId} [GET]
func GetConnection(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connection := &models.GitlabConnection{}
	err := connectionHelper.First(connection, input.Params)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Body: connection}, err
}

// @Summary pipelines plan for gitlab
// @Description pipelines plan for gitlab
// @Tags plugins/gitlab
// @Accept application/json
// @Param blueprint body GitlabPipelinePlan true "json"
// @Router /pipelines/gitlab/pipeline-plan [post]
func PostGitlabPipeline(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	blueprint := &GitlabPipelinePlan{}
	return &plugin.ApiResourceOutput{Body: blueprint, Status: http.StatusOK}, nil
}

type GitlabPipelinePlan [][]struct {
	Plugin   string   `json:"plugin"`
	Subtasks []string `json:"subtasks"`
	Options  struct {
		ConnectionID   int `json:"connectionId"`
		ProjectId      int `json:"projectId"`
		Since          string
		Transformation models.GitlabTransformationRule `json:"transformation"`
	} `json:"options"`
}

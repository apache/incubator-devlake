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
	"encoding/json"
	"fmt"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	aha "github.com/apache/incubator-devlake/helpers/pluginhelper/api/apihelperabstract"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/apache/incubator-devlake/plugins/jira/tasks"
	"github.com/apache/incubator-devlake/plugins/jira/tasks/apiv2models"
	"io"
	"net/http"
)

type ScopeRes struct {
	models.JiraBoard
	TransformationRuleName string `json:"transformationRuleName,omitempty"`
}

type ScopeReq api.ScopeReq[models.JiraBoard]

// PutScope create or update jira board
// @Summary create or update jira board
// @Description Create or update Jira board
// @Tags plugins/jira
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scope body ScopeReq true "json"
// @Success 200  {object} []models.JiraBoard
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/scopes [PUT]
func PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.Put(input)
}

// UpdateScope patch to jira board
// @Summary patch to jira board
// @Description patch to jira board
// @Tags plugins/jira
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scopeId path int false "board ID"
// @Param scope body models.JiraBoard true "json"
// @Success 200  {object} models.JiraBoard
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/scopes/{scopeId} [PATCH]
func UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.Update(input, "board_id")
}

// GetScopeList get Jira boards
// @Summary get Jira boards
// @Description get Jira boards
// @Tags plugins/jira
// @Param connectionId path int false "connection ID"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} []ScopeRes
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/scopes/ [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.GetScopeList(input)
}

// GetScope get one Jira board
// @Summary get one Jira board
// @Description get one Jira board
// @Tags plugins/jira
// @Param connectionId path int false "connection ID"
// @Param scopeId path int false "board ID"
// @Success 200  {object} ScopeRes
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/scopes/{scopeId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	return scopeHelper.GetScope(input, "board_id")
}

func GetApiJira(op *tasks.JiraOptions, apiClient aha.ApiClientAbstract) (*apiv2models.Board, errors.Error) {
	boardRes := &apiv2models.Board{}
	res, err := apiClient.Get(fmt.Sprintf("agile/1.0/board/%d", op.BoardId), nil, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.HttpStatus(res.StatusCode).New(fmt.Sprintf("unexpected status code when requesting repo detail from %s", res.Request.URL.String()))
	}
	body, err := errors.Convert01(io.ReadAll(res.Body))
	if err != nil {
		return nil, err
	}
	err = errors.Convert(json.Unmarshal(body, boardRes))
	if err != nil {
		return nil, err
	}
	return boardRes, nil
}

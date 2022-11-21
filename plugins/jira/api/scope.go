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
	"net/http"
	"strconv"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/mitchellh/mapstructure"
)

type putBoardRequest struct {
	ProjectId uint   `json:"projectId"`
	Name      string `json:"name"`
	Self      string `json:"self"`
	Type      string `json:"type"`
}

// PutScope create or update jira board
// @Summary create or update jira board
// @Description Create or update Jira board
// @Tags plugins/jira
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param boardId path int false "board ID"
// @Param scope body putBoardRequest true "json"
// @Success 200  {object} core.ApiResourceOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/scopes/{boardId} [PUT]
func PutScope(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connectionId, boardId := extractParam(input.Params)
	if connectionId*boardId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	var req putBoardRequest
	err := mapstructure.Decode(input.Body, &req)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error decoding map into putBoardRequest")
	}
	board := &models.JiraBoard{
		ConnectionId: connectionId,
		BoardId:      boardId,
		ProjectId:    req.ProjectId,
		Name:         req.Name,
		Self:         req.Self,
		Type:         req.Type,
	}
	err = basicRes.GetDal().CreateOrUpdate(&board)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving JiraBoard")
	}
	return &core.ApiResourceOutput{Status: http.StatusOK}, nil
}

// GetScopeList get Jira boards
// @Summary get Jira boards
// @Description get Jira boards
// @Tags plugins/jira
// @Param connectionId path int false "connection ID"
// @Success 200  {object} []models.JiraBoard
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/scopes/ [GET]
func GetScopeList(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var boards []models.JiraBoard
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	err := basicRes.GetDal().All(&boards, dal.Where("connection_id = ?", connectionId))
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: boards, Status: http.StatusOK}, nil
}

// GetScope get one Jira board
// @Summary get one Jira board
// @Description get one Jira board
// @Tags plugins/jira
// @Param connectionId path int false "connection ID"
// @Param boardId path int false "board ID"
// @Success 200  {object} models.JiraBoard
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/connections/{connectionId}/scopes/{boardId} [GET]
func GetScope(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var board models.JiraBoard
	connectionId, boardId := extractParam(input.Params)
	if connectionId*boardId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	err := basicRes.GetDal().First(&board, dal.Where("connection_id = ? AND board_id = ?", connectionId, board))
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: board, Status: http.StatusOK}, nil
}

func extractParam(params map[string]string) (uint64, uint64) {
	connectionId, _ := strconv.ParseUint(params["connectionId"], 10, 64)
	boardId, _ := strconv.ParseUint(params["boardId"], 10, 64)
	return connectionId, boardId
}

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
	"github.com/apache/incubator-devlake/plugins/trello/models"
	"net/http"
	"strconv"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/mitchellh/mapstructure"
)

type apiBoard struct {
	models.TrelloBoard
	TransformationRuleName string `json:"transformationRuleName,omitempty"`
}

type req struct {
	Data []*models.TrelloBoard `json:"data"`
}

// PutScope create or update trello board
// @Summary create or update trello board
// @Description Create or update trello board
// @Tags plugins/trello
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scope body req true "json"
// @Success 200  {object} []models.TrelloBoard
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/trello/connections/{connectionId}/scopes [PUT]
func PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, _ := extractParam(input.Params)

	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}
	var boards req
	err := errors.Convert(mapstructure.Decode(input.Body, &boards))
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "decoding Trello board error")
	}
	keeper := make(map[string]struct{})
	for _, board := range boards.Data {
		if _, ok := keeper[board.BoardId]; ok {
			return nil, errors.BadInput.New("duplicated item")
		} else {
			keeper[board.BoardId] = struct{}{}
		}
		board.ConnectionId = connectionId
		err = verifyBoard(board)
		if err != nil {
			return nil, err
		}
	}
	err = basicRes.GetDal().CreateOrUpdate(boards.Data)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving TrelloBoard")
	}
	return &plugin.ApiResourceOutput{Body: boards.Data, Status: http.StatusOK}, nil
}

// UpdateScope patch to trello board
// @Summary patch to trello board
// @Description patch to trello board
// @Tags plugins/trello
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param boardId path string false "board ID"
// @Param scope body models.TrelloBoard true "json"
// @Success 200  {object} models.TrelloBoard
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/trello/connections/{connectionId}/scopes/{boardId} [PATCH]
func UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, boardId := extractParam(input.Params)
	if connectionId == 0 || boardId == "" {
		return nil, errors.BadInput.New("invalid connectionId or boardId")
	}
	var board models.TrelloBoard
	err := basicRes.GetDal().First(&board, dal.Where("connection_id = ? AND board_id = ?", connectionId, boardId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "getting TrelloBoard error")
	}
	err = api.DecodeMapStruct(input.Body, &board, true)
	if err != nil {
		return nil, errors.Default.Wrap(err, "patch trello board error")
	}
	err = verifyBoard(&board)
	if err != nil {
		return nil, err
	}
	err = basicRes.GetDal().Update(board)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving TrelloBoard")
	}
	return &plugin.ApiResourceOutput{Body: board, Status: http.StatusOK}, nil
}

// GetScopeList get Trello boards
// @Summary get Trello boards
// @Description get Trello boards
// @Tags plugins/trello
// @Param connectionId path int false "connection ID"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} []apiBoard
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/trello/connections/{connectionId}/scopes/ [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var boards []models.TrelloBoard
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	limit, offset := api.GetLimitOffset(input.Query, "pageSize", "page")
	err := basicRes.GetDal().All(&boards, dal.Where("connection_id = ?", connectionId), dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, err
	}
	var ruleIds []uint64
	for _, board := range boards {
		if board.TransformationRuleId > 0 {
			ruleIds = append(ruleIds, board.TransformationRuleId)
		}
	}
	var rules []models.TrelloTransformationRule
	if len(ruleIds) > 0 {
		err = basicRes.GetDal().All(&rules, dal.Where("id IN (?)", ruleIds))
		if err != nil {
			return nil, err
		}
	}
	names := make(map[uint64]string)
	for _, rule := range rules {
		names[rule.ID] = rule.Name
	}
	var apiBoards []apiBoard
	for _, board := range boards {
		apiBoards = append(apiBoards, apiBoard{board, names[board.TransformationRuleId]})
	}
	return &plugin.ApiResourceOutput{Body: apiBoards, Status: http.StatusOK}, nil
}

// GetScope get one Trello board
// @Summary get one Trello board
// @Description get one Trello board
// @Tags plugins/trello
// @Param connectionId path int false "connection ID"
// @Param boardId path string false "board ID"
// @Success 200  {object} models.TrelloBoard
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/trello/connections/{connectionId}/scopes/{boardId} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var board models.TrelloBoard
	connectionId, boardId := extractParam(input.Params)
	if connectionId == 0 || boardId == "" {
		return nil, errors.BadInput.New("invalid path params")
	}
	db := basicRes.GetDal()
	err := db.First(&board, dal.Where("connection_id = ? AND board_id = ?", connectionId, boardId))
	if db.IsErrorNotFound(err) {
		return nil, errors.NotFound.New("record not found")
	}
	if err != nil {
		return nil, err
	}
	var rule models.TrelloTransformationRule
	if board.TransformationRuleId > 0 {
		err = basicRes.GetDal().First(&rule, dal.Where("id = ?", board.TransformationRuleId))
		if err != nil {
			return nil, err
		}
	}
	return &plugin.ApiResourceOutput{Body: apiBoard{board, rule.Name}, Status: http.StatusOK}, nil
}

func extractParam(params map[string]string) (uint64, string) {
	connectionId, _ := strconv.ParseUint(params["connectionId"], 10, 64)
	return connectionId, params["boardId"]
}

func verifyBoard(board *models.TrelloBoard) errors.Error {
	if board.ConnectionId == 0 {
		return errors.BadInput.New("invalid connectionId")
	}
	if board.BoardId == "" {
		return errors.BadInput.New("invalid boardId")
	}
	return nil
}

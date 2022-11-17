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
	"net/url"
	"strconv"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/core/dal"
	"github.com/apache/incubator-devlake/plugins/jira/models"
	"github.com/mitchellh/mapstructure"
)

type putBoardRequest struct {
	ConnectionId uint64 `json:"connectionId"`
	BoardId      uint64 `json:"boardId"`
	ProjectId    uint   `json:"projectId"`
	Name         string `json:"name"`
	Self         string `json:"self"`
	Type         string `json:"type"`
}

// PutScope create or update jira board
// @Summary create or update jira board
// @Description Create or update Jira board
// @Tags plugins/jira
// @Accept application/json
// @Param scope body putBoardRequest true "json"
// @Success 200  {object} core.ApiResourceOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/scopes [PUT]
func PutScope(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	// update from request and save to database
	var req putBoardRequest
	err := mapstructure.Decode(input.Body, &req)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error decoding map into putBoardRequest")
	}
	board := &models.JiraBoard{
		ConnectionId: req.ConnectionId,
		BoardId:      req.BoardId,
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

// DeleteScope delete a jira board
// @Summary delete a jira board
// @Description delete a jira board
// @Tags plugins/jira
// @Param connectionId query int false "connection ID"
// @Param boardId query int false "board ID"
// @Success 200  {object} core.ApiResourceOutput
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/scopes [DELETE]
func DeleteScope(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	connectionId, boardId := extractQuery(input.Query)
	if connectionId == 0 {
		return nil, errors.Default.New("invalid connectionId")
	}
	if boardId == 0 {
		return nil, errors.Default.New("invalid boardId")
	}
	err := basicRes.GetDal().Delete(&models.JiraBoard{}, dal.Where("connection_id = ? AND board_id = ?", connectionId, boardId))
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on deleting JiraBoard")
	}
	return &core.ApiResourceOutput{Status: http.StatusOK}, nil
}

// GetScope get Jira board
// @Summary get Jira board
// @Description get Jira board
// @Tags plugins/jira
// @Param connectionId query int false "connection ID"
// @Param boardId query int false "board ID"
// @Success 200  {object} []models.JiraBoard
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/jira/scopes [GET]
func GetScope(input *core.ApiResourceInput) (*core.ApiResourceOutput, errors.Error) {
	var boards []models.JiraBoard
	var clauses []dal.Clause
	connectionId, boardId := extractQuery(input.Query)
	if connectionId > 0 {
		clauses = append(clauses, dal.Where("connection_id = ?", connectionId))
	}
	if boardId > 0 {
		clauses = append(clauses, dal.Where("board_id = ?", boardId))
	}
	err := basicRes.GetDal().All(&boards, clauses...)
	if err != nil {
		return nil, err
	}
	return &core.ApiResourceOutput{Body: boards, Status: http.StatusOK}, nil
}

func extractQuery(query url.Values) (uint64, uint64) {
	var connectionId, boardId uint64
	cid := query["connectionId"]
	if len(cid) > 0 {
		if connectionId, _ = strconv.ParseUint(cid[0], 10, 64); connectionId > 0 {
		}
	}
	bid := query["boardId"]
	if len(bid) > 0 {
		if boardId, _ = strconv.ParseUint(bid[0], 10, 64); boardId > 0 {
		}
	}
	return connectionId, boardId
}

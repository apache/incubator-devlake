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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/google/uuid"
)

type paginatedTeams struct {
	Count int64      `json:"count"`
	Teams []teamTree `json:"teams"`
}

type teamTree struct {
	Id           string     `json:"id"`
	Name         string     `json:"name"`
	Alias        string     `json:"alias"`
	ParentId     string     `json:"parentId"`
	SortingIndex int        `json:"sortingIndex"`
	UserCount    int        `json:"userCount"`
	Children     []teamTree `json:"children,omitempty"`
}

// ListTeams returns teams with pagination support
// @Summary      List teams
// @Description  GET /plugins/org/teams?page=1&pageSize=50
// @Tags 		 plugins/org
// @Produce      json
// @Param        page      query  int     false  "page number (default 1)"
// @Param        pageSize  query  int     false  "page size (default 50)"
// @Param        name      query  string  false  "filter by name (case-insensitive, partial match)"
// @Param        grouped   query  bool    false  "when true, returns parent teams with nested children"
// @Success      200  {object}  paginatedTeams
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/teams [get]
func (h *Handlers) ListTeams(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	page, pageSize := 1, 50
	if p := input.Query.Get("page"); p != "" {
		if v, e := strconv.Atoi(p); e == nil && v > 0 {
			page = v
		}
	}
	if ps := input.Query.Get("pageSize"); ps != "" {
		if v, e := strconv.Atoi(ps); e == nil && v > 0 {
			pageSize = v
		}
	}
	grouped := false
	if groupedParam := input.Query.Get("grouped"); groupedParam != "" {
		groupedValue, parseErr := strconv.ParseBool(groupedParam)
		if parseErr != nil {
			return nil, errors.BadInput.Wrap(parseErr, "grouped must be a boolean value")
		}
		grouped = groupedValue
	}
	nameFilter := input.Query.Get("name")
	teams, count, err := h.store.findTeamsPaginated(page, pageSize, nameFilter, grouped)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{
		Body:   paginatedTeams{Count: count, Teams: teams},
		Status: http.StatusOK,
	}, nil
}

// GetTeamById returns a single team by ID
// @Summary      Get a team
// @Description  get a team by ID
// @Tags 		 plugins/org
// @Produce      json
// @Param        teamId  path  string  true  "team ID"
// @Success      200  {object}  team
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 404  {object} shared.ApiBody "Not Found"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/teams/{teamId} [get]
func (h *Handlers) GetTeamById(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	teamId := input.Params["teamId"]
	if teamId == "" {
		return nil, errors.BadInput.New("teamId is required")
	}
	t, err := h.store.findTeamById(teamId)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{
		Body: team{
			Id:           t.Id,
			Name:         t.Name,
			Alias:        t.Alias,
			ParentId:     t.ParentId,
			SortingIndex: t.SortingIndex,
		},
		Status: http.StatusOK,
	}, nil
}

type createTeamsRequest struct {
	Teams []team `json:"teams"`
}

// CreateTeams creates one or more teams
// @Summary      Create teams
// @Description  create one or more teams
// @Tags 		 plugins/org
// @Accept       json
// @Produce      json
// @Param        body  body  createTeamsRequest  true  "teams to create"
// @Success      201  {array}  team
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/teams [post]
func (h *Handlers) CreateTeams(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var req createTeamsRequest
	err := helper.Decode(input.Body, &req, nil)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid request body")
	}
	if len(req.Teams) == 0 {
		return nil, errors.BadInput.New("at least one team is required")
	}
	var created []team
	for _, t := range req.Teams {
		id := uuid.New().String()
		domainTeam := &crossdomain.Team{
			DomainEntity: domainlayer.DomainEntity{Id: id},
			Name:         t.Name,
			Alias:        t.Alias,
			ParentId:     t.ParentId,
			SortingIndex: t.SortingIndex,
		}
		if err := h.store.createTeam(domainTeam); err != nil {
			return nil, err
		}
		t.Id = id
		created = append(created, t)
	}
	return &plugin.ApiResourceOutput{Body: created, Status: http.StatusCreated}, nil
}

// UpdateTeamById updates a team by ID
// @Summary      Update a team
// @Description  update a team by ID
// @Tags 		 plugins/org
// @Accept       json
// @Produce      json
// @Param        teamId  path  string  true  "team ID"
// @Param        body    body  team    true  "team fields to update"
// @Success      200  {object}  team
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 404  {object} shared.ApiBody "Not Found"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/teams/{teamId} [put]
func (h *Handlers) UpdateTeamById(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	teamId := input.Params["teamId"]
	if teamId == "" {
		return nil, errors.BadInput.New("teamId is required")
	}
	existing, err := h.store.findTeamById(teamId)
	if err != nil {
		return nil, err
	}
	var t team
	if e := helper.Decode(input.Body, &t, nil); e != nil {
		return nil, errors.BadInput.Wrap(e, "invalid request body")
	}
	existing.Name = t.Name
	existing.Alias = t.Alias
	existing.ParentId = t.ParentId
	existing.SortingIndex = t.SortingIndex
	if err := h.store.updateTeam(existing); err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{
		Body: team{
			Id:           existing.Id,
			Name:         existing.Name,
			Alias:        existing.Alias,
			ParentId:     existing.ParentId,
			SortingIndex: existing.SortingIndex,
		},
		Status: http.StatusOK,
	}, nil
}

// DeleteTeamById deletes a team by ID and its associated team_users
// @Summary      Delete a team
// @Description  delete a team by ID (cascades to team_users)
// @Tags 		 plugins/org
// @Param        teamId  path  string  true  "team ID"
// @Success      200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/teams/{teamId} [delete]
func (h *Handlers) DeleteTeamById(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	teamId := input.Params["teamId"]
	if teamId == "" {
		return nil, errors.BadInput.New("teamId is required")
	}
	if err := h.store.deleteTeam(teamId); err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Status: http.StatusOK}, nil
}

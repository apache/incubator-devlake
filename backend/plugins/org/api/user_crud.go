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
	"strings"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models/domainlayer"
	"github.com/apache/incubator-devlake/core/models/domainlayer/crossdomain"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/google/uuid"
)

type paginatedUsers struct {
	Count int64           `json:"count"`
	Users []userWithTeams `json:"users"`
}

type userWithTeams struct {
	Id             string   `json:"id"`
	Name           string   `json:"name"`
	Email          string   `json:"email"`
	TeamIds        string   `json:"teamIds"`
	TeamCount      int      `json:"teamCount"`
	TeamNames      []string `json:"teamNames"`
	AccountCount   int      `json:"accountCount"`
	AccountSources []string `json:"accountSources"`
}

type createUsersRequest struct {
	Users []user `json:"users"`
}

// ListUsers returns users with pagination support
// @Summary      List users
// @Description  GET /plugins/org/users?page=1&pageSize=50&email=example
// @Tags 		 plugins/org
// @Produce      json
// @Param        page      query  int     false  "page number (default 1)"
// @Param        pageSize  query  int     false  "page size (default 50)"
// @Param        email     query  string  false  "filter by email (case-insensitive, partial match)"
// @Success      200  {object}  paginatedUsers
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/users [get]
func (h *Handlers) ListUsers(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
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
	emailFilter := input.Query.Get("email")
	users, count, err := h.store.findUsersPaginated(page, pageSize, emailFilter)
	if err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{
		Body:   paginatedUsers{Count: count, Users: users},
		Status: http.StatusOK,
	}, nil
}

// GetUserById returns a single user by ID
// @Summary      Get a user
// @Description  get a user by ID
// @Tags 		 plugins/org
// @Produce      json
// @Param        userId  path  string  true  "user ID"
// @Success      200  {object}  user
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 404  {object} shared.ApiBody "Not Found"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/users/{userId} [get]
func (h *Handlers) GetUserById(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	userId := input.Params["userId"]
	if userId == "" {
		return nil, errors.BadInput.New("userId is required")
	}
	u, err := h.store.findUserById(userId)
	if err != nil {
		return nil, err
	}
	// fetch team associations
	var tus []crossdomain.TeamUser
	teamIds := ""
	tus, err = h.store.findTeamUsersByUserId(userId)
	if err == nil && len(tus) > 0 {
		var ids []string
		for _, tu := range tus {
			ids = append(ids, tu.TeamId)
		}
		teamIds = strings.Join(ids, ";")
	}
	return &plugin.ApiResourceOutput{
		Body: user{
			Id:      u.Id,
			Name:    u.Name,
			Email:   u.Email,
			TeamIds: teamIds,
		},
		Status: http.StatusOK,
	}, nil
}

// CreateUsers creates one or more users
// @Summary      Create users
// @Description  create one or more users
// @Tags 		 plugins/org
// @Accept       json
// @Produce      json
// @Param        body  body  createUsersRequest  true  "users to create"
// @Success      201  {array}  user
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/users [post]
func (h *Handlers) CreateUsers(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var req createUsersRequest
	err := helper.Decode(input.Body, &req, nil)
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid request body")
	}
	if len(req.Users) == 0 {
		return nil, errors.BadInput.New("at least one user is required")
	}
	var created []user
	for _, u := range req.Users {
		id := uuid.New().String()
		domainUser := &crossdomain.User{
			DomainEntity: domainlayer.DomainEntity{Id: id},
			Name:         u.Name,
			Email:        u.Email,
		}
		if err := h.store.createUser(domainUser); err != nil {
			return nil, err
		}
		u.Id = id
		created = append(created, u)
	}
	return &plugin.ApiResourceOutput{Body: created, Status: http.StatusCreated}, nil
}

// UpdateUserById updates a user by ID
// @Summary      Update a user
// @Description  update a user by ID
// @Tags 		 plugins/org
// @Accept       json
// @Produce      json
// @Param        userId  path  string  true  "user ID"
// @Param        body    body  user    true  "user fields to update"
// @Success      200  {object}  user
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 404  {object} shared.ApiBody "Not Found"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/users/{userId} [put]
func (h *Handlers) UpdateUserById(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	userId := input.Params["userId"]
	if userId == "" {
		return nil, errors.BadInput.New("userId is required")
	}
	existing, err := h.store.findUserById(userId)
	if err != nil {
		return nil, err
	}
	var u user
	if e := helper.Decode(input.Body, &u, nil); e != nil {
		return nil, errors.BadInput.Wrap(e, "invalid request body")
	}
	existing.Name = u.Name
	existing.Email = u.Email
	if err := h.store.updateUser(existing); err != nil {
		return nil, err
	}
	// replace team associations
	if err := h.store.replaceTeamUsersForUser(userId, strings.Split(u.TeamIds, ";")); err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{
		Body: user{
			Id:      existing.Id,
			Name:    existing.Name,
			Email:   existing.Email,
			TeamIds: u.TeamIds,
		},
		Status: http.StatusOK,
	}, nil
}

// DeleteUserById deletes a user by ID and its associated team_users and user_accounts
// @Summary      Delete a user
// @Description  delete a user by ID (cascades to team_users and user_accounts)
// @Tags 		 plugins/org
// @Param        userId  path  string  true  "user ID"
// @Success      200
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/users/{userId} [delete]
func (h *Handlers) DeleteUserById(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	userId := input.Params["userId"]
	if userId == "" {
		return nil, errors.BadInput.New("userId is required")
	}
	if err := h.store.deleteUser(userId); err != nil {
		return nil, err
	}
	return &plugin.ApiResourceOutput{Status: http.StatusOK}, nil
}

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
	"sort"

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	helper "github.com/apache/incubator-devlake/helpers/pluginhelper/api"
)

type teamUsersResponse struct {
	TeamId  string   `json:"teamId"`
	UserIds []string `json:"userIds"`
	Count   int      `json:"count"`
}

type updateTeamUsersRequest struct {
	UserIds []string `json:"userIds"`
}

type userTeamsResponse struct {
	UserId    string   `json:"userId"`
	TeamIds   []string `json:"teamIds"`
	TeamNames []string `json:"teamNames"`
	Count     int      `json:"count"`
}

type updateUserTeamsRequest struct {
	TeamIds []string `json:"teamIds"`
}

func (h *Handlers) sanitizeTeamUserIds(userIds []string) ([]string, errors.Error) {
	if len(userIds) == 0 {
		return []string{}, nil
	}

	uniqueUserIds := make([]string, 0, len(userIds))
	seen := make(map[string]struct{}, len(userIds))
	for _, userId := range userIds {
		if userId == "" {
			continue
		}
		if _, exists := seen[userId]; exists {
			continue
		}
		seen[userId] = struct{}{}
		uniqueUserIds = append(uniqueUserIds, userId)
	}

	if len(uniqueUserIds) == 0 {
		return []string{}, nil
	}

	users, err := h.store.findUsersByIds(uniqueUserIds)
	if err != nil {
		return nil, err
	}

	existingUserIds := make(map[string]struct{}, len(users))
	for _, u := range users {
		existingUserIds[u.Id] = struct{}{}
	}

	filteredUserIds := make([]string, 0, len(uniqueUserIds))
	for _, userId := range uniqueUserIds {
		if _, exists := existingUserIds[userId]; exists {
			filteredUserIds = append(filteredUserIds, userId)
		}
	}

	return filteredUserIds, nil
}

func (h *Handlers) listTeamUserIds(teamId string) ([]string, errors.Error) {
	teamUsers, err := h.store.findTeamUsersByTeamId(teamId)
	if err != nil {
		return nil, err
	}

	userIds := make([]string, 0, len(teamUsers))
	seen := make(map[string]struct{}, len(teamUsers))
	for _, teamUser := range teamUsers {
		if teamUser.UserId == "" {
			continue
		}
		if _, exists := seen[teamUser.UserId]; exists {
			continue
		}
		seen[teamUser.UserId] = struct{}{}
		userIds = append(userIds, teamUser.UserId)
	}
	sort.Strings(userIds)
	return userIds, nil
}

func (h *Handlers) sanitizeUserTeamIds(teamIds []string) ([]string, errors.Error) {
	if len(teamIds) == 0 {
		return []string{}, nil
	}

	uniqueTeamIds := make([]string, 0, len(teamIds))
	seen := make(map[string]struct{}, len(teamIds))
	for _, teamId := range teamIds {
		if teamId == "" {
			continue
		}
		if _, exists := seen[teamId]; exists {
			continue
		}
		seen[teamId] = struct{}{}
		uniqueTeamIds = append(uniqueTeamIds, teamId)
	}

	if len(uniqueTeamIds) == 0 {
		return []string{}, nil
	}

	teams, err := h.store.findTeamsByIds(uniqueTeamIds)
	if err != nil {
		return nil, err
	}

	existingTeamIds := make(map[string]struct{}, len(teams))
	for _, t := range teams {
		existingTeamIds[t.Id] = struct{}{}
	}

	filteredTeamIds := make([]string, 0, len(uniqueTeamIds))
	for _, teamId := range uniqueTeamIds {
		if _, exists := existingTeamIds[teamId]; exists {
			filteredTeamIds = append(filteredTeamIds, teamId)
		}
	}

	sort.Strings(filteredTeamIds)
	return filteredTeamIds, nil
}

func (h *Handlers) listUserTeamIds(userId string) ([]string, errors.Error) {
	teamUsers, err := h.store.findTeamUsersByUserId(userId)
	if err != nil {
		return nil, err
	}

	teamIds := make([]string, 0, len(teamUsers))
	seen := make(map[string]struct{}, len(teamUsers))
	for _, teamUser := range teamUsers {
		if teamUser.TeamId == "" {
			continue
		}
		if _, exists := seen[teamUser.TeamId]; exists {
			continue
		}
		seen[teamUser.TeamId] = struct{}{}
		teamIds = append(teamIds, teamUser.TeamId)
	}
	sort.Strings(teamIds)
	return teamIds, nil
}

func (h *Handlers) listUserTeamData(userId string) ([]string, []string, errors.Error) {
	teamIds, err := h.listUserTeamIds(userId)
	if err != nil {
		return nil, nil, err
	}
	if len(teamIds) == 0 {
		return []string{}, []string{}, nil
	}

	teams, err := h.store.findTeamsByIds(teamIds)
	if err != nil {
		return nil, nil, err
	}

	teamNameById := make(map[string]string, len(teams))
	for _, t := range teams {
		teamNameById[t.Id] = t.Name
	}

	teamNames := make([]string, 0, len(teamIds))
	for _, teamId := range teamIds {
		if teamName, exists := teamNameById[teamId]; exists && teamName != "" {
			teamNames = append(teamNames, teamName)
		}
	}
	sort.Strings(teamNames)

	return teamIds, teamNames, nil
}

// GetTeamUsersByTeamId returns all user IDs assigned to a team
// @Summary      List team users
// @Description  get all user IDs assigned to a team
// @Tags 			 plugins/org
// @Produce      json
// @Param        teamId  path  string  true  "team ID"
// @Success      200  {object}  teamUsersResponse
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 404  {object} shared.ApiBody "Not Found"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/teams/{teamId}/users [get]
func (h *Handlers) GetTeamUsersByTeamId(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	teamId := input.Params["teamId"]
	if teamId == "" {
		return nil, errors.BadInput.New("teamId is required")
	}

	if _, err := h.store.findTeamById(teamId); err != nil {
		return nil, err
	}

	userIds, err := h.listTeamUserIds(teamId)
	if err != nil {
		return nil, err
	}

	return &plugin.ApiResourceOutput{
		Body: teamUsersResponse{
			TeamId:  teamId,
			UserIds: userIds,
			Count:   len(userIds),
		},
		Status: http.StatusOK,
	}, nil
}

// UpdateTeamUsersByTeamId replaces user assignments for a team
// @Summary      Replace team users
// @Description  replace user assignments for a team by team ID
// @Tags 			 plugins/org
// @Accept       json
// @Produce      json
// @Param        teamId  path  string                   true  "team ID"
// @Param        body    body  updateTeamUsersRequest   true  "list of user IDs to assign to the team"
// @Success      200  {object}  teamUsersResponse
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 404  {object} shared.ApiBody "Not Found"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/teams/{teamId}/users [put]
func (h *Handlers) UpdateTeamUsersByTeamId(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	teamId := input.Params["teamId"]
	if teamId == "" {
		return nil, errors.BadInput.New("teamId is required")
	}

	if _, err := h.store.findTeamById(teamId); err != nil {
		return nil, err
	}

	var req updateTeamUsersRequest
	if err := helper.Decode(input.Body, &req, nil); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid request body")
	}
	validUserIds, err := h.sanitizeTeamUserIds(req.UserIds)
	if err != nil {
		return nil, err
	}

	if err := h.store.replaceTeamUsersForTeam(teamId, validUserIds); err != nil {
		return nil, err
	}

	userIds, err := h.listTeamUserIds(teamId)
	if err != nil {
		return nil, err
	}

	return &plugin.ApiResourceOutput{
		Body: teamUsersResponse{
			TeamId:  teamId,
			UserIds: userIds,
			Count:   len(userIds),
		},
		Status: http.StatusOK,
	}, nil
}

// GetUserTeamsByUserId returns all team IDs and names assigned to a user
// @Summary      List user teams
// @Description  get all team IDs and names assigned to a user
// @Tags 			 plugins/org
// @Produce      json
// @Param        userId  path  string  true  "user ID"
// @Success      200  {object}  userTeamsResponse
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 404  {object} shared.ApiBody "Not Found"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/users/{userId}/teams [get]
func (h *Handlers) GetUserTeamsByUserId(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	userId := input.Params["userId"]
	if userId == "" {
		return nil, errors.BadInput.New("userId is required")
	}

	if _, err := h.store.findUserById(userId); err != nil {
		return nil, err
	}

	teamIds, teamNames, err := h.listUserTeamData(userId)
	if err != nil {
		return nil, err
	}

	return &plugin.ApiResourceOutput{
		Body: userTeamsResponse{
			UserId:    userId,
			TeamIds:   teamIds,
			TeamNames: teamNames,
			Count:     len(teamIds),
		},
		Status: http.StatusOK,
	}, nil
}

// UpdateUserTeamsByUserId replaces team assignments for a user
// @Summary      Replace user teams
// @Description  replace team assignments for a user by user ID
// @Tags 			 plugins/org
// @Accept       json
// @Produce      json
// @Param        userId  path  string                  true  "user ID"
// @Param        body    body  updateUserTeamsRequest  true  "list of team IDs to assign to the user"
// @Success      200  {object}  userTeamsResponse
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 404  {object} shared.ApiBody "Not Found"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router       /plugins/org/users/{userId}/teams [put]
func (h *Handlers) UpdateUserTeamsByUserId(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	userId := input.Params["userId"]
	if userId == "" {
		return nil, errors.BadInput.New("userId is required")
	}

	if _, err := h.store.findUserById(userId); err != nil {
		return nil, err
	}

	var req updateUserTeamsRequest
	if err := helper.Decode(input.Body, &req, nil); err != nil {
		return nil, errors.BadInput.Wrap(err, "invalid request body")
	}

	validTeamIds, err := h.sanitizeUserTeamIds(req.TeamIds)
	if err != nil {
		return nil, err
	}

	if err := h.store.replaceTeamUsersForUser(userId, validTeamIds); err != nil {
		return nil, err
	}

	teamIds, teamNames, err := h.listUserTeamData(userId)
	if err != nil {
		return nil, err
	}

	return &plugin.ApiResourceOutput{
		Body: userTeamsResponse{
			UserId:    userId,
			TeamIds:   teamIds,
			TeamNames: teamNames,
			Count:     len(teamIds),
		},
		Status: http.StatusOK,
	}, nil
}

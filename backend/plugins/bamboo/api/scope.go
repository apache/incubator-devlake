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
	"net/http"
	"strconv"
	"time"

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/bamboo/models"
	"github.com/mitchellh/mapstructure"
)

type apiProject struct {
	models.BambooProject
	TransformationRuleName string `json:"transformationRuleName,omitempty"`
}

type req struct {
	Data []*models.BambooProject `json:"data"`
}

// PutScope create or update bamboo project
// @Summary create or update bamboo project
// @Description Create or update bamboo project
// @Tags plugins/bamboo
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scope body req true "json"
// @Success 200  {object} []models.BambooProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scopes [PUT]
func PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, _ := extractParam(input.Params)

	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}
	var projects req
	err := errors.Convert(mapstructure.Decode(input.Body, &projects))
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "decoding Bamboo project error")
	}
	keeper := make(map[string]struct{})
	now := time.Now()
	for _, project := range projects.Data {
		if _, ok := keeper[project.ProjectKey]; ok {
			return nil, errors.BadInput.New("duplicated item")
		} else {
			keeper[project.ProjectKey] = struct{}{}
		}
		project.ConnectionId = connectionId
		project.CreatedAt = now
		project.UpdatedAt = now
		err = verifyProject(project)
		if err != nil {
			return nil, err
		}
	}
	err = basicRes.GetDal().CreateOrUpdate(projects.Data)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving BambooProject")
	}
	return &plugin.ApiResourceOutput{Body: projects.Data, Status: http.StatusOK}, nil
}

// UpdateScope patch to bamboo project
// @Summary patch to bamboo project
// @Description patch to bamboo project
// @Tags plugins/bamboo
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param projectKey path int false "project ID"
// @Param scope body models.BambooProject true "json"
// @Success 200  {object} models.BambooProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scopes/{projectKey} [PATCH]
func UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, projectKey := extractParam(input.Params)
	if connectionId == 0 || projectKey == "" {
		return nil, errors.BadInput.New("invalid path params")
	}
	var project models.BambooProject
	err := basicRes.GetDal().First(&project, dal.Where("connection_id = ? AND project_key = ?", connectionId, projectKey))
	if err != nil {
		return nil, errors.Default.Wrap(err, "getting BambooProject error")
	}
	err = api.DecodeMapStruct(input.Body, &project)
	if err != nil {
		return nil, errors.Default.Wrap(err, "patch bamboo project error")
	}
	err = verifyProject(&project)
	if err != nil {
		return nil, err
	}
	err = basicRes.GetDal().Update(project)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving BambooProject")
	}
	return &plugin.ApiResourceOutput{Body: project, Status: http.StatusOK}, nil
}

// GetScopeList get Bamboo projects
// @Summary get Bamboo projects
// @Description get Bamboo projects
// @Tags plugins/bamboo
// @Param connectionId path int false "connection ID"
// @Success 200  {object} []apiProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scopes/ [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var projects []models.BambooProject
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	limit, offset := api.GetLimitOffset(input.Query, "pageSize", "page")
	err := basicRes.GetDal().All(&projects, dal.Where("connection_id = ?", connectionId), dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, err
	}
	var ruleIds []uint64
	for _, proj := range projects {
		if proj.TransformationRuleId > 0 {
			ruleIds = append(ruleIds, proj.TransformationRuleId)
		}
	}
	var rules []models.BambooTransformationRule
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
	var apiProjects []apiProject
	for _, proj := range projects {
		apiProjects = append(apiProjects, apiProject{proj, names[proj.TransformationRuleId]})
	}
	return &plugin.ApiResourceOutput{Body: apiProjects, Status: http.StatusOK}, nil
}

// GetScope get one Bamboo project
// @Summary get one Bamboo project
// @Description get one Bamboo project
// @Tags plugins/bamboo
// @Param connectionId path int false "connection ID"
// @Param projectKey path int false "project ID"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} apiProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/bamboo/connections/{connectionId}/scopes/{projectKey} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var project models.BambooProject
	connectionId, projectKey := extractParam(input.Params)
	if connectionId == 0 || projectKey == "" {
		return nil, errors.BadInput.New("invalid path params")
	}
	db := basicRes.GetDal()
	err := db.First(&project, dal.Where("connection_id = ? AND project_key = ?", connectionId, projectKey))
	if err != nil && db.IsErrorNotFound(err) {
		var scope models.BambooProject
		connection := &models.BambooConnection{}
		err = connectionHelper.First(connection, input.Params)
		if err != nil {
			return nil, err
		}
		apiClient, err := api.NewApiClientFromConnection(context.TODO(), basicRes, connection)
		if err != nil {
			return nil, err
		}

		apiProject, err := GetApiProject(projectKey, apiClient)
		if err != nil {
			return nil, err
		}

		scope.Convert(apiProject)
		scope.ConnectionId = connectionId
		err = db.CreateIfNotExist(&scope)
		if err != nil {
			return nil, err
		}
		return nil, errors.NotFound.New("record not found")
	} else if err != nil {
		return nil, err
	}

	var rule models.BambooTransformationRule
	if project.TransformationRuleId > 0 {
		err = basicRes.GetDal().First(&rule, dal.Where("id = ?", project.TransformationRuleId))
		if err != nil {
			return nil, err
		}
	}
	return &plugin.ApiResourceOutput{Body: apiProject{project, rule.Name}, Status: http.StatusOK}, nil
}

func extractParam(params map[string]string) (uint64, string) {
	connectionId, _ := strconv.ParseUint(params["connectionId"], 10, 64)
	projectKey := params["projectKey"]
	return connectionId, projectKey
}

func verifyProject(project *models.BambooProject) errors.Error {
	if project.ConnectionId == 0 {
		return errors.BadInput.New("invalid connectionId")
	}
	if project.ProjectKey == "" {
		return errors.BadInput.New("invalid projectKey")
	}
	return nil
}

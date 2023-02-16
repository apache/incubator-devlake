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

	"github.com/apache/incubator-devlake/core/dal"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/plugin"
	"github.com/apache/incubator-devlake/helpers/pluginhelper/api"
	"github.com/apache/incubator-devlake/plugins/sonarqube/models"
)

type req struct {
	Data []*models.SonarqubeProject `json:"data"`
}

// PutScope create or update sonarqube project
// @Summary create or update sonarqube project
// @Description Create or update sonarqube project
// @Tags plugins/sonarqube
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param scope body req true "json"
// @Success 200  {object} []models.SonarqubeProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId}/scopes [PUT]
func PutScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid connectionId")
	}
	var projects req
	// As we need to process *api.Iso8601Time, we need to use DecodeMapStruct instead of mapstructure.Decode
	err := errors.Convert(api.DecodeMapStruct(input.Body, &projects))
	if err != nil {
		return nil, errors.BadInput.Wrap(err, "decoding Sonarqube project error")
	}
	keeper := make(map[string]struct{})
	for _, project := range projects.Data {
		if _, ok := keeper[project.ProjectKey]; ok {
			return nil, errors.BadInput.New("duplicated item")
		} else {
			keeper[project.ProjectKey] = struct{}{}
		}
		project.ConnectionId = connectionId
		err = verifyProject(project)
		if err != nil {
			return nil, err
		}
	}
	err = basicRes.GetDal().CreateOrUpdate(projects.Data)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving SonarqubeProject")
	}
	return &plugin.ApiResourceOutput{Body: projects.Data, Status: http.StatusOK}, nil
}

// UpdateScope patch to sonarqube project
// @Summary patch to sonarqube project
// @Description patch to sonarqube project
// @Tags plugins/sonarqube
// @Accept application/json
// @Param connectionId path int false "connection ID"
// @Param projectKey path string false "project Key"
// @Param scope body models.SonarqubeProject true "json"
// @Success 200  {object} models.SonarqubeProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId}/scopes/{projectKey} [PATCH]
func UpdateScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	connectionId, projectKey := extractParam(input.Params)
	if connectionId*uint64(len(projectKey)) == 0 {
		return nil, errors.BadInput.New("invalid connectionId or projectKey")
	}
	var project models.SonarqubeProject
	err := basicRes.GetDal().First(&project, dal.Where("connection_id = ? AND project_key = ?", connectionId, projectKey))
	if err != nil {
		return nil, errors.Default.Wrap(err, "getting SonarqubeProject error")
	}
	err = api.DecodeMapStruct(input.Body, &project)
	if err != nil {
		return nil, errors.Default.Wrap(err, "patch sonarqube project error")
	}
	err = verifyProject(&project)
	if err != nil {
		return nil, err
	}
	err = basicRes.GetDal().Update(project)
	if err != nil {
		return nil, errors.Default.Wrap(err, "error on saving SonarqubeProject")
	}
	return &plugin.ApiResourceOutput{Body: project, Status: http.StatusOK}, nil
}

// GetScopeList get Sonarqube projects
// @Summary get Sonarqube projects
// @Description get Sonarqube projects
// @Tags plugins/sonarqube
// @Param connectionId path int false "connection ID"
// @Success 200  {object} []apiProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId}/scopes/ [GET]
func GetScopeList(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var projects []models.SonarqubeProject
	connectionId, _ := extractParam(input.Params)
	if connectionId == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	limit, offset := api.GetLimitOffset(input.Query, "pageSize", "page")
	err := basicRes.GetDal().All(&projects, dal.Where("connection_id = ?", connectionId), dal.Limit(limit), dal.Offset(offset))
	if err != nil {
		return nil, err
	}

	return &plugin.ApiResourceOutput{Body: projects, Status: http.StatusOK}, nil
}

// GetScope get one Sonarqube project
// @Summary get one Sonarqube project
// @Description get one Sonarqube project
// @Tags plugins/sonarqube
// @Param connectionId path int false "connection ID"
// @Param projectKey path string false "project key"
// @Param pageSize query int false "page size, default 50"
// @Param page query int false "page size, default 1"
// @Success 200  {object} apiProject
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /plugins/sonarqube/connections/{connectionId}/scopes/{projectKey} [GET]
func GetScope(input *plugin.ApiResourceInput) (*plugin.ApiResourceOutput, errors.Error) {
	var project models.SonarqubeProject
	connectionId, projectKey := extractParam(input.Params)
	if connectionId*uint64(len(projectKey)) == 0 {
		return nil, errors.BadInput.New("invalid path params")
	}
	db := basicRes.GetDal()
	err := db.First(&project, dal.Where("connection_id = ? AND project_key = ?", connectionId, projectKey))
	if db.IsErrorNotFound(err) {
		return nil, errors.NotFound.New("record not found")
	}
	if err != nil {
		return nil, err
	}

	return &plugin.ApiResourceOutput{Body: project, Status: http.StatusOK}, nil
}

func extractParam(params map[string]string) (uint64, string) {
	connectionId, _ := strconv.ParseUint(params["connectionId"], 10, 64)
	projectKey := params["projectKey"]
	return connectionId, projectKey
}

func verifyProject(project *models.SonarqubeProject) errors.Error {
	if project.ConnectionId == 0 {
		return errors.BadInput.New("invalid connectionId")
	}
	if len(project.ProjectKey) == 0 {
		return errors.BadInput.New("invalid project key")
	}
	return nil
}

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

package project

import (
	"fmt"
	"net/http"

	"github.com/apache/incubator-devlake/api/shared"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/services"
	"github.com/gin-gonic/gin"
)

type PaginatedProjects struct {
	Projects []*models.BaseProject `json:"projects"`
	Count    int64                 `json:"count"`
}

// @Summary Create and run a new project
// @Description Create and run a new project
// @Tags framework/projects
// @Accept application/json
// @Param projectName path string true "project name"
// @Success 200  {object} models.ApiOutputProject
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /projects/:projectName [get]
func GetProject(c *gin.Context) {
	projectOutput := &models.ApiOutputProject{}
	projectName := c.Param("projectName")

	project, err := services.GetProject(projectName)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error getting project"))
		return
	}

	projectOutput.BaseProject = project.BaseProject
	err = services.LoadBluePrintAndMetrics(projectOutput)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, fmt.Sprintf("Failed to LoadBluePrintAndMetrics on GetProject for %s", projectOutput.Name)))
		return
	}

	shared.ApiOutputSuccess(c, projectOutput, http.StatusOK)
}

// @Summary Get list of projects
// @Description GET /projects?page=1&pageSize=10
// @Tags framework/projects
// @Param page query int true "query"
// @Param pageSize query int true "query"
// @Success 200  {object} PaginatedProjects
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /projects [get]
func GetProjects(c *gin.Context) {
	var query services.ProjectQuery
	err := c.ShouldBindQuery(&query)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, shared.BadRequestBody))
		return
	}
	projects, count, err := services.GetProjects(&query)
	if err != nil {
		shared.ApiOutputAbort(c, errors.Default.Wrap(err, "error getting projects"))
		return
	}

	baseProjects := make([]*models.BaseProject, count)
	for i, project := range projects {
		baseProjects[i] = &project.BaseProject
	}

	shared.ApiOutputSuccess(c, PaginatedProjects{
		Projects: baseProjects,
		Count:    count,
	}, http.StatusOK)
}

// @Summary Create a new project
// @Description Create a new project
// @Tags framework/projects
// @Accept application/json
// @Param project body models.ApiInputProject true "json"
// @Success 200  {object} models.ApiOutputProject
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /projects [post]
func PostProject(c *gin.Context) {
	projectInput := &models.ApiInputProject{}
	projectOutput := &models.ApiOutputProject{}

	err := c.ShouldBind(projectInput)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, shared.BadRequestBody))
		return
	}

	err = services.CreateProject(&models.Project{BaseProject: projectInput.BaseProject})
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, "error creating project"))
		return
	}

	// check if need to changed the blueprint setting
	if projectInput.Enable != nil {
		_, err = services.PatchBlueprintEnableByProjectName(projectInput.Name, *projectInput.Enable)
		if err != nil {
			shared.ApiOutputError(c, errors.BadInput.Wrap(err, "Failed to set if project enable"))
			return
		}
	}

	// check if need flush the Metrics
	if projectInput.Metrics != nil {
		err = services.FlushProjectMetrics(projectInput.Name, projectInput.Metrics)
		if err != nil {
			shared.ApiOutputError(c, errors.BadInput.Wrap(err, "Failed to flush project metrics"))
			return
		}
	}

	projectOutput.BaseProject = projectInput.BaseProject
	err = services.LoadBluePrintAndMetrics(projectOutput)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, fmt.Sprintf("Failed to LoadBluePrintAndMetrics on PostProject for %s", projectOutput.Name)))
		return
	}

	shared.ApiOutputSuccess(c, projectOutput, http.StatusCreated)
}

// @Summary Patch a project
// @Description Patch a project
// @Tags framework/projects
// @Accept application/json
// @Param project body models.ApiInputProject true "json"
// @Success 200  {object} models.ApiOutputProject
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /projects/:projectName [patch]
func PatchProject(c *gin.Context) {
	projectName := c.Param("projectName")

	var body map[string]interface{}
	err := c.ShouldBind(&body)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, shared.BadRequestBody))
		return
	}

	projectOutput, err := services.PatchProject(projectName, body)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, "error patch project"))
		return
	}

	shared.ApiOutputSuccess(c, projectOutput, http.StatusCreated)
}

// @Cancel a project
// @Description Cancel a project
// @Tags framework/projects
// @Success 200
// @Failure 400  {string} er2rcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /projects/:projectName [delete]
//func DeleteProject(c *gin.Context) {
//}

// @Summary Get a ProjectMetrics
// @Description Get a ProjectMetrics
// @Tags framework/ProjectMetrics
// @Param projectName path string true "project name"
// @Param pluginName path string true "plugin name"
// @Success 200  {object} models.BaseProjectMetric
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /projects/:projectName/metrics/:pluginName [get]
func GetProjectMetrics(c *gin.Context) {
	projectName := c.Param("projectName")
	pluginName := c.Param("pluginName")

	projectMetric, err := services.GetProjectMetric(projectName, pluginName)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, "error getting project metric"))
		return
	}

	shared.ApiOutputSuccess(c, projectMetric.BaseProjectMetric, http.StatusOK)
}

// @Summary Create a new ProjectMetrics
// @Description Create  a new ProjectMetrics
// @Tags framework/ProjectMetrics
// @Accept application/json
// @Param project body models.BaseProjectMetric true "json"
// @Success 200  {object} models.BaseProjectMetric
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /projects/:projectName/metrics [post]
func PostProjectMetrics(c *gin.Context) {
	projectMetric := &models.BaseProjectMetric{}

	projectName := c.Param("projectName")

	_, err1 := services.GetProject(projectName)
	if err1 != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err1, shared.BadRequestBody))
		return
	}

	err := c.ShouldBind(projectMetric)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, shared.BadRequestBody))
		return
	}

	projectMetric.ProjectName = projectName
	err = services.CreateProjectMetric(&models.ProjectMetric{BaseProjectMetric: *projectMetric})
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, "error creating project metric"))
		return
	}

	shared.ApiOutputSuccess(c, projectMetric, http.StatusCreated)
}

// @Summary Patch a ProjectMetrics
// @Description Patch a ProjectMetrics
// @Tags framework/ProjectMetrics
// @Accept application/json
// @Param ProjectMetrics body models.BaseProjectMetric true "json"
// @Success 200  {object} models.BaseProjectMetric
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /projects/:projectName/metrics/:pluginName  [patch]
func PatchProjectMetrics(c *gin.Context) {
	projectName := c.Param("projectName")
	pluginName := c.Param("pluginName")

	var body map[string]interface{}
	err := c.ShouldBind(&body)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, shared.BadRequestBody))
		return
	}

	projectMetric, err := services.PatchProjectMetric(projectName, pluginName, body)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, "error patch project metric"))
		return
	}

	shared.ApiOutputSuccess(c, projectMetric.BaseProjectMetric, http.StatusCreated)
}

// @delete a ProjectMetrics
// @Description delete a ProjectMetrics
// @Tags framework/ProjectMetrics
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /project_metrics/:projectName/:pluginName [delete]
//func DeleteProjectMetrics(c *gin.Context) {
//}

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
	"net/http"

	"github.com/apache/incubator-devlake/api/shared"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/services"
	"github.com/gin-gonic/gin"
)

// @Summary Create and run a new project
// @Description Create and run a new project
// @Tags framework/projects
// @Accept application/json
// @Param project body models.Project true "json"
// @Success 200  {object} models.Project
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /projects/:projectName [get]
func GetProject(c *gin.Context) {
	projectName := c.Param("projectName")

	project, err := services.GetProject(projectName)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error getting project"))
		return
	}

	shared.ApiOutputSuccess(c, project, http.StatusOK)
}

// @Summary Get list of projects
// @Description GET /projects?page=1&pagesize=10
// @Tags framework/projects
// @Param page query int true "query"
// @Param pagesize query int true "query"
// @Success 200  {object} gin.H
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
	shared.ApiOutputSuccess(c, gin.H{"project": projects, "count": count}, http.StatusOK)
}

// @Summary Create a new project
// @Description Create a new project
// @Tags framework/projects
// @Accept application/json
// @Param project body models.Project true "json"
// @Success 200  {object} models.Project
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /projects [post]
func PostProject(c *gin.Context) {
	project := &models.Project{}

	err := c.ShouldBind(project)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, shared.BadRequestBody))
		return
	}

	err = services.CreateProject(project)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error creating project"))
		return
	}

	shared.ApiOutputSuccess(c, project, http.StatusCreated)
}

// @Summary Patch a project
// @Description Patch a project
// @Tags framework/projects
// @Accept application/json
// @Param project body models.Project true "json"
// @Success 200  {object} models.Project
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

	project, err := services.PatchProject(projectName, body)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error patch project"))
		return
	}

	shared.ApiOutputSuccess(c, project, http.StatusCreated)
}

// @Summary Get a ProjectMetrics
// @Description Get a ProjectMetrics
// @Tags framework/ProjectMetrics
// @Param page query int true "query"
// @Param pagesize query int true "query"
// @Success 200  {object} models.ProjectMetric
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /project_metrics/:projectName/:pluginName [get]
func GetProjectMetric(c *gin.Context) {
	projectName := c.Param("projectName")
	pluginName := c.Param("pluginName")

	projectMetric, err := services.GetProjectMetric(projectName, pluginName)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error getting project metric"))
		return
	}

	shared.ApiOutputSuccess(c, projectMetric, http.StatusOK)
}

// @Summary Create a new ProjectMetrics
// @Description Create  a new ProjectMetrics
// @Tags framework/ProjectMetrics
// @Accept application/json
// @Param project body models.Project true "json"
// @Success 200  {object} models.ProjectMetric
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /project_metrics [post]
func PostProjectMetrics(c *gin.Context) {
	projectMetric := &models.ProjectMetric{}

	err := c.ShouldBind(projectMetric)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, shared.BadRequestBody))
		return
	}

	err = services.CreateProjectMetric(projectMetric)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error creating project"))
		return
	}

	shared.ApiOutputSuccess(c, projectMetric, http.StatusCreated)
}

// @Summary Patch a ProjectMetrics
// @Description Patch a ProjectMetrics
// @Tags framework/ProjectMetrics
// @Accept application/json
// @Param ProjectMetrics body models.ProjectMetric true "json"
// @Success 200  {object} models.ProjectMetric
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /project_metrics/:projectName/:pluginName  [patch]
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
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error patch project"))
		return
	}

	shared.ApiOutputSuccess(c, projectMetric, http.StatusCreated)
}

/*
// @Cancel a pending ProjectMetrics
// @Description Cancel a pending ProjectMetrics
// @Tags framework/ProjectMetrics
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /project_metrics/:projectName/:pluginName [delete]
func Delete(c *gin.Context) {
	projectName := c.Param("projectName")
	pluginName := c.Param("pluginName")

	err := services.CancelProjectMetrics(id)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error cancelling pipeline"))
		return
	}
	shared.ApiOutputSuccess(c, nil, http.StatusOK)
}
*/

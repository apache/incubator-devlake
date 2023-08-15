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

	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/server/api/shared"
	"github.com/apache/incubator-devlake/server/services"

	"github.com/gin-gonic/gin"
)

type PaginatedProjects struct {
	Projects []*models.ApiOutputProject `json:"projects"`
	Count    int64                      `json:"count"`
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
	projectName := c.Param("projectName")[1:]

	projectOutput, err := services.GetProject(projectName)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error getting project"))
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
// @Failure 500  {string} errcode.Error "Internal Error"
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

	shared.ApiOutputSuccess(c, PaginatedProjects{
		Projects: projects,
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
	err := c.ShouldBind(projectInput)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, shared.BadRequestBody))
		return
	}

	projectOutput, err := services.CreateProject(projectInput)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error creating project"))
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
	projectName := c.Param("projectName")[1:]

	var body map[string]interface{}
	err := c.ShouldBind(&body)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, shared.BadRequestBody))
		return
	}

	projectOutput, err := services.PatchProject(projectName, body)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error patch project"))
		return
	}

	shared.ApiOutputSuccess(c, projectOutput, http.StatusCreated)
}

// @Summary Delete a project
// @Description Delete a project
// @Tags framework/projects
// @Accept application/json
// @Success 200
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internal Error"
// @Router /projects/:projectName [delete]
func DeleteProject(c *gin.Context) {
	projectName := c.Param("projectName")[1:]
	err := services.DeleteProject(projectName)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error deleting project"))
		return
	}
	shared.ApiOutputSuccess(c, nil, http.StatusOK)
}

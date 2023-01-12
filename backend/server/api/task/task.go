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

package task

import (
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/core/models"
	"github.com/apache/incubator-devlake/server/api/shared"
	"github.com/apache/incubator-devlake/server/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Delete(c *gin.Context) {
	taskId := c.Param("taskId")
	id, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, "invalid task ID format"))
		return
	}
	err = services.CancelTask(id)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error cancelling task"))
		return
	}
	shared.ApiOutputSuccess(c, nil, http.StatusOK)
}

type getTaskResponse struct {
	Tasks []*models.Task `json:"tasks"`
	Count int            `json:"count"`
}

// GetTaskByPipeline return most recent tasks
// @Summary Get tasks, only the most recent tasks will be returned
// @Tags framework/tasks
// @Accept application/json
// @Param pipelineId path int true "pipelineId"
// @Success 200  {object} getTaskResponse
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /pipelines/{pipelineId}/tasks [get]
func GetTaskByPipeline(c *gin.Context) {
	pipelineId, err := strconv.ParseUint(c.Param("pipelineId"), 10, 64)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, "invalid pipeline ID format"))
		return
	}
	tasks, err := services.GetTasksWithLastStatus(pipelineId)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error getting tasks"))
		return
	}
	shared.ApiOutputSuccess(c, getTaskResponse{Tasks: tasks, Count: len(tasks)}, http.StatusOK)
}

// RerunTask rerun the specified task.
// @Summary rerun task
// @Tags framework/tasks
// @Accept application/json
// @Success 200  {object} models.Task
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /tasks/{taskId}/rerun [post]
func PostRerun(c *gin.Context) {
	taskId := c.Param("taskId")
	id, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, "bad taskId format supplied"))
		return
	}
	task, err := services.RerunTask(id)
	if err != nil {
		shared.ApiOutputError(c, err)
		return
	}
	shared.ApiOutputSuccess(c, task, http.StatusOK)
}

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
	"net/http"
	"strconv"

	"github.com/apache/incubator-devlake/api/shared"
	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models"
	"github.com/apache/incubator-devlake/services"
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
	Tasks []models.Task `json:"tasks"`
	Count int           `json:"count"`
}

// GetTaskByPipeline return most recent tasks
// @Summary Get tasks, only the most recent tasks will be returned
// @Tags framework/task
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

type rerunRequest struct {
	TaskId uint64 `json:"taskId"`
}

// RerunTask rerun the specified the task. If taskId is 0, all failed tasks of this pipeline will rerun
// @Summary rerun tasks
// @Tags framework/task
// @Accept application/json
// @Param pipelineId path int true "pipelineId"
// @Param request body rerunRequest false "specify the task to rerun. If it's 0, all failed tasks of this pipeline will rerun"
// @Success 200  {object} shared.ApiBody
// @Failure 400  {object} shared.ApiBody "Bad Request"
// @Failure 500  {object} shared.ApiBody "Internal Error"
// @Router /pipelines/{pipelineId}/tasks [post]
func RerunTask(c *gin.Context) {
	var request rerunRequest
	err := c.BindJSON(&request)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, "invalid task ID format"))
		return
	}
	pipelineId, err := strconv.ParseUint(c.Param("pipelineId"), 10, 64)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, "invalid pipeline ID format"))
		return
	}
	pipeline, err := services.GetPipeline(pipelineId)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error get pipeline"))
		return
	}
	if pipeline.Status == models.TASK_RUNNING {
		shared.ApiOutputError(c, errors.BadInput.New("pipeline is running"))
		return
	}
	if pipeline.Status == models.TASK_CREATED || pipeline.Status == models.TASK_RERUN {
		shared.ApiOutputError(c, errors.BadInput.New("pipeline is waiting to run"))
		return
	}

	var failedTasks []models.Task
	if request.TaskId > 0 {
		failedTask, err := services.GetTask(request.TaskId)
		if err != nil || failedTask == nil {
			shared.ApiOutputError(c, errors.Default.Wrap(err, "error getting failed task"))
			return
		}
		if failedTask.PipelineId != pipelineId {
			shared.ApiOutputError(c, errors.BadInput.New("the task ID and pipeline ID doesn't match"))
			return
		}
		failedTasks = append(failedTasks, *failedTask)
	} else {
		tasks, err := services.GetTasksWithLastStatus(pipelineId)
		if err != nil {
			shared.ApiOutputError(c, errors.Default.Wrap(err, "error getting tasks"))
			return
		}
		for _, task := range tasks {
			if task.Status == models.TASK_FAILED {
				failedTasks = append(failedTasks, task)
			}
		}
	}
	if len(failedTasks) == 0 {
		shared.ApiOutputSuccess(c, nil, http.StatusOK)
		return
	}
	err = services.DeleteCreatedTasks(pipelineId)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error delete tasks"))
		return
	}
	_, err = services.SpawnTasks(failedTasks)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error create tasks"))
		return
	}
	err = services.UpdateDbPipelineStatus(pipelineId, models.TASK_RERUN)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error create tasks"))
		return
	}
	shared.ApiOutputSuccess(c, nil, http.StatusOK)
}

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
	"github.com/apache/incubator-devlake/services"
	"github.com/gin-gonic/gin"
)

/*
Get list of pipelines
GET /pipelines/pipeline:id/tasks?status=TASK_RUNNING&pending=1&page=1&=pagesize=10
{
	"tasks": [
		{"id": 1, "plugin": "", ...}
	],
	"count": 5
}
*/
// @Summary Get task
// @Description get task
// @Description SAMPLE
// @Description {
// @Description 	"tasks": [
// @Description 		{"id": 1, "plugin": "", ...}
// @Description 	],
// @Description 	"count": 5
// @Description }
// @Tags framework/task
// @Accept application/json
// @Param pipelineId path int true "pipelineId"
// @Success 200  {string} gin.H "{"tasks": tasks, "count": count}"
// @Failure 400  {string} errcode.Error "Bad Request"
// @Failure 500  {string} errcode.Error "Internel Error"
// @Router /pipelines/{pipelineId}/tasks [get]
func Index(c *gin.Context) {
	var query services.TaskQuery
	err := c.ShouldBindQuery(&query)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, shared.BadRequestBody))
		return
	}
	err = c.ShouldBindUri(&query)
	if err != nil {
		shared.ApiOutputError(c, errors.BadInput.Wrap(err, "bad request URI format"))
		return
	}
	tasks, count, err := services.GetTasks(&query)
	if err != nil {
		shared.ApiOutputError(c, errors.Default.Wrap(err, "error getting tasks"))
		return
	}
	shared.ApiOutputSuccess(c, gin.H{"tasks": tasks, "count": count}, http.StatusOK)
}

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

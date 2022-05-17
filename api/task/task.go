package task

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/api/shared"
	"github.com/merico-dev/lake/services"
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
// @Tags task
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
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	err = c.ShouldBindUri(&query)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	tasks, count, err := services.GetTasks(&query)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	shared.ApiOutputSuccess(c, gin.H{"tasks": tasks, "count": count}, http.StatusOK)
}

func Delete(c *gin.Context) {
	taskId := c.Param("taskId")
	id, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid task id")
		return
	}
	err = services.CancelTask(id)
	if err != nil {
		shared.ApiOutputError(c, err, http.StatusBadRequest)
		return
	}
	shared.ApiOutputSuccess(c, nil, http.StatusOK)
}

package task

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
func Index(ctx *gin.Context) {
	var query services.TaskQuery
	err := ctx.BindQuery(&query)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	err = ctx.BindUri(&query)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	tasks, count, err := services.GetTasks(&query)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks, "count": count})
}

func Delete(ctx *gin.Context) {
	taskId := ctx.Param("taskId")
	id, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid task id")
		return
	}
	err = services.CancelTask(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	}
}

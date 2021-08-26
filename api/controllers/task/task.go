package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/merico-dev/lake/api/services"
	"github.com/merico-dev/lake/api/types"
	"github.com/merico-dev/lake/logger"
)

func Post(ctx *gin.Context) {
	var data types.CreateTask
	err := ctx.MustBindWith(&data, binding.JSON)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	logger.Debug("Create Task", data)
	task, err := services.NewTask(data)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, task)
	// TODO: trigger plugin task
}

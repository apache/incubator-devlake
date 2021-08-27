package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/merico-dev/lake/services"
)

func Post(ctx *gin.Context) {
	var data services.NewTask
	err := ctx.MustBindWith(&data, binding.JSON)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	task, err := services.CreateTask(data)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusCreated, task)
}

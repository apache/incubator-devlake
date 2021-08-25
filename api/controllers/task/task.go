package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/merico-dev/lake/api/services"
	"github.com/merico-dev/lake/api/types"
	"github.com/merico-dev/lake/logger"
)

// PostTask godoc
// @Summary create a plugin task
// @Description create and trigger a plugin task
// @ID create-task
// @Accept  json
// @Produce  json
// @Param param body types.CreateTask true "task info"
// @Success 200 {object} models.Task
// @Header 200 {string} Token "qwerty"
// @Router /task [post]
func Post(ctx *gin.Context) {
	var data types.CreateTask
	err := ctx.MustBindWith(&data, binding.JSON)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	logger.Debug("CreateTask", data)
	task, err := services.NewTask(data)
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, task)
	// TODO: trigger plugin task
}

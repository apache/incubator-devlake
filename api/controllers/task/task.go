package task

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/merico-dev/lake/api/services"
	"github.com/merico-dev/lake/api/types"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins"
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

	// trigger plugins
	data.Options["ID"] = task.ID

	go func() {
		progress := make(chan float32)
		_ = plugins.RunPlugin(task.Plugin, data.Options, progress)
		for p := range progress {
			fmt.Printf("running plugin %v, progress: %v\n", task.Plugin, p*100)
		}
	}()

	ctx.JSON(http.StatusCreated, task)
}

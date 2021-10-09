package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/services"
)

func Post(ctx *gin.Context) {
	// We use a 2D array because the request body must be an array of a set of tasks
	// to be executed concurrently, while each set is to be executed sequentially.
	var data [][]services.NewTask

	err := ctx.MustBindWith(&data, binding.JSON)
	if err != nil {
		logger.Error("", err)
		ctx.JSON(http.StatusBadRequest, "You must send down an array of objects")
		return
	}

	tasks := services.CreateTasksInDBFromJSON(data)
	// Return all created tasks to the User
	ctx.JSON(http.StatusCreated, tasks)

	go func() {
		err := services.RunAllTasks(data, tasks)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}
	}()
}

func Get(ctx *gin.Context) {
	tasks, err := services.GetTasks(ctx.Query("status"))
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

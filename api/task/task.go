package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
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

	var tasks []models.Task

	// This double for loop executes each set of tasks sequentially while
	// executing the set of tasks concurrently.
	for _, array := range data {
		taskComplete := make(chan bool)
		count := 0
		for _, element := range array {
			task, err := services.CreateTask(element, taskComplete)
			if err != nil {
				_ = ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			tasks = append(tasks, *task)
		}
		for range taskComplete {
			count++
			if count == len(array) {
				close(taskComplete)
			}
		}
	}

	ctx.JSON(http.StatusCreated, tasks)
}

func Get(ctx *gin.Context) {
	tasks, err := services.GetTasks(ctx.Query("status"))
	if err != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

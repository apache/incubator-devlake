package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/services"
	"github.com/nats-io/nats.go"
)

func CancelTask(ctx *gin.Context) {
	// pull the task name from the query params
	taskName := ctx.Query("taskName")

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		logger.Error("err connecting to nats server", err)
	}

	// Simple Publisher
	nc.Publish(taskName, []byte("Hello World"))
}

func Post(ctx *gin.Context) {
	var data []services.NewTask

	err := ctx.MustBindWith(&data, binding.JSON)
	if err != nil {
		logger.Error("", err)
		ctx.JSON(http.StatusBadRequest, "You must send down an array of objects")
		return
	}

	var tasks []models.Task

	for _, value := range data {
		task, err := services.CreateTask(value)
		if err != nil {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		tasks = append(tasks, *task)
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

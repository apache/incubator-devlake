package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/services"
)

func CreateTasksInDBFromJSON(data [][]services.NewTask) []models.Task {
	// create all the tasks in the db without running the tasks
	var tasks []models.Task

	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			task, _ := services.CreateTaskInDB(data[i][j])
			tasks = append(tasks, *task)
		}
	}

	return tasks
}

func RunAllTasks(data [][]services.NewTask, ctx *gin.Context) {
	// This double for loop executes each set of tasks sequentially while
	// executing the set of tasks concurrently.
	for _, array := range data {
		taskComplete := make(chan bool)
		count := 0
		for _, taskFromRequest := range array {
			_, err := services.RunTask(taskFromRequest, taskComplete)
			if err != nil {
				_ = ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}
		for range taskComplete {
			count++
			if count == len(array) {
				close(taskComplete)
			}
		}
	}
}

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

	tasks := CreateTasksInDBFromJSON(data)
	// Return all created tasks to the User
	ctx.JSON(http.StatusCreated, tasks)

	go func() {
		RunAllTasks(data, ctx)
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

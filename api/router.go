package api

import (
	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/api/task"
)

func RegisterRouter(r *gin.Engine) {
	r.POST("/task", task.Post)
}

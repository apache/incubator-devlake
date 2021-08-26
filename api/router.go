package api

import (
	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/api/controllers/source"
	"github.com/merico-dev/lake/api/controllers/task"
)

func RegisterRouter(r *gin.Engine) {
	r.POST("/source", source.Post)
	r.POST("/task", task.Post)
}

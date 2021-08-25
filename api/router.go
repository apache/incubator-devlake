package api

import (
	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/api/controllers/source"
)

func RegisterRouter(r *gin.Engine) {
	r.POST("/source", source.Post)
}

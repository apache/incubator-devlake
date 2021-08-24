package api

import (
	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/api/controllers/source"
	"github.com/merico-dev/lake/config"
)

func CreateApiService() {
	gin.SetMode(config.V.GetString("MODE"))
	r := gin.Default()
	r.POST("/source", source.Post)
	r.Run(config.V.GetString("PORT"))
}

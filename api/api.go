package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/config"
)

func CreateApiService() {
	gin.SetMode(config.V.GetString("MODE"))
	router := gin.Default()

	// CORS CONFIG
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           96 * time.Hour,
	}))

	RegisterRouter(router)
	err := router.Run(config.V.GetString("PORT"))
	if err != nil {
		panic(err)
	}
}

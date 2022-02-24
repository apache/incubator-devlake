package api

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/config"
)

func CreateApiService() {
	v := config.GetConfig()
	gin.SetMode(v.GetString("MODE"))
	router := gin.Default()

	// CORS CONFIG
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           120 * time.Hour,
	}))

	RegisterRouter(router)
	err := router.Run(v.GetString("PORT"))
	if err != nil {
		panic(err)
	}
}

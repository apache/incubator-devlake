package main

import (
	_ "net/http"

	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/api"
	"github.com/merico-dev/lake/config"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/merico-dev/lake/api/docs"
)

func CreateApiService() {
	gin.SetMode(config.V.GetString("MODE"))
	r := gin.Default()

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	api.RegisterRouter(r)
	err := r.Run(config.V.GetString("PORT"))
	if err != nil {
		panic(err)
	}
}

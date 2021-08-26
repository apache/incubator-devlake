package main

import (
	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/api"
	"github.com/merico-dev/lake/config"
)

func CreateApiService() {
	gin.SetMode(config.V.GetString("MODE"))
	r := gin.Default()
	api.RegisterRouter(r)
	err := r.Run(config.V.GetString("PORT"))
	if err != nil {
		panic(err)
	}
}

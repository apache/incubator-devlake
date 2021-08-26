package main

import (
	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/api"
	"github.com/spf13/viper"
)

func CreateApiService() {
	gin.SetMode(viper.GetString("MODE"))
	r := gin.Default()
	api.RegisterRouter(r)
	err := r.Run(viper.GetString("PORT"))
	if err != nil {
		panic(err)
	}
}

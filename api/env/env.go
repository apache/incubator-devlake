package env

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/spf13/viper"
)

var V *viper.Viper

func Get(ctx *gin.Context) {
	configJson, err := config.GetConfigJson()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Your json is malformed")
		return
	}
	ctx.JSON(http.StatusOK, configJson)
}

func Set(ctx *gin.Context) {
	var envVars map[string]interface{}
	err := ctx.BindJSON(&envVars)
	if err != nil {
		logger.Error("", err)
		ctx.JSON(http.StatusBadRequest, "Your json is malformed")
		return
	}

	V := config.LoadConfigFile()

	for key, value := range envVars {
		V.Set(key, value)
	}

	err = V.WriteConfig()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "Could not write config file")
		return
	}

	ctx.JSON(http.StatusOK, envVars)
}

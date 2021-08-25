package config

import (
	"github.com/merico-dev/lake/logger"
	"github.com/spf13/viper"
)

var V *viper.Viper

func init() {
	V = viper.New()
	V.SetConfigFile(".env")
	V.AutomaticEnv()
	err := V.ReadInConfig()
	if err != nil {
		logger.Error(err)
	}
	V.SetDefault("PORT", ":8080")
}

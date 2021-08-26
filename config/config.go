package config

import (
	"github.com/merico-dev/lake/logger"
	"github.com/spf13/viper"
)

func ReadConfig() {
	logger.Info("loading config", true)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// TODO: dirname path
	viper.AddConfigPath("/Users/jonathanodonnell/go/src/github.com/merico-dev/lake/")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("failed to read in config", err)
	}
	viper.SetDefault("PORT", ":8080")
}

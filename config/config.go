package config

import (
	"fmt"
	"os"

	"github.com/merico-dev/lake/logger"
	"github.com/spf13/viper"
)

func getDirName() string {
	dirname, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current directory: %v\n", dirname)

	fmt.Printf("Name of ../../: %v\n", dirname)
	return dirname
}

func ReadConfig() {
	logger.Info("loading config", true)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(getDirName())
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("failed to read in config", err)
	}
	viper.SetDefault("PORT", ":8080")
}

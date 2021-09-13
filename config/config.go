package config

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
	"os"
	"path"
)

var V *viper.Viper

func init() {
	V = viper.New()
	configFile := os.Getenv("ENV_FILE")
	if configFile != "" {
		if !path.IsAbs(configFile) {
			panic("Please set ENV_FILE with absolute path. " +
				"Currently it should only be used for go test to load ENVs.")
		}
		V.SetConfigFile(configFile)
		V.Set("WORKING_DIRECTORY", path.Dir(configFile))
	} else {
		V.SetConfigName(".env")
		V.SetConfigType("env")

		V.AddConfigPath(".")
		V.AddConfigPath("conf")
		V.AddConfigPath("etc")

		execPath, execErr := os.Executable()
		if execErr == nil {
			V.AddConfigPath(path.Dir(execPath))
			V.AddConfigPath(path.Join(path.Dir(execPath), "conf"))
			V.AddConfigPath(path.Join(path.Dir(execPath), "etc"))
		}

		wdPath, _ := os.Getwd()
		V.Set("WORKING_DIRECTORY", wdPath)
	}

	_ = V.ReadInConfig()
	V.AutomaticEnv()
	V.SetDefault("PORT", ":8080")
}

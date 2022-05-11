package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Lowcase V for private this. You can use it by call GetConfig.
var v *viper.Viper = nil

func GetConfig() *viper.Viper {
	return v
}

// Set default value for no .env or .env not set it
func setDefaultValue() {
	v.SetDefault("DB_URL", "mysql://merico:merico@mysql:3306/lake?charset=utf8mb4&parseTime=True")
	v.SetDefault("PORT", ":8080")
	v.SetDefault("PLUGIN_DIR", "bin/plugins")
	v.SetDefault("TEMPORAL_TASK_QUEUE", "DEVLAKE_TASK_QUEUE")
	v.SetDefault("GITLAB_ENDPOINT", "https://gitlab.com/api/v4/")
	v.SetDefault("GITHUB_ENDPOINT", "https://api.github.com/")
}

func init() {
	// create the object and load the .env file
	v = viper.New()
	v.SetConfigFile(".env")
	err := v.ReadInConfig()
	if err != nil {
		logrus.Warn("Failed to read [.env] file:", err)
	}
	v.AutomaticEnv()

	setDefaultValue()
	// This line is essential for reading and writing
	v.WatchConfig()
}

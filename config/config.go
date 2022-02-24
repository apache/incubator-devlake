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
	v.SetDefault("PORT", ":8080")
	v.SetDefault("PLUGIN_DIR", "bin/plugins")
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

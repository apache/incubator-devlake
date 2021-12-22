package config

import (
	"github.com/spf13/viper"
)

var V *viper.Viper

func LoadConfigFile() *viper.Viper {
	V = viper.New()
	V.SetConfigFile(".env")
	_ = V.ReadInConfig()
	V.AutomaticEnv()
	return V
}

func init() {
	V := LoadConfigFile()
	V.SetDefault("PORT", ":8080")
	V.SetDefault("PLUGIN_DIR", "bin/plugins")
	// This line is essential for reading and writing
	V.WatchConfig()
}

func GetConfigJson() (map[string]string, error) {
	var configJson map[string]string
	err := V.Unmarshal(&configJson)
	if err != nil {
		return nil, err
	}
	return configJson, nil
}

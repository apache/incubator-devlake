package config

import (
	"github.com/spf13/viper"
)

var V *viper.Viper

func init() {
	V = viper.New()
	V.AutomaticEnv()
	V.SetDefault("PORT", ":8080")
}

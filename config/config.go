package config

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"
)

var V *viper.Viper

func init() {
	V = viper.New()
	V.AutomaticEnv()
	V.SetDefault("PORT", ":8080")
}

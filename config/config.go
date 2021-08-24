package config

import "github.com/spf13/viper"

var V *viper.Viper

func init() {
	V = viper.New()
	V.SetConfigFile(".env")
	err := V.ReadInConfig()
	if err != nil {
		panic(err)
	}
	V.SetDefault("PORT", ":8080")
}

package e2e

import "github.com/spf13/viper"

func LoadConfigFile() *viper.Viper {
	V := viper.New()
	V.SetConfigFile("../.env")
	_ = V.ReadInConfig()
	V.AutomaticEnv()
	return V
}

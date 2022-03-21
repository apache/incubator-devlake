package api

import (
	"github.com/merico-dev/lake/plugins/core"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var db *gorm.DB
var cfg *viper.Viper

func Init(config *viper.Viper, logger core.Logger, database *gorm.DB) {
	db = database
	cfg = config
}

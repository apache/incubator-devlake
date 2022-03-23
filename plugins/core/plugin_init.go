package core

import (
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Implement this interface if plugin needed some initialization
type PluginInit interface {
	Init(config *viper.Viper, logger Logger, db *gorm.DB) error
}

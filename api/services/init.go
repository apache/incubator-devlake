package services

import (
	"github.com/merico-dev/lake/api/models"
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	config.ReadConfig()
	var connectionString = viper.GetString("DB_URL")
	logger.Info("connectionString", connectionString)
	var err error
	db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	migrateDB()
}

func migrateDB() {
	err := db.AutoMigrate(&models.Source{})
	if err != nil {
		panic(err)
	}
	// TODO: create customer migration here
}

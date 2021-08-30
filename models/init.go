package models

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	var connectionString = config.V.GetString("DB_URL")
	logger.Info("JON >>> connectionString", connectionString)
	var err error
	Db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	migrateDB()
}

func migrateDB() {
	err := Db.AutoMigrate(&Task{})
	if err != nil {
		panic(err)
	}
	// TODO: create customer migration here
}

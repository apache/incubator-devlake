package main

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func Init() {
	config.ReadConfig()
	var connectionString = viper.GetString("DB_URL")
	logger.Info("connectionString", connectionString)
	var err error
	db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "jira_plugin_",
		},
	})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.Issue{}, &models.Board{})
}

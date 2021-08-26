package main

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/jira/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func (jira Jira) Init() {
	var connectionString = config.V.GetString("DB_URL")
	var err error
	db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(models.Issue{}, models.Board{})
}

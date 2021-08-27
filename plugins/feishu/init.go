package main

import (
	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/plugins/feishu/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func (plugin Feishu) Init() {
	var connectionString = config.V.GetString("DB_URL")
	var err error
	db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "feishu_plugin_",
		},
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&models.MeetingTopUserItem{})
	if err != nil {
		panic(err)
	}
}

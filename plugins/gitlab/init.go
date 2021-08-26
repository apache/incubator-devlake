package main

import (
	"github.com/merico-dev/lake/config"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

func (plugin Gitlab) Init() {
	config.ReadConfig()
	var connectionString = viper.GetString("DB_URL")
	var err error
	db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "gitlab_plugin_",
		},
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate()
	if err != nil {
		panic(err)
	}
}

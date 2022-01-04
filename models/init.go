package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/merico-dev/lake/config"
	lakeDb "github.com/merico-dev/lake/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func init() {
	config.LoadConfigFile()

	connectionString := lakeDb.GetConnectionString(map[string]string{}, false)

	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,         // Disable color
		},
	)

	Db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println("ERROR: >>> Mysql failed to connect")
		panic(err)
	}

	// TODO: create customer migration here
	err = lakeDb.MigrateDB("lake")
	if err != nil {
		fmt.Println("INFO: >>> This is shown when a database is already up to date. Please check migrations if you find an error")
	}
}

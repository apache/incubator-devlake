package e2e

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	mysqlGorm "gorm.io/driver/mysql"
)

func InitializeDb() (*sql.DB, error) {
	cfg := mysql.Config{
		User:   "merico",
		Passwd: "merico",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "lake",
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
		return nil, err
	}
	fmt.Println("Connected!")
	return db, nil
}

func InitializeGormDb() (*gorm.DB, error ){
	connectionString := "merico:merico@tcp(localhost:3306)/lake"
	db, err := gorm.Open(mysqlGorm.Open(connectionString))
	if err != nil {
		return nil, err
	}
	
	return db, nil
}
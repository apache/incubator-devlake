package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/merico-dev/lake/config"
)

var Db *sql.DB

func InitDb() {
	V := config.LoadConfigFile("../.env")
	connectionString := V.GetString("DB_URL")
	var err error
	Db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	pingErr := Db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

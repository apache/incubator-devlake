package db

import (
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
)

var ROOT_CONNECTION_STRING string = "mysql://root:admin@tcp(localhost:3306)/%v"
var MIGRATIONS_PATH string = "file://./db/migration"

func MigrateDB(dbName string) {
	err := RunMigrationsUp(dbName)
	if err != nil {
		fmt.Println("INFO: ", err)
	}
}

func RunMigrationsUp(dbName string) error {
	m, err := migrate.New(
		MIGRATIONS_PATH,
		fmt.Sprintf(ROOT_CONNECTION_STRING, dbName))

	if err != nil {
		fmt.Println("ERROR: Could not init migrate for UP: ", err)
		return err
	}
	err = m.Up()
	if err != nil {
		fmt.Println("ERROR: Could not run migrations UP: ", err)
		return err
	}
	return nil
}

func RunMigrationsDown(dbName string) error {
	m, err := migrate.New(
		MIGRATIONS_PATH,
		fmt.Sprintf(ROOT_CONNECTION_STRING, dbName))

	if err != nil {
		fmt.Println("ERROR: Could not init migrate for DOWN: ", err)
		return err
	}
	err = m.Down()
	if err != nil {
		fmt.Println("ERROR: Could not run migrations DOWN: ", err)
		return err
	}
	return nil
}

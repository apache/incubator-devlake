package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

var ROOT_CONNECTION_STRING string = "mysql://root:admin@tcp(localhost:3306)/lake"
var MIGRATIONS_PATH string = "file://../../db/migration"

func RunMigrationsUp() error {
	m, err := migrate.New(
		MIGRATIONS_PATH,
		ROOT_CONNECTION_STRING)

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

func RunMigrationsDown() error {
	m, err := migrate.New(
		MIGRATIONS_PATH,
		ROOT_CONNECTION_STRING)

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

package db

import (
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/merico-dev/lake/config"

	"github.com/golang-migrate/migrate/v4"
)

var MIGRATIONS_PATH string = "file://./db/migration"

func MigrateDB(dbName string) {
	err := RunMigrationsUp(dbName)
	if err != nil {
		fmt.Println("INFO: ", err)
	}
	RunPluginMigrationsUp(dbName, "github")
}

func RunPluginMigrationsUp(dbName string, pluginName string) error {
	connectionString := config.V.GetString("DB_URL")

	dbParams := fmt.Sprintf("x-migrations-table=schema_migrations_%v&x-migrations-table-quoted=1", pluginName)

	m, err := migrate.New(
		fmt.Sprintf("file://./plugins/%v/migration", pluginName),
		fmt.Sprintf("%v/%v?%v", connectionString, dbName, dbParams))

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

func RunMigrationsUp(dbName string) error {
	connectionString := config.V.GetString("DB_URL")
	m, err := migrate.New(
		MIGRATIONS_PATH,
		fmt.Sprintf("%v/%v", connectionString, dbName))

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
	connectionString := config.V.GetString("DB_URL")

	m, err := migrate.New(
		MIGRATIONS_PATH,
		fmt.Sprintf("%v/%v", connectionString, dbName))

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

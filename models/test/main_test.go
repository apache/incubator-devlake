package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var ROOT_CONNECTION_STRING string = "mysql://root:admin@tcp(localhost:3306)/lake"
var MIGRATIONS_PATH string = "file://../../db/migration"

func TestMain(m *testing.M) {
	runMigrationsDown()
	runMigrationsUp()
	setup()
	code := m.Run()
	os.Exit(code)
}

func runMigrationsUp() {
	m, err := migrate.New(
		MIGRATIONS_PATH,
		ROOT_CONNECTION_STRING)

	if err != nil {
		fmt.Println("ERROR: Could not run migrations UP: ", err)
	}
	m.Up()
}

func runMigrationsDown() {
	m, err := migrate.New(
		MIGRATIONS_PATH,
		ROOT_CONNECTION_STRING)

	if err != nil {
		fmt.Println("ERROR: Could not run migrations DOWN: ", err)
	}
	m.Down()
}

func setup() {
	fmt.Println("JON >>> setup", "setup")
	// TODO: DB setup...

	// Drop DB (if exists)
	// Create DB fresh
}

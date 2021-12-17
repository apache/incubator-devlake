package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// This file runs before ALL tests.
// This gives us the opportunity to run setup() and shutdown() functions...
// ...before and after m.Run()
// http://cs-guy.com/blog/2015/01/test-main/

var ROOT_CONNECTION_STRING string = "mysql://root:admin@tcp(localhost:3306)/lake"
var MIGRATIONS_PATH string = "file://../../db/migration"

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		os.Exit(1)
	}
	code := m.Run()
	os.Exit(code)
}

func setup() error {
	// Comment out because it caused the following error...
	// ERROR: Could not run migrations DOWN:  no change
	// Scripts are not behaving as expected. Needs more troubleshooting.

	err := runMigrationsDown()
	if err != nil {
		return err
	}
	err = runMigrationsUp()
	if err != nil {
		return err
	}
	return nil
}

func runMigrationsUp() error {
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

func runMigrationsDown() error {
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

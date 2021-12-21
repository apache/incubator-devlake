package test

import (
	"os"
	"testing"

	"github.com/merico-dev/lake/db"
)

// This file runs before ALL tests.
// This gives us the opportunity to run setup() and shutdown() functions...
// ...before and after m.Run()
// http://cs-guy.com/blog/2015/01/test-main/

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		os.Exit(1)
	}
	code := m.Run()
	os.Exit(code)
}

func setup() error {
	err := db.RunMigrationsDown("lake_test")
	if err != nil {
		return err
	}
	err = db.RunMigrationsUp("lake_test")
	if err != nil {
		return err
	}
	return nil
}

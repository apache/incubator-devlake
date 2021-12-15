package test

import (
	"os"
	"testing"

	"github.com/merico-dev/lake/models"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	// TODO: DB setup...

	// Drop DB (if exists)
	// Create DB fresh

	// Temporarily, connect to test db created manually...
	models.Init("merico:merico@tcp(localhost:3306)/lake_test?charset=utf8mb4&loc=Asia%2fShanghai&parseTime=True")
}

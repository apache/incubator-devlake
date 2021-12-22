package db

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/merico-dev/lake/plugins"

	"github.com/golang-migrate/migrate/v4"
)

var MIGRATIONS_PATH string = "file://./db/migration"

func MigrateDB(dbName string) {
	err := RunMigrationsUp(dbName)
	if err != nil {
		fmt.Println("INFO: ", err)
	}
	MigrateAllPluginDBSchemas(dbName, plugins.PluginDir())

}

func GolangMigrateDBString(pluginName string) string {
	dbParams := fmt.Sprintf("x-migrations-table=schema_migrations_%v&x-migrations-table-quoted=1", pluginName)
	connectionString := GetConnectionString(dbParams, true)
	fmt.Println("JON >>> connectionString", connectionString)
	return connectionString
}

func MigrateAllPluginDBSchemas(dbName string, pluginsDir string) error {
	walkErr := filepath.WalkDir(pluginsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fileName := d.Name()
		if strings.HasSuffix(fileName, ".so") {
			pluginName := fileName[0 : len(d.Name())-3]
			fmt.Println("JON >>> pluginName", pluginName)
			RunPluginMigrationsUp(dbName, pluginName)
		}
		return nil
	})
	return walkErr
}

func RunPluginMigrationsUp(dbName string, pluginName string) error {

	connectionString := GolangMigrateDBString(pluginName)
	path := fmt.Sprintf("file://./plugins/%v/migration", pluginName)

	fmt.Println("JON >>> path", path)

	m, err := migrate.New(path, connectionString)

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
	m, err := migrate.New(
		MIGRATIONS_PATH,
		GetConnectionString("", true))

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
		GetConnectionString("", true))

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

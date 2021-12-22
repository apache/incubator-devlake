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
	err := RunDomainLayerMigrationsUp(dbName)
	if err != nil {
		fmt.Println("INFO: ", err)
	}
	RunPluginMigrations(dbName, plugins.PluginDir())
}

// We need to maintain separate tables for migration tracking for each plugin.
func GolangMigrateDBString(pluginName string) string {
	dbParams := fmt.Sprintf("x-migrations-table=schema_migrations_%v&x-migrations-table-quoted=1", pluginName)
	connectionString := GetConnectionString(dbParams, true)
	return connectionString
}

// Run the migration folder for all plugins that have a compiled .so file
func RunPluginMigrations(dbName string, pluginsDir string) {
	walkErr := filepath.WalkDir(pluginsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fileName := d.Name()
		if strings.HasSuffix(fileName, ".so") {
			pluginName := fileName[0 : len(d.Name())-3]
			RunPluginMigrationsUp(dbName, pluginName)
		}
		return nil
	})
	if walkErr != nil {
		fmt.Println("ERROR >>> walkErr", walkErr)
	}
}

// Run the plugins/<pluginName>/migration folder in order
func RunPluginMigrationsUp(dbName string, pluginName string) {
	connectionString := GolangMigrateDBString(pluginName)
	path := fmt.Sprintf("file://./plugins/%v/migration", pluginName)

	m, err := migrate.New(path, connectionString)

	if err != nil {
		fmt.Println("INFO: RunPluginMigrationsUp: Could not init migrate for UP: ", pluginName, err)
	}
	err = m.Up()
	if err != nil {
		fmt.Println("INFO: RunPluginMigrationsUp: Could not run migrations UP: ", pluginName, err)
	}
}

func RunDomainLayerMigrationsUp(dbName string) error {
	m, err := migrate.New(
		MIGRATIONS_PATH,
		GetConnectionString("", true))

	if err != nil {
		fmt.Println("INFO: RunDomainLayerMigrationsUp: Could not init migrate for UP: ", err)
		return err
	}
	err = m.Up()
	if err != nil {
		fmt.Println("INFO: RunDomainLayerMigrationsUp: Could not run migrations UP: ", err)
		return err
	}
	return nil
}

func RunDomainLayerMigrationsDown(dbName string) error {
	m, err := migrate.New(
		MIGRATIONS_PATH,
		GetConnectionString("", true))

	if err != nil {
		fmt.Println("INFO: RunDomainLayerMigrationsDown: Could not init migrate for DOWN: ", err)
		return err
	}
	err = m.Down()
	if err != nil {
		fmt.Println("INFO: RunDomainLayerMigrationsDown: Could not run migrations DOWN: ", err)
		return err
	}
	return nil
}

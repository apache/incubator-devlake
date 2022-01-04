package db

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
)

// var MIGRATIONS_PATH string = "file://./db/migration"
var MIGRATIONS_PATH string = GetMigrationPath()

func GetRootPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	rootPath := filepath.Dir(exPath)
	return rootPath
}

func GetMigrationPath() string {
	return fmt.Sprintf("file://%v/db/migration", GetRootPath())
}

func GetBinPath() string {
	return fmt.Sprintf("%v/bin/plugins", GetRootPath())
}

func GetPluginsPath() string {
	return fmt.Sprintf("%v/plugins", GetRootPath())
}

func MigrateDB(dbName string) error {
	err := RunDomainLayerMigrationsUp(dbName)
	if err != nil {
		fmt.Println("INFO: ", err)
		return err
	}
	err = RunPluginMigrations(dbName, GetBinPath())
	if err != nil {
		return err
	}
	return nil
}

// We need to maintain separate tables for migration tracking for each plugin.
func GolangMigrateDBString(pluginName string) string {
	params := map[string]string{
		"x-migrations-table": fmt.Sprintf("schema_migrations_%v", pluginName),
	}
	connectionString := GetConnectionString(params, true)
	return connectionString
}

// Run the migration folder for all plugins that have a compiled .so file
func RunPluginMigrations(dbName string, binDir string) error {
	walkErr := filepath.WalkDir(binDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		fileName := d.Name()
		if strings.HasSuffix(fileName, ".so") {
			pluginName := fileName[0 : len(d.Name())-3]
			err = RunPluginMigrationsUp(dbName, pluginName)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if walkErr != nil {
		fmt.Println("ERROR >>> walkErr", walkErr)
		return walkErr
	}
	return nil
}

// Run the plugins/<pluginName>/migration folder in order
func RunPluginMigrationsUp(dbName string, pluginName string) error {
	connectionString := GolangMigrateDBString(pluginName)
	path := fmt.Sprintf("file://./plugins/%v/migration", pluginName)

	m, err := migrate.New(path, connectionString)

	if err != nil {
		fmt.Println("INFO: RunPluginMigrationsUp: No migration scripts found for plugin: ", pluginName, err)
		return nil
	}
	err = m.Up()
	if err != nil {
		fmt.Println("INFO: RunPluginMigrationsUp: Database already up to date: ", pluginName, err)
		return nil
	}
	return nil
}

func RunDomainLayerMigrationsUp(dbName string) error {
	m, err := migrate.New(
		MIGRATIONS_PATH,
		GetConnectionString(map[string]string{}, true))

	if err != nil {
		fmt.Println("INFO: RunDomainLayerMigrationsUp: No migration scripts found: ", err)
		return nil
	}
	err = m.Up()
	if err != nil {
		fmt.Println("INFO: RunDomainLayerMigrationsUp: Database already up to date: ", err)
		return nil
	}
	return nil
}

func RunDomainLayerMigrationsDown(dbName string) error {
	m, err := migrate.New(
		MIGRATIONS_PATH,
		GetConnectionString(map[string]string{}, true))

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

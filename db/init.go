package db

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/merico-dev/lake/plugins"

	"github.com/golang-migrate/migrate/v4"
)

// var MIGRATIONS_PATH string = "file://./db/migration"
var MIGRATIONS_PATH string = GetFilePath()
var PLUGINS_PATH string = GetPluginsPath()

func GetPluginsPath() string {
	relPath := plugins.PluginDir()
	exPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%v/../../../%v", exPath, relPath)
}
func GetFilePath() string {
	// ex, err := os.Executable()
	// if err != nil {
	// 	panic(err)
	// }
	exPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// exPath := filepath.Dir(ex)
	// The reason for the "../"s is the working directory starts at lake/test/api/task (due to init() functions I assume)
	return fmt.Sprintf("file://%v/../../../../lake/db/migration", exPath)
}

func MigrateDB(dbName string) {
	err := RunDomainLayerMigrationsUp(dbName)
	if err != nil {
		fmt.Println("INFO: ", err)
	}

	RunPluginMigrations(dbName)
}

// We need to maintain separate tables for migration tracking for each plugin.
func GolangMigrateDBString(pluginName string) string {
	dbParams := fmt.Sprintf("x-migrations-table=schema_migrations_%v&x-migrations-table-quoted=1", pluginName)
	connectionString := GetConnectionString(dbParams, true)
	return connectionString
}

// Run the migration folder for all plugins that have a compiled .so file
func RunPluginMigrations(dbName string) {
	fmt.Println("KEVIN >>> PLUGINS_PATH", PLUGINS_PATH)
	walkErr := filepath.WalkDir(PLUGINS_PATH, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("KEVIN >>> err", err)
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

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Run the plugins/<pluginName>/migration folder in order
func RunPluginMigrationsUp(dbName string, pluginName string) {
	connectionString := GolangMigrateDBString(pluginName)
	path := fmt.Sprintf("file://./../../../plugins/%v/migration", pluginName)
	// check if path exists
	pathExists, err := exists(path)
	if err != nil {
		panic(err)
	}
	if pathExists {
		m, err := migrate.New(path, connectionString)

		if err != nil {
			fmt.Println("INFO: RunPluginMigrationsUp: Could not init migrate for UP: ", pluginName, err)
		}
		err = m.Up()
		if err != nil {
			fmt.Println("INFO: RunPluginMigrationsUp: Could not run migrations UP: ", pluginName, err)
		}
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

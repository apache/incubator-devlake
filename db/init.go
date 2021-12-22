package db

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

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
	MigrateAllPluginDBSchemas("../plugins")
	RunPluginMigrationsUp(dbName, "github")
}

func formatDBString(dbName string, pluginName string) string {
	connectionString := config.V.GetString("DB_URL")

	dbParams := fmt.Sprintf("x-migrations-table=schema_migrations_%v&x-migrations-table-quoted=1", pluginName)
	dbConnectionWithParams := fmt.Sprintf("mysql://%v/%v?%v", connectionString, dbName, dbParams)

	fmt.Println("JON >>> dbConnectionWithParams", dbConnectionWithParams)
	return dbConnectionWithParams
}

func MigrateAllPluginDBSchemas(pluginsDir string) error {
	walkErr := filepath.WalkDir(pluginsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("JON >>> err", err)
			return err
		}
		fileName := d.Name()
		fmt.Println("JON >>> d.Name()", d.Name())
		if strings.HasSuffix(fileName, ".so") {
			pluginName := fileName[0 : len(d.Name())-3]
			fmt.Println("JON >>> pluginName", pluginName)
			// plug, loadErr := plugin.Open(path)
			// if loadErr != nil {
			// 	return loadErr
			// }
			// symPluginEntry, pluginEntryError := plug.Lookup("PluginEntry")
			// if pluginEntryError != nil {
			// 	return pluginEntryError
			// }
			// plugEntry, ok := symPluginEntry.(Plugin)
			// if !ok {
			// 	return fmt.Errorf("%v PluginEntry must implement Plugin interface", pluginName)
			// }
			// plugEntry.Init()
			// logger.Info(`[plugins] init a plugin success`, pluginName)
			// err = RegisterPlugin(pluginName, plugEntry)
			// if err != nil {
			// 	return nil
			// }
			// logger.Info("[plugins] plugin loaded", pluginName)
		}
		return nil
	})
	return walkErr
}

func RunPluginMigrationsUp(dbName string, pluginName string) error {

	dbConnectionWithParams := formatDBString(dbName, pluginName)
	path := fmt.Sprintf("file://./plugins/%v/migration", pluginName)

	fmt.Println("JON >>> path", path)

	m, err := migrate.New(path, dbConnectionWithParams)

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
		fmt.Sprintf("mysql://%v/%v", connectionString, dbName))

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

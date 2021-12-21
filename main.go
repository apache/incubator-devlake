package main

import (
	"fmt"

	"github.com/merico-dev/lake/api"
	"github.com/merico-dev/lake/db"
	"github.com/merico-dev/lake/plugins"
)

func main() {
	migrateDB()
	startAPI()
}

func startAPI() {
	err := plugins.LoadPlugins(plugins.PluginDir())
	if err != nil {
		panic(err)
	}
	api.CreateApiService()
	println("Hello, lake")
}

func migrateDB() {
	err := db.RunMigrationsUp("lake")
	if err != nil {
		fmt.Println("INFO: ", err)
	}
}

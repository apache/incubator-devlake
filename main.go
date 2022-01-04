package main

import (
	"github.com/merico-dev/lake/api"
	"github.com/merico-dev/lake/db"

	"github.com/merico-dev/lake/plugins"
)

func main() {
	startAPI()
}

func startAPI() {
	err := plugins.LoadPlugins(db.GetPluginsPath())
	if err != nil {
		panic(err)
	}
	api.CreateApiService()
	println("Hello, lake")
}

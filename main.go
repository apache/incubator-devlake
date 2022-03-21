package main

import (
	"github.com/merico-dev/lake/api"
	"github.com/merico-dev/lake/config"
)

func main() {
	err := worker.LoadPlugins(config.GetConfig().GetString("PLUGIN_DIR"))
	if err != nil {
		panic(err)
	}
	api.CreateApiService()
}

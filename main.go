package main

import (
	"fmt"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins"
)

func main() {
	config.ReadConfig()
	loadPlugins()
	runJiraPlugin()
	CreateApiService()
}

func loadPlugins() {
	err := plugins.LoadPlugins("plugins")
	if err != nil {
		logger.Error("Failed to LoadPlugins %v", err)
	}
	if len(plugins.Plugins) == 0 {
		logger.Error("No plugin found", false)
		return
	}
}

func runJiraPlugin() {
	name := "jira"
	options := map[string]interface{}{
		"boardId": 20,
	}
	progress := make(chan float32)
	fmt.Printf("start runing plugin %v\n", name)
	go func() {
		_ = plugins.RunPlugin(name, options, progress)
	}()
	for p := range progress {
		fmt.Printf("running plugin %v, progress: %v\n", name, p*100)
	}
	fmt.Printf("end running plugin %v\n", name)
}

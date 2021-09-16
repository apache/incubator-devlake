package main // must be main for plugin entry point

import (
	"github.com/merico-dev/lake/logger" // A pseudo type for Plugin Interface implementation
	"github.com/merico-dev/lake/plugins/github/tasks"
)

type Github string

func (plugin Github) Description() string {
	return "To collect and enrich data from GitHub"
}

func (plugin Github) Execute(options map[string]interface{}, progress chan<- float32) {
	logger.Print("start github plugin execution")

	owner, ok := options["owner"]
	if !ok {
		logger.Print("owner is required for github execution")
		return
	}
	ownerString := owner.(string)

	repositoryName, ok := options["repositoryName"]
	if !ok {
		logger.Print("repositoryName is required for github execution")
		return
	}
	repositoryNameString := repositoryName.(string)

	if err := tasks.CollectRepository(ownerString, repositoryNameString); err != nil {
		logger.Error("Could not collect repositories: ", err)
		return
	}

	progress <- 1

	close(progress)

}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Github //nolint

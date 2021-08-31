package main // must be main for plugin entry point

import (
	"time"

	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/pokemon/tasks"
)

// A pseudo type for Plugin Interface implementation
type Pokemon string

func (plugin Pokemon) Description() string {
	return "To collect and enrich data from https://pokeapi.co/"
}

func (plugin Pokemon) Execute(options map[string]interface{}, progress chan<- float32) {

	err := tasks.CollectAllPokemon()

	if err != nil {
		logger.Error("Something went wrong: ", err.Error())
	}

	time.Sleep(1 * time.Second)
	progress <- 0.1
	time.Sleep(1 * time.Second)
	progress <- 0.5
	time.Sleep(1 * time.Second)
	progress <- 1
	logger.Print("end pokemon plugin execution")
	close(progress)
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Pokemon //nolint

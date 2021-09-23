package main // must be main for plugin entry point

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/domainlayer/models/code"
)

// plugin interface
type DomainLayer string

func (plugin DomainLayer) Init() {
	err := lakeModels.Db.AutoMigrate(
		&code.Repo{},
		&code.Commit{},
		&code.Pr{},
		&code.Note{},
	)
	if err != nil {
		panic(err)
	}
}

func (plugin DomainLayer) Description() string {
	return "Domain Layer"
}

func (plugin DomainLayer) Execute(options map[string]interface{}, progress chan<- float32) {
	progress <- 1
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry DomainLayer //nolint

package main // must be main for plugin entry point

import (
	"context"

	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/merico-analysis-engine-domain/tasks"
	"github.com/mitchellh/mapstructure"
)

type AEDomainOptions struct {
	Tasks []string `json:"tasks,omitempty"`
}

// plugin interface
type AEDomain string

func (plugin AEDomain) Init() {
}

func (plugin AEDomain) Description() string {
	return "Convert AE Entities to Domain Layer Entities"
}

func (plugin AEDomain) Execute(
	options map[string]interface{},
	progress chan<- float32,
	ctx context.Context,
) error {
	// process options
	var op AEDomainOptions
	var err error
	err = mapstructure.Decode(options, &op)
	if err != nil {
		return err
	}

	// run tasks
	logger.Print("start AEDomain plugin execution")
	err = tasks.ConvertCommits()
	if err != nil {
		return err
	}
	progress <- 1
	logger.Print("end AEDomain plugin execution")
	close(progress)
	return nil
}

func (plugin AEDomain) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/aedomain"
}

func (plugin AEDomain) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return make(map[string]map[string]core.ApiResourceHandler)
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry AEDomain //nolint

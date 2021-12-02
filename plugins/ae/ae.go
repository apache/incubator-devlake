package main // must be main for plugin entry point

import (
	"context"
	"fmt"

	"github.com/merico-dev/lake/logger" // A pseudo type for Plugin Interface implementation
	"github.com/merico-dev/lake/plugins/ae/api"
	"github.com/merico-dev/lake/plugins/ae/tasks"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/mitchellh/mapstructure"
)

type AEOptions struct {
	Tasks []string `json:"tasks,omitempty"`
}
type AE string

func (plugin AE) Description() string {
	return "To collect and enrich data from AE"
}

func (plugin AE) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	logger.Print("start ae plugin execution")

	projectId, ok := options["projectId"]
	if !ok {
		return fmt.Errorf("projectId is required for ae execution")
	}

	projectIdInt := int(projectId.(float64))
	if projectIdInt < 0 {
		return fmt.Errorf("projectId is invalid")
	}

	var op AEOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return err
	}

	progress <- 0.1
	if err := tasks.CollectProject(projectIdInt); err != nil {
		return fmt.Errorf("could not collect project: %v", err)
	}

	progress <- 0.25

	if err := tasks.CollectCommits(projectIdInt); err != nil {
		return fmt.Errorf("could not collect commits: %v", err)
	}

	progress <- 0.75

	if err := tasks.SetDevEqOnCommits(); err != nil {
		return fmt.Errorf("could not enhance commits with AE dev equivalent: %v", err)
	}

	progress <- 1
	return nil
}

func (plugin AE) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/ae"
}

func (plugin AE) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"GET": api.TestConnection,
		},
		"sources": {
			"GET":  api.ListSources,
			"POST": api.PutSource,
		},
		"sources/:sourceId": {
			"GET": api.GetSource,
			"PUT": api.PutSource,
		},
	}
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry AE //nolint

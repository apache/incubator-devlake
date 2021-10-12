package main // must be main for plugin entry point

import (
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jenkinsdomain/tasks"
	"github.com/mitchellh/mapstructure"
)

type JenkinsDomainOptions struct {
	Tasks []string `json:"tasks,omitempty"`
}

// plugin interface
type JenkinsDomain string

func (plugin JenkinsDomain) Init() {
}

func (plugin JenkinsDomain) Description() string {
	return "Convert Jenkins Entities to Domain Layer Entities"
}

func (plugin JenkinsDomain) Execute(options map[string]interface{}, progress chan<- float32) {
	// process options
	var op JenkinsDomainOptions
	var err error
	err = mapstructure.Decode(options, &op)
	if err != nil {
		logger.Error("Error: ", err)
		return
	}
	tasksToRun := make(map[string]bool, len(op.Tasks))
	for _, task := range op.Tasks {
		tasksToRun[task] = true
	}
	if len(tasksToRun) == 0 {
		tasksToRun = map[string]bool{
			"convertJobs":   true,
			"convertBuilds": true,
		}
	}

	// run tasks
	logger.Print("start JenkinsDomain plugin execution")
	if tasksToRun["convertJobs"] {
		err := tasks.ConvertJobs()
		if err != nil {
			logger.Error("Error: ", err)
			return
		}
	}
	progress <- 0.2
	if tasksToRun["convertBuilds"] {
		err = tasks.ConvertBuilds()
		if err != nil {
			logger.Error("Error: ", err)
			return
		}
	}
	progress <- 1
	logger.Print("end JenkinsDomain plugin execution")
	close(progress)
}

func (plugin JenkinsDomain) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/jenkinsdomain"
}

func (plugin JenkinsDomain) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return make(map[string]map[string]core.ApiResourceHandler)
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry JenkinsDomain //nolint

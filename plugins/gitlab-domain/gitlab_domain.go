package main // must be main for plugin entry point

import (
	"context"

	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab-domain/tasks"
	"github.com/mitchellh/mapstructure"
)

type GitlabDomainOptions struct {
	Tasks []string `json:"tasks,omitempty"`
}

// plugin interface
type GitlabDomain string

func (plugin GitlabDomain) Init() {}

func (plugin GitlabDomain) Description() string {
	return "Convert Gitlab Entities to Domain Layer Entities"
}

func (plugin GitlabDomain) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	// process options
	var op GitlabDomainOptions
	var err error
	err = mapstructure.Decode(options, &op)
	if err != nil {
		return err
	}

	tasksToRun := make(map[string]bool, len(op.Tasks))
	for _, task := range op.Tasks {
		tasksToRun[task] = true
	}
	if len(tasksToRun) == 0 {
		tasksToRun = map[string]bool{
			"convertRepos":   true,
			"convertMrs":     true,
			"convertCommits": true,
			"convertNotes":   true,
		}
	}

	// run tasks
	logger.Print("start GitlabDomain plugin execution")
	if tasksToRun["convertRepos"] {
		progress <- 0.01
		err = tasks.ConvertRepos()
		if err != nil {
			return err
		}
	}
	if tasksToRun["convertMrs"] {
		progress <- 0.05
		err = tasks.ConvertPrs()
		if err != nil {
			return err
		}
	}
	if tasksToRun["convertCommits"] {
		progress <- 0.07
		err = tasks.ConvertCommits()
		if err != nil {
			return err
		}
	}
	if tasksToRun["convertNotes"] {
		progress <- 0.09
		err = tasks.ConvertNotes()
		if err != nil {
			return err
		}
	}
	progress <- 1
	logger.Print("end GitlabDomain plugin execution")
	close(progress)
	return nil
}

func (plugin GitlabDomain) RootPkgPath() string {
	return "gitlab.com/merico-dev/lake/plugins/gitlab-domain"
}

func (plugin GitlabDomain) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return make(map[string]map[string]core.ApiResourceHandler)
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry GitlabDomain //nolint

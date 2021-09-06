package main // must be main for plugin entry point

import (
	"github.com/merico-dev/lake/logger" // A pseudo type for Plugin Interface implementation
	"github.com/merico-dev/lake/plugins/gitlab/tasks"
)

type Gitlab string

func (plugin Gitlab) Description() string {
	return "To collect and enrich data from Gitlab"
}

func (plugin Gitlab) Execute(options map[string]interface{}, progress chan<- float32) {
	logger.Print("start gitlab plugin execution")

	projectId, ok := options["projectId"]
	if !ok {
		logger.Print("projectId is required for gitlab execution")
		return
	}

	projectIdInt := int(projectId.(float64))
	if projectIdInt < 0 {
		logger.Print("boardId is invalid")
		return
	}
	c := make(chan bool)
	go func() {
		err := tasks.CollectProjects(projectIdInt, c)
		if err != nil {
			logger.Error("Could not collect projects: ", err)
			return
		}
	}()
	<-c
	err := tasks.CollectCommits(projectIdInt)
	if err != nil {
		logger.Error("Could not collect commits: ", err)
		return
	}

	mergeRequestErr := tasks.CollectMergeRequests(projectIdInt)
	if mergeRequestErr != nil {
		logger.Error("Could not collect merge requests: ", mergeRequestErr)
		return
	}
	progress <- 1

	close(progress)
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Gitlab //nolint

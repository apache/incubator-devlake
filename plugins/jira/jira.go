package main // must be main for plugin entry point

import (
	"time"

	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/jira/tasks"
)

// A pseudo type for Plugin Interface implementation
type Jira string

func (plugin Jira) Description() string {
	return "To collect and enrich data from JIRA"
}

func (plugin Jira) Execute(options map[string]interface{}, progress chan<- float32) {
	boardId, ok := options["boardId"]
	if !ok {
		logger.Print("boardId is required for jira execution")
		return
	}
	boardIdInt := uint64(boardId.(float64))
	if boardIdInt == 0 {
		logger.Print("boardId is invalid")
		return
	}
	logger.Print("start jira plugin execution")
	err := tasks.CollectBoard(boardIdInt)
	if err != nil {
		logger.Error("Error: ", err)
		return
	}
	progress <- 0.01
	err = tasks.CollectIssues(boardIdInt)
	if err != nil {
		logger.Error("Error: ", err)
		return
	}
	progress <- 0.8
	time.Sleep(1 * time.Second)
	progress <- 1
	logger.Print("end jira plugin execution")
	close(progress)
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Jira //nolint

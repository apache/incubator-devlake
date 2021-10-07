package main // must be main for plugin entry point

import (
	"github.com/merico-dev/lake/logger"
	"github.com/merico-dev/lake/plugins/jiradomain/tasks"
	"github.com/mitchellh/mapstructure"
)

type JiraDomainOptions struct {
	BoardId uint64   `json:"boardId"`
	Tasks   []string `json:"tasks,omitempty"`
}

// plugin interface
type JiraDomain string

func (plugin JiraDomain) Init() {
}

func (plugin JiraDomain) Description() string {
	return "Convert Jira Entities to Domain Layer Entities"
}

func (plugin JiraDomain) Execute(options map[string]interface{}, taskId uint64, progress chan<- float32) {
	// process options
	var op JiraDomainOptions
	var err error
	err = mapstructure.Decode(options, &op)
	if err != nil {
		logger.Error("Error: ", err)
		return
	}
	if op.BoardId == 0 {
		logger.Print("boardId is invalid")
		return
	}
	boardId := op.BoardId
	tasksToRun := make(map[string]bool, len(op.Tasks))
	for _, task := range op.Tasks {
		tasksToRun[task] = true
	}
	if len(tasksToRun) == 0 {
		tasksToRun = map[string]bool{
			"convertBoard":      true,
			"convertIssues":     true,
			"convertChangelogs": true,
		}
	}

	// run tasks
	logger.Print("start JiraDomain plugin execution")
	if tasksToRun["convertBoard"] {
		err := tasks.ConvertBoard(boardId)
		if err != nil {
			logger.Error("Error: ", err)
			return
		}
	}
	progress <- 0.01
	if tasksToRun["convertIssues"] {
		err = tasks.ConvertIssues(boardId)
		if err != nil {
			logger.Error("Error: ", err)
			return
		}
	}
	progress <- 0.7
	if tasksToRun["convertChangelogs"] {
		err = tasks.ConvertChangelogs(boardId)
		if err != nil {
			logger.Error("Error: ", err)
			return
		}
	}
	progress <- 1
	logger.Print("end JiraDomain plugin execution")
	close(progress)
}

func (plugin JiraDomain) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/jiradomain"
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry JiraDomain //nolint

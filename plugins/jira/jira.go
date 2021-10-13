package main // must be main for plugin entry point

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/plugins/jira/tasks"
	"github.com/merico-dev/lake/utils"
	"github.com/mitchellh/mapstructure"
)

type JiraOptions struct {
	BoardId uint64   `json:"boardId"`
	Tasks   []string `json:"tasks,omitempty"`
	Since   string
}

// plugin interface
type Jira string

func (plugin Jira) Init() {
	err := lakeModels.Db.AutoMigrate(
		&models.JiraIssue{},
		&models.JiraBoard{},
		&models.JiraBoardIssue{},
		&models.JiraChangelog{},
		&models.JiraChangelogItem{},
	)
	if err != nil {
		panic(err)
	}
}

func (plugin Jira) Description() string {
	return "To collect and enrich data from JIRA"
}

func (plugin Jira) Execute(options map[string]interface{}, progress chan<- float32) {
	// process options
	var op JiraOptions
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
	convertedSince := utils.ConvertStringToTime(op.Since)
	if convertedSince.IsZero() {
		fmt.Println("ERROR >>> Since value is in the wrong format")
		return
	}
	boardId := op.BoardId
	tasksToRun := make(map[string]bool, len(op.Tasks))
	for _, task := range op.Tasks {
		tasksToRun[task] = true
	}
	if len(tasksToRun) == 0 {
		tasksToRun = map[string]bool{
			"collectBoard":      true,
			"collectIssues":     true,
			"collectChangelogs": true,
			"enrichIssues":      true,
		}
	}

	// run tasks
	logger.Print("start jira plugin execution")
	if tasksToRun["collectBoard"] {
		err := tasks.CollectBoard(boardId)
		if err != nil {
			logger.Error("Error: ", err)
			return
		}
	}
	progress <- 0.01
	if tasksToRun["collectIssues"] {
		err = tasks.CollectIssues(boardId, convertedSince)
		if err != nil {
			logger.Error("Error: ", err)
			return
		}
	}
	progress <- 0.5
	if tasksToRun["collectChangelogs"] {
		err = tasks.CollectChangelogs(boardId, convertedSince)
		if err != nil {
			logger.Error("Error: ", err)
			return
		}
	}
	progress <- 0.8
	if tasksToRun["enrichIssues"] {
		err = tasks.EnrichIssues(boardId)
		if err != nil {
			logger.Error("Error: ", err)
			return
		}
	}
	progress <- 1
	logger.Print("end jira plugin execution")
	close(progress)
}

func (plugin Jira) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/jira"
}

func (plugin Jira) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"echo": {
			"POST": func(input *core.ApiResourceInput) (*core.ApiResourceOutput, error) {
				return &core.ApiResourceOutput{Body: input.Body}, nil
			},
		},
	}
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Jira //nolint

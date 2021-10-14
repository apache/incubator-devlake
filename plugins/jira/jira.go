package main // must be main for plugin entry point

import (
	"context"
	"fmt"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/api"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/plugins/jira/tasks"
	"github.com/merico-dev/lake/utils"
	"github.com/mitchellh/mapstructure"
)

type JiraOptions struct {
	SourceId uint64   `json:"sourceId"`
	BoardId  uint64   `json:"boardId"`
	Tasks    []string `json:"tasks,omitempty"`
	Since    string
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
		&models.JiraSource{},
		&models.JiraIssueTypeMapping{},
		&models.JiraIssueStatusMapping{},
	)
	if err != nil {
		panic(err)
	}
}

func (plugin Jira) Description() string {
	return "To collect and enrich data from JIRA"
}

func (plugin Jira) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) {
	// process options
	var op JiraOptions
	var err error
	var boardIds []uint64
	err = mapstructure.Decode(options, &op)
	if err != nil {
		logger.Error("Error: ", err)
		return
	}
	if op.SourceId == 0 {
		// sourceId is required
		logger.Print("sourceId is invalid")
		return
	}
	source := &models.JiraSource{}
	err = lakeModels.Db.Find(source, op.SourceId).Error
	if err != nil {
		logger.Print("jira source not found")
		return
	}
	if op.BoardId == 0 {
		// boardId omitted: to collect all boards of the data source
		err = lakeModels.Db.Model(&models.JiraBoard{}).Where("source_id = ?", op.SourceId).Pluck("id", &boardIds).Error
		if err != nil {
			logger.Error("Error: ", err)
			return
		}
	} else {
		boardIds = []uint64{op.BoardId}
	}
	if len(boardIds) == 0 {
		logger.Error("no board to collect", op)
		return
	}
	convertedSince := utils.ConvertStringToTime(op.Since)
	if convertedSince.IsZero() {
		fmt.Println("ERROR >>> Since value is in the wrong format")
		return
	}
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

	setBoardProgress := func(boardIndex int, boardProgress float32) {
		boardPart := 1.0 / float32(len(boardIds))
		progress <- boardPart*float32(boardIndex) + boardPart*boardProgress
	}

	// run tasks
	logger.Print("start jira plugin execution")

	jiraApiClient, err := tasks.NewJiraApiClientBySourceId(op.SourceId)
	if err != nil {
		logger.Error("failed to create jira api client", err)
		return
	}
	for i, boardId := range boardIds {
		if tasksToRun["collectBoard"] {
			err := tasks.CollectBoard(jiraApiClient, source, boardId)
			if err != nil {
				logger.Error("Error: ", err)
				return
			}
		}
		setBoardProgress(i, 0.01)
		if tasksToRun["collectIssues"] {
			err = tasks.CollectIssues(jiraApiClient, source, boardId, convertedSince, ctx)
			if err != nil {
				logger.Error("Error: ", err)
				return
			}
		}
		setBoardProgress(i, 0.5)
		if tasksToRun["collectChangelogs"] {
			err = tasks.CollectChangelogs(jiraApiClient, boardId, convertedSince, ctx)
			if err != nil {
				logger.Error("Error: ", err)
				return
			}
		}
		setBoardProgress(i, 0.8)
		if tasksToRun["enrichIssues"] {
			err = tasks.EnrichIssues(source, boardId)
			if err != nil {
				logger.Error("Error: ", err)
				return
			}
		}
		setBoardProgress(i, 1.0)
	}
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
		"sources": {
			"POST": api.PostSources,
			"GET":  api.ListSources,
		},
		"sources/:sourceId": {
			"PUT":    api.PutSource,
			"DELETE": api.DeleteSource,
			"GET":    api.GetSource,
		},
		"sources/:sourceId/type-mappings": {
			"POST": api.PostIssueTypeMappings,
			"GET":  api.ListIssueTypeMappings,
		},
		"sources/:sourceId/type-mappings/:userType": {
			"PUT":    api.PutIssueTypeMapping,
			"DELETE": api.DeleteIssueTypeMapping,
		},
		"sources/:sourceId/type-mappings/:userType/status-mappings": {
			"POST": api.PostIssueStatusMappings,
			"GET":  api.ListIssueStatusMappings,
		},
		"sources/:sourceId/type-mappings/:userType/status-mappings/:userStatus": {
			"PUT":    api.PutIssueStatusMapping,
			"DELETE": api.DeleteIssueStatusMapping,
		},
	}
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Jira //nolint

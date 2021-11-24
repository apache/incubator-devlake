package main // must be main for plugin entry point

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/jira/api"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/plugins/jira/tasks"
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
		&models.JiraProject{},
		&models.JiraUser{},
		&models.JiraIssue{},
		&models.JiraBoard{},
		&models.JiraBoardIssue{},
		&models.JiraChangelog{},
		&models.JiraChangelogItem{},
		&models.JiraSource{},
		&models.JiraIssueTypeMapping{},
		&models.JiraIssueStatusMapping{},
		&models.JiraSprint{},
		&models.JiraBoardSprint{},
		&models.JiraSprintIssue{},
		&models.JiraWorklog{},
	)
	if err != nil {
		panic(err)
	}
}

func (plugin Jira) Description() string {
	return "To collect and enrich data from JIRA"
}

func (plugin Jira) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	// process options
	var op JiraOptions
	var err error
	var boardIds []uint64
	err = mapstructure.Decode(options, &op)
	if err != nil {
		return err
	}
	if op.SourceId == 0 {
		return fmt.Errorf("sourceId is invalid")
	}
	source := &models.JiraSource{}
	err = lakeModels.Db.Find(source, op.SourceId).Error
	if err != nil {
		return err
	}
	if op.BoardId == 0 {
		// boardId omitted: to collect all boards of the data source
		err = lakeModels.Db.Model(&models.JiraBoard{}).Where("source_id = ?", op.SourceId).Pluck("id", &boardIds).Error
		if err != nil {
			return err
		}
	} else {
		boardIds = []uint64{op.BoardId}
	}
	if len(boardIds) == 0 {
		return fmt.Errorf("no board to collect")
	}

	var since time.Time
	if op.Since != "" {
		since, err = time.Parse("2006-01-02T15:04:05Z", op.Since)
		if err != nil {
			return err
		}
	}
	tasksToRun := make(map[string]bool, len(op.Tasks))
	for _, task := range op.Tasks {
		tasksToRun[task] = true
	}
	if len(tasksToRun) == 0 {
		tasksToRun = map[string]bool{
			"collectBoard":      true,
			"collectProjects":   true,
			"collectIssues":     true,
			"collectChangelogs": true,
			"enrichIssues":      true,
			"collectSprints":    true,
			"collectUsers":      true,
			"convertBoard":      true,
			"convertIssues":     true,
			"convertWorklogs":   true,
			"convertChangelogs": true,
			"convertUsers":      true,
			"convertSprints":    true,
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
		return fmt.Errorf("failed to create jira api client: %v", err)
	}
	for i, boardId := range boardIds {
		if tasksToRun["collectProjects"] {
			err := tasks.CollectProjects(jiraApiClient, op.SourceId)
			if err != nil {
				return err
			}
		}
		if tasksToRun["collectUsers"] {
			err := tasks.CollectUsers(jiraApiClient, op.SourceId)
			if err != nil {
				return err
			}
		}
		if tasksToRun["collectBoard"] {
			err := tasks.CollectBoard(jiraApiClient, source, boardId)
			if err != nil {
				return err
			}
		}
		setBoardProgress(i, 0.01)
		if tasksToRun["collectIssues"] {
			err = tasks.CollectIssues(jiraApiClient, source, boardId, since, ctx)
			if err != nil {
				return err
			}
		}
		setBoardProgress(i, 0.1)
		if tasksToRun["collectChangelogs"] {
			err = tasks.CollectChangelogs(jiraApiClient, source, boardId, ctx)
			if err != nil {
				return err
			}
		}
		setBoardProgress(i, 0.2)
		if tasksToRun["enrichIssues"] {
			err = tasks.EnrichIssues(source, boardId)
			if err != nil {
				return err
			}
		}
		setBoardProgress(i, 0.3)
		if tasksToRun["collectSprints"] {
			err = tasks.CollectSprint(jiraApiClient, source, boardId)
			if err != nil {
				return err
			}
		}
		setBoardProgress(i, 0.4)
		if tasksToRun["convertBoard"] {
			err := tasks.ConvertBoard(op.SourceId, boardId)
			if err != nil {
				return err
			}
		}
		setBoardProgress(i, 0.5)
		if tasksToRun["convertUsers"] {
			err := tasks.ConvertUsers(op.SourceId)
			if err != nil {
				return err
			}
		}
		setBoardProgress(i, 0.6)
		if tasksToRun["convertIssues"] {
			err = tasks.ConvertIssues(op.SourceId, boardId)
			if err != nil {
				return err
			}
		}
		setBoardProgress(i, 0.7)
		if tasksToRun["convertWorklogs"] {
			err = tasks.ConvertWorklog(op.SourceId, boardId)
			if err != nil {
				return err
			}
		}
		setBoardProgress(i, 0.8)
		if tasksToRun["convertChangelogs"] {
			err = tasks.ConvertChangelogs(op.SourceId, boardId)
			if err != nil {
				return err
			}
		}
		setBoardProgress(i, 0.9)
		if tasksToRun["convertSprints"] {
			err = tasks.ConvertSprint(op.SourceId, boardId)
			if err != nil {
				logger.Error("convertSprints", err)
				return err
			}
		}
		setBoardProgress(i, 1.0)

	}
	logger.Print("end jira plugin execution")
	return nil
}

func (plugin Jira) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/jira"
}

func (plugin Jira) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return map[string]map[string]core.ApiResourceHandler{
		"test": {
			"GET": api.TestConnection,
		},
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
		"sources/:sourceId/epics": {
			"GET": api.GetEpicsBySourceId,
		},
		"sources/:sourceId/granularities": {
			"GET": api.GetGranularitiesBySourceId,
		},
		"sources/:sourceId/boards": {
			"GET": api.GetBoardsBySourceId,
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
		"sources/:sourceId/proxy/rest/*path": {
			"GET": api.Proxy,
		},
	}
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Jira //nolint

// standalone mode for debugging
func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		panic(fmt.Errorf("Usage: jira <source_id> <board_id>"))
	}
	sourceId, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		panic(fmt.Errorf("error paring source_id: %w", err))
	}
	boardId, err := strconv.ParseUint(args[1], 10, 64)
	if err != nil {
		panic(fmt.Errorf("error paring board_id: %w", err))
	}

	PluginEntry.Init()
	progress := make(chan float32)
	go func() {
		err := PluginEntry.Execute(
			map[string]interface{}{
				"sourceId": sourceId,
				"boardId":  boardId,
				//"tasks":    []string{"enrichIssues"},
			},
			progress,
			context.Background(),
		)
		if err != nil {
			panic(err)
		}
	}()
	for p := range progress {
		fmt.Println(p)
	}
}

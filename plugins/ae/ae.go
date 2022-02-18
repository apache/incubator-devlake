package main // must be main for plugin entry point

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/logger" // A pseudo type for Plugin Interface implementation
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/ae/api"
	"github.com/merico-dev/lake/plugins/ae/models"
	"github.com/merico-dev/lake/plugins/ae/tasks"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/mitchellh/mapstructure"
)

var _ core.Plugin = (*AE)(nil)

type AEOptions struct {
	ProjectId int
	Tasks     []string `json:"tasks,omitempty"`
}
type AE string

func (plugin AE) Init() {
	logger.Info("INFO >>> init go plugin", true)
	err := lakeModels.Db.AutoMigrate(
		&models.AEProject{},
		&models.AECommit{})
	if err != nil {
		logger.Error("Error migrating ae: ", err)
		panic(err)
	}
}

func (plugin AE) Description() string {
	return "To collect and enrich data from AE"
}

func (plugin AE) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	logger.Print("start ae plugin execution")

	var op AEOptions
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return err
	}

	tasksToRun := make(map[string]bool, len(op.Tasks))
	for _, task := range op.Tasks {
		tasksToRun[task] = true
	}
	if len(tasksToRun) == 0 {
		tasksToRun = map[string]bool{
			"collectProject": true,
			"collectCommits": true,
			"convertCommits": true,
		}
	}

	progress <- 0.1
	if tasksToRun["collectProject"] {
		if err := tasks.CollectProject(op.ProjectId, ctx); err != nil {
			return &errors.SubTaskError{
				SubTaskName: "collectProject",
				Message:     fmt.Sprintf("could not collect project: %v", err),
			}
		}
	}

	progress <- 0.25

	if tasksToRun["collectCommits"] {
		if err := tasks.CollectCommits(op.ProjectId, ctx); err != nil {
			return &errors.SubTaskError{
				SubTaskName: "collectCommitDevEqs",
				Message:     fmt.Sprintf("could not collect commits: %v", err),
			}
		}
	}

	progress <- 0.75

	if tasksToRun["convertCommits"] {
		if err := tasks.ConvertCommits(ctx); err != nil {
			return &errors.SubTaskError{
				SubTaskName: "convertCommitDevEqs",
				Message:     fmt.Sprintf("could not enhance commits with AE dev equivalent: %v", err),
			}
		}
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

// standalone mode for debugging
func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		panic(fmt.Errorf("Usage: ae <project_id> [sub_task_name]"))
	}
	projectId, err := strconv.ParseInt(args[0], 10, 32)
	if err != nil {
		panic(fmt.Errorf("error paring project_id: %w", err))
	}
	options := map[string]interface{}{
		"projectId": projectId,
	}
	if len(args) > 1 {
		options["tasks"] = []string{args[1]}
	}

	err = core.RegisterPlugin("ae", PluginEntry)
	if err != nil {
		panic(err)
	}
	PluginEntry.Init()
	progress := make(chan float32)
	go func() {
		err := PluginEntry.Execute(
			options,
			progress,
			context.Background(),
		)
		if err != nil {
			panic(err)
		}
		close(progress)
	}()
	for p := range progress {
		fmt.Println(p)
	}
}

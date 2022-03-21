package main

import (
	"context"

<<<<<<< HEAD
	lakeErrors "github.com/merico-dev/lake/errors"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitextractor/tasks"
	"github.com/mitchellh/mapstructure"
=======
	"github.com/merico-dev/lake/plugins/core"
>>>>>>> feat: decouple global variables
)

type GitExtractor struct{}

func (plugin GitExtractor) Description() string {
	return "extract infos from git repository"
}

// TODO: remove
func (plugin GitExtractor) Init() {
<<<<<<< HEAD
	logger := helper.NewDefaultTaskLogger(nil, "git extractor")
	logger.Info("INFO >>> init git extractor")
}

func (plugin GitExtractor) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	logger := helper.NewDefaultTaskLogger(nil, "git extractor")
	logger.Info("INFO >>> start git extractor plugin execution")

	// decode options into op
	var op tasks.GitExtractorOptions
=======
}

func (plugin GitExtractor) Execute(options map[string]interface{}, progress chan<- float32, ctx context.Context) error {
	/* TODO: adopt new interface
	logger.Print("start gitlab plugin execution")
	var op GitExtractorOptions
>>>>>>> feat: decouple global variables
	err := mapstructure.Decode(options, &op)
	if err != nil {
		return err
	}
	err = op.Valid()
	if err != nil {
		return err
	}

	// construct task context
	subtasksToRun := map[string]bool{"CollectGitRepo": true}
	taskCtx := helper.NewDefaultTaskContext("git", ctx, logger, op, subtasksToRun)
	// only 1 subtask to be executed, set current progress to 0 and total to 1
	taskCtx.SetProgress(0, 1)

	// execute subtasks, only one subtask for now
	c, err := taskCtx.SubTaskContext("CollectGitRepo")
	if err != nil {
		return err
	}
<<<<<<< HEAD
	err = tasks.CollectGitRepo(c)
	if err != nil {
		return &lakeErrors.SubTaskError{
			SubTaskName: "collectGitRepo",
			Message: err.Error(),
		}
	}
	taskCtx.IncProgress(1)
=======
	progress <- 1
	*/
>>>>>>> feat: decouple global variables
	return nil
}

func (plugin GitExtractor) RootPkgPath() string {
	return "github.com/merico-dev/lake/plugins/gitextractor"
}

func (plugin GitExtractor) ApiResources() map[string]map[string]core.ApiResourceHandler {
	return nil
}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry GitExtractor //nolint

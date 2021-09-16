package main // must be main for plugin entry point

import (
	"github.com/merico-dev/lake/logger" // A pseudo type for Plugin Interface implementation
	lakeModels "github.com/merico-dev/lake/models"
	gitlabModels "github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/plugins/gitlab/tasks"
	"github.com/merico-dev/lake/utils"
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

	if err := tasks.CollectPipelines(projectIdInt); err != nil {
		logger.Error("Could not collect projects: ", err)
		return
	}

	if err := tasks.CollectProject(projectIdInt); err != nil {
		logger.Error("Could not collect projects: ", err)
		return
	}

	progress <- 0.1

	if err := tasks.CollectCommits(projectIdInt); err != nil {
		logger.Error("Could not collect commits: ", err)
		return
	}

	progress <- 0.3

	mergeRequestErr := tasks.CollectMergeRequests(projectIdInt)
	if mergeRequestErr != nil {
		logger.Error("Could not collect merge requests: ", mergeRequestErr)
		return
	}

	progress <- 0.4

	// find all mrs from db
	var mrs []gitlabModels.GitlabMergeRequest
	lakeModels.Db.Find(&mrs)

	// Gilab's authenticated api rate limit is 2000 per min
	// 15 tasks/s* ~2 requests/task * 60s/min = 1800 per min < 2000 per min
	scheduler, err := utils.NewWorkerScheduler(50, 15)
	if err != nil {
		logger.Error("Could not create work scheduler", err)
	}

	for i := 0; i < len(mrs); i++ {
		mr := (mrs)[i]

		err := scheduler.Submit(func() error {
			notesErr := tasks.CollectMergeRequestNotes(projectIdInt, &mr)
			if notesErr != nil {
				logger.Error("Could not collect MR Notes", notesErr)
				return notesErr
			}

			commitsErr := tasks.CollectMergeRequestCommits(projectIdInt, &mr)
			if commitsErr != nil {
				logger.Error("Could not collect MR Commits", commitsErr)
				return commitsErr
			}
			return nil
		})
		if err != nil {
			logger.Error("err", err)
			return
		}
	}

	scheduler.WaitUntilFinish()
	progress <- 0.8

	enrichErr := tasks.EnrichMergeRequests()
	if enrichErr != nil {
		logger.Error("Could not enrich merge requests", enrichErr)
		return
	}
	progress <- 1

	close(progress)

}

// Export a variable named PluginEntry for Framework to search and load
var PluginEntry Gitlab //nolint

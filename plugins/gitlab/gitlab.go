package main // must be main for plugin entry point

import (
	"github.com/merico-dev/lake/logger" // A pseudo type for Plugin Interface implementation
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

	c1 := make(chan *tasks.ApiMergeRequestResponse)
	go func() {
		mergeRequestErr := tasks.CollectMergeRequests(projectIdInt, c1)
		if mergeRequestErr != nil {
			logger.Error("Could not collect merge requests: ", mergeRequestErr)
			return
		}
	}()
	mergeRequests := <-c1

	scheduler, err := utils.NewWorkerScheduler(50, 10)
	if err != nil {
		logger.Error("Could not create work scheduler", err)
	}

	for _, mergeRequest := range *mergeRequests {
		scheduler.Submit(func() error {
			notesErr := tasks.CollectMergeRequestNotes(projectIdInt, &mergeRequest)
			if notesErr != nil {
				logger.Error("Could not collect MR Notes", notesErr)
				return notesErr
			}

			commitsErr := tasks.CollectMergeRequestCommits(projectIdInt, &mergeRequest)
			if commitsErr != nil {
				logger.Error("Could not collect MR Commits", commitsErr)
				return commitsErr
			}
			return nil
		})
	}

	scheduler.WaitUntilFinish()

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

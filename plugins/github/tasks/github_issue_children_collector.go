package tasks

import (
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/utils"
)

func CollectChildrenOnIssues(owner string, repositoryName string, repositoryId int, scheduler *utils.WorkerScheduler) error {
	var issues []models.GithubIssue
	lakeModels.Db.Find(&issues)
	for i := 0; i < len(issues); i++ {
		issue := (issues)[i]
		eventsErr := CollectIssueEvents(owner, repositoryName, &issue, scheduler)
		if eventsErr != nil {
			logger.Error("Could not collect issue events", eventsErr)
			return eventsErr
		}
		commentsErr := CollectIssueComments(owner, repositoryName, &issue, scheduler)
		if commentsErr != nil {
			logger.Error("Could not collect issue Comments", commentsErr)
			return commentsErr
		}
	}
	return nil
}

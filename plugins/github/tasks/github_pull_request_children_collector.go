package tasks

import (
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/utils"
)

func CollectChildrenOnPullRequests(owner string, repositoryName string, repositoryId int) {
	var prs []models.GithubPullRequest
	lakeModels.Db.Find(&prs)

	maxWorkersPerSecond := 5 // Needs work - this is a temporary value
	scheduler, err := utils.NewWorkerScheduler(50, maxWorkersPerSecond)
	if err != nil {
		logger.Error("Could not create work scheduler for GitHub Pull Requests", err)
	}

	for i := 0; i < len(prs); i++ {
		pr := (prs)[i]
		err := scheduler.Submit(func() error {
			reviewErr := CollectPullRequestReviews(owner, repositoryName, repositoryId, &pr)
			if reviewErr != nil {
				logger.Error("Could not collect PR reviews", reviewErr)
				return reviewErr
			}
			commentsErr := CollectPullRequestComments(owner, repositoryName, &pr)
			if commentsErr != nil {
				logger.Error("Could not collect PR Comments", commentsErr)
				return commentsErr
			}
			commitsErr := CollectPullRequestCommits(owner, repositoryName, &pr)
			if commitsErr != nil {
				logger.Error("Could not collect PR Comments", commitsErr)
				return commitsErr
			}
			
			// This call is to update the details of the individual pull request with additions / deletions / etc.
			// prErr := CollectPullRequest(owner, repositoryName, repositoryId, &pr)
			// if prErr != nil {
			// 	logger.Error("Could not collect PRs to update details", reviewErr)
			// 	return reviewErr
			// }

			return nil
		})
		if err != nil {
			logger.Error("INFO >>> Scheduler error: ", err)
			return
		}
	}

	scheduler.WaitUntilFinish()
}

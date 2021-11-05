package tasks

import (
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/utils"
)

func CollectChildrenOnPullRequests(owner string, repositoryName string, repositoryId int, scheduler *utils.WorkerScheduler, githubApiClient *GithubApiClient) error {
	var prs []models.GithubPullRequest
	lakeModels.Db.Find(&prs)
	for i := 0; i < len(prs); i++ {
		pr := (prs)[i]
		reviewErr := CollectPullRequestReviews(owner, repositoryName, repositoryId, &pr, scheduler, githubApiClient)
		if reviewErr != nil {
			logger.Error("Could not collect PR Reviews", reviewErr)
			return reviewErr
		}
		commentsErr := CollectPullRequestComments(owner, repositoryName, &pr, scheduler, githubApiClient)
		if commentsErr != nil {
			logger.Error("Could not collect PR Comments", commentsErr)
			return commentsErr
		}
		commitsErr := CollectPullRequestCommits(owner, repositoryName, &pr, scheduler, githubApiClient)
		if commitsErr != nil {
			logger.Error("Could not collect PR Commits", commitsErr)
			return commitsErr
		}
		// Please Note: There is no difference between Issue Labels and Pull Request Labels - they are the same.
		labelsErr := CollectIssueLabelsForSinglePullRequest(owner, repositoryName, &pr, scheduler, githubApiClient)
		if labelsErr != nil {
			logger.Error("Could not collect PR Labels", labelsErr)
			return labelsErr
		}
	}
	return nil
}

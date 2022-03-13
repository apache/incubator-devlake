package tasks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/utils"
	"gorm.io/gorm/clause"
)

func CollectPullRequestCommits(ctx context.Context, owner string, repo string, repoId int, rateLimitPerSecondInt int, apiClient *GithubApiClient) error {
	scheduler, err := utils.NewWorkerScheduler(rateLimitPerSecondInt*2, rateLimitPerSecondInt, ctx)
	if err != nil {
		return err
	}
	cursor, err := lakeModels.Db.Model(&models.GithubPullRequest{}).Where("repo_id = ?", repoId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	for cursor.Next() {
		githubPr := &models.GithubPullRequest{}
		err = lakeModels.Db.ScanRows(cursor, githubPr)
		if err != nil {
			return err
		}
		err = scheduler.Submit(func() error {
			processErr := ProcessCollection(owner, repo, githubPr, apiClient)
			if processErr != nil {
				logger.Error("Could not collect PR Commits", err)
				return processErr
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	//scheduler.WaitUntilFinish()

	return nil
}

func ProcessCollection(owner string, repo string, pr *models.GithubPullRequest, apiClient *GithubApiClient) error {
	getUrl := fmt.Sprintf("repos/%v/%v/pulls/%v/commits", owner, repo, pr.Number)
	err := lakeModels.Db.Where("pull_request_id = ?",
		pr.GithubId).Delete(&models.GithubPullRequestCommit{}).Error
	if err != nil {
		logger.Error("Could not delete: ", err)
		return err
	}
	return apiClient.FetchPages(getUrl, nil, 100,
		func(res *http.Response) error {
			githubApiResponse := &ApiPullRequestCommitResponse{}
			if res.StatusCode == 200 {
				err := core.UnmarshalResponse(res, githubApiResponse)
				if err != nil {
					logger.Error("Error: ", err)
					return err
				}
				for _, pullRequestCommit := range *githubApiResponse {
					githubCommit, err := convertPullRequestCommit(&pullRequestCommit)
					if err != nil {
						return err
					}
					err = lakeModels.Db.Clauses(clause.OnConflict{
						DoNothing: true,
					}).Create(&githubCommit).Error
					if err != nil {
						logger.Error("Could not upsert: ", err)
					}
					githubPullRequestCommit := &models.GithubPullRequestCommit{
						CommitSha:     pullRequestCommit.Sha,
						PullRequestId: pr.GithubId,
					}
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&githubPullRequestCommit).Error

					if err != nil {
						logger.Error("Could not upsert: ", err)
					}
				}
			} else {
				fmt.Println("INFO: PR PrCommit collection >>> res.Status: ", res.Status)
			}
			return nil
		})
}

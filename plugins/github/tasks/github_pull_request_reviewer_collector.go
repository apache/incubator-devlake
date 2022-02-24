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

type ApiPullRequestReviewResponse []PullRequestReview

type PullRequestReview struct {
	GithubId int `json:"id"`
	User     struct {
		Id    int
		Login string
	}
	Body        string
	State       string
	SubmittedAt core.Iso8601Time `json:"submitted_at"`
}

func CollectPullRequestReviews(owner string, repo string, apiClient *GithubApiClient, rateLimitPerSecondInt int, ctx context.Context) error {
	scheduler, err := utils.NewWorkerScheduler(rateLimitPerSecondInt*2, rateLimitPerSecondInt, ctx)
	if err != nil {
		return err
	}
	cursor, err := lakeModels.Db.Model(&models.GithubPullRequest{}).Rows()
	if err != nil {
		return nil
	}
	defer cursor.Close()

	for cursor.Next() {
		githubPr := &models.GithubPullRequest{}
		err = lakeModels.Db.ScanRows(cursor, githubPr)
		if err != nil {
			return nil
		}
		err = scheduler.Submit(func() error {
			reviewErr := processPullRequestReviewsCollection(owner, repo, githubPr, apiClient)
			if reviewErr != nil {
				logger.Error("Could not collect PR Reviews", reviewErr)
				return reviewErr
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	scheduler.WaitUntilFinish()

	return nil
}
func processPullRequestReviewsCollection(owner string, repo string, pull *models.GithubPullRequest, apiClient *GithubApiClient) error {
	getUrl := fmt.Sprintf("repos/%v/%v/pulls/%v/reviews", owner, repo, pull.Number)
	return apiClient.FetchPages(getUrl, nil, 100,
		func(res *http.Response) error {
			githubApiResponse := &ApiPullRequestReviewResponse{}
			if res.StatusCode == 200 {
				err := core.UnmarshalResponse(res, githubApiResponse)
				if err != nil {
					logger.Error("Error: ", err)
					return err
				}
				for _, review := range *githubApiResponse {
					githubReviewer := &models.GithubReviewer{
						GithubId:      review.User.Id,
						Login:         review.User.Login,
						PullRequestId: pull.GithubId,
					}
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&githubReviewer).Error
					if err != nil {
						logger.Error("Could not upsert: ", err)
					}
				}
			} else {
				fmt.Println("INFO: PR Review collection >>> res.Status: ", res.Status)
			}
			return nil
		})
}

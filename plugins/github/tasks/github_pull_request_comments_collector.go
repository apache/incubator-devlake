package tasks

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/utils"
	"gorm.io/gorm/clause"
)

type ApiPullRequestCommentResponse []PullRequestComment

type PullRequestComment struct {
	GithubId int `json:"id"`
	Body     string
	User     struct {
		Login string
	}
	GithubCreatedAt core.Iso8601Time `json:"created_at"`
}

func CollectPullRequestComments(owner string, repo string, scheduler *utils.WorkerScheduler, apiClient *GithubApiClient) error {
	cursor, err := lakeModels.Db.Model(&models.GithubPullRequest{}).Rows()
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
		schedulerNew := scheduler

		err = scheduler.Submit(func() error {
			commentsErr := processPullRequestCommentsCollection(owner, repo, githubPr, schedulerNew, apiClient)
			if commentsErr != nil {
				logger.Error("Could not collect PR Comments", commentsErr)
				return commentsErr
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

func processPullRequestCommentsCollection(owner string, repo string, pull *models.GithubPullRequest, scheduler *utils.WorkerScheduler, apiClient *GithubApiClient) error {
	getUrl := fmt.Sprintf("repos/%v/%v/issues/%v/comments", owner, repo, pull.Number)
	return apiClient.FetchPages(getUrl, nil, 100, scheduler,
		func(res *http.Response) error {
			githubApiResponse := &ApiPullRequestCommentResponse{}
			if res.StatusCode == 200 {
				err := core.UnmarshalResponse(res, githubApiResponse)
				if err != nil {
					logger.Error("Error: ", err)
					return err
				}
				for _, comment := range *githubApiResponse {
					githubComment, err := convertGithubPullRequestComment(&comment, pull.GithubId)
					if err != nil {
						return err
					}
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&githubComment).Error
					if err != nil {
						logger.Error("Could not upsert: ", err)
					}
				}
			} else {
				fmt.Println("INFO: PR Comment collection >>> res.Status: ", res.Status)
			}
			return nil
		})
}
func convertGithubPullRequestComment(comment *PullRequestComment, pullId int) (*models.GithubPullRequestComment, error) {
	githubComment := &models.GithubPullRequestComment{
		GithubId:        comment.GithubId,
		PullRequestId:   pullId,
		Body:            comment.Body,
		AuthorUsername:  comment.User.Login,
		GithubCreatedAt: comment.GithubCreatedAt.ToTime(),
	}
	return githubComment, nil
}

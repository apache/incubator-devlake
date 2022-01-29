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

type ApiIssueCommentResponse []IssueComment

type IssueComment struct {
	GithubId int `json:"id"`
	Body     string
	User     struct {
		Login string
	}
	GithubCreatedAt core.Iso8601Time `json:"created_at"`
}

func CollectIssueComments(owner string, repo string, scheduler *utils.WorkerScheduler, apiClient *GithubApiClient) error {
	var issues []models.GithubIssue
	lakeModels.Db.Find(&issues)
	for i := 0; i < len(issues); i++ {
		issue := (issues)[i]

		commentsErr := processCommentsCollection(owner, repo, &issue, scheduler, apiClient)
		if commentsErr != nil {
			logger.Error("Could not collect issue Comments", commentsErr)
			return commentsErr
		}
	}
	return nil
}
func processCommentsCollection(
	owner string,
	repo string,
	issue *models.GithubIssue,
	scheduler *utils.WorkerScheduler,
	apiClient *GithubApiClient,
) error {
	getUrl := fmt.Sprintf("repos/%v/%v/issues/%v/comments", owner, repo, issue.Number)
	return apiClient.FetchWithPaginationAnts(getUrl, nil, 100, 1, scheduler,
		func(res *http.Response) error {
			githubApiResponse := &ApiIssueCommentResponse{}
			if res.StatusCode == 200 {
				err := core.UnmarshalResponse(res, githubApiResponse)
				if err != nil {
					logger.Error("Error: ", err)
					return err
				}
				for _, comment := range *githubApiResponse {
					githubComment, err := convertGithubComment(&comment, issue.GithubId)
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

func convertGithubComment(comment *IssueComment, issueId int) (*models.GithubIssueComment, error) {
	githubComment := &models.GithubIssueComment{
		GithubId:        comment.GithubId,
		IssueId:         issueId,
		Body:            comment.Body,
		AuthorUsername:  comment.User.Login,
		GithubCreatedAt: comment.GithubCreatedAt.ToTime(),
	}
	return githubComment, nil
}

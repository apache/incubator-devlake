package tasks

import (
	"fmt"
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	githubUtils "github.com/merico-dev/lake/plugins/github/utils"
	"gorm.io/gorm/clause"
	"net/http"
)

type ApiIssueCommentResponse []IssueComment

type IssueComment struct {
	GithubId int `json:"id"`
	Body     string
	User     struct {
		Login string
	}
	IssueUrl        string           `json:"issue_url"`
	GithubCreatedAt core.Iso8601Time `json:"created_at"`
}

func CollectIssueComments(owner string, repo string, repoId int, apiClient *GithubApiClient) error {
	commentsErr := processCommentsCollection(owner, repo, repoId, apiClient)
	if commentsErr != nil {
		logger.Error("Could not collect issue Comments", commentsErr)
		return commentsErr
	}
	return nil
}
func processCommentsCollection(
	owner string,
	repo string,
	repoId int,
	apiClient *GithubApiClient,
) error {
	getUrl := fmt.Sprintf("repos/%v/%v/issues/comments", owner, repo)
	return apiClient.FetchPages(getUrl, nil, 100,
		func(res *http.Response) error {
			githubApiResponse := &ApiIssueCommentResponse{}
			if res.StatusCode == 200 {
				err := core.UnmarshalResponse(res, githubApiResponse)
				if err != nil {
					logger.Error("Error: ", err)
					return err
				}
				for _, comment := range *githubApiResponse {
					githubComment, err := convertGithubComment(repoId, &comment)
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

func convertGithubComment(repoId int, comment *IssueComment) (*models.GithubIssueComment, error) {
	issueId, err := githubUtils.GetIssueIdByIssueUrl(comment.IssueUrl)
	if err != nil {
		return nil, err
	}
	githubComment := &models.GithubIssueComment{
		GithubId:        comment.GithubId,
		IssueNumber:     issueId,
		RepoId:          repoId,
		Body:            comment.Body,
		AuthorUsername:  comment.User.Login,
		GithubCreatedAt: comment.GithubCreatedAt.ToTime(),
	}
	return githubComment, nil
}

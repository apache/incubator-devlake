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

type ApiIssuesResponse []IssuesResponse
type IssuesResponse struct {
	GithubId    int `json:"id"`
	Number      int
	State       string
	Title       string
	Body        string
	PullRequest struct {
		Url     string `json:"url"`
		HtmlUrl string `json:"html_url"`
	} `json:"pull_request"`
	ClosedAt        string `json:"closed_at"`
	GithubCreatedAt string `json:"created_at"`
	GithubUpdatedAt string `json:"updated_at"`
}

func CollectIssues(owner string, repositoryName string, repositoryId int, scheduler *utils.WorkerScheduler) error {
	githubApiClient := CreateApiClient()
	getUrl := fmt.Sprintf("repos/%v/%v/issues?state=all", owner, repositoryName)
	return githubApiClient.FetchWithPaginationAnts(getUrl, 100, 20, scheduler,
		func(res *http.Response) error {
			githubApiResponse := &ApiIssuesResponse{}
			err := core.UnmarshalResponse(res, githubApiResponse)
			if err != nil {
				logger.Error("Error: ", err)
				return err
			}

			for _, issue := range *githubApiResponse {
				if issue.PullRequest.Url == "" {
					// This is an issue from github
					githubIssue, err := convertGithubIssue(&issue)
					if err != nil {
						return err
					}
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&githubIssue).Error
					if err != nil {
						logger.Error("Could not upsert: ", err)
					}
				} else {
					// This is a pull request from github
					githubPull, err := convertGithubPullRequest(&issue, repositoryId)
					if err != nil {
						return err
					}
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&githubPull).Error
					if err != nil {
						logger.Error("Could not upsert: ", err)
					}
				}
			}
			return nil
		})
}
func convertGithubIssue(issue *IssuesResponse) (*models.GithubIssue, error) {
	convertedClosedAt := utils.ConvertStringToSqlNullTime(issue.ClosedAt)
	convertedCreatedAt, err := utils.ConvertStringToTime(issue.GithubCreatedAt)
	convertedUpdatedAt, err := utils.ConvertStringToTime(issue.GithubCreatedAt)
	if err != nil {
		return nil, err
	}
	githubIssue := &models.GithubIssue{
		GithubId:        issue.GithubId,
		Number:          issue.Number,
		State:           issue.State,
		Title:           issue.Title,
		Body:            issue.Body,
		ClosedAt:        *convertedClosedAt,
		GithubCreatedAt: *convertedCreatedAt,
		GithubUpdatedAt: *convertedUpdatedAt,
	}
	return githubIssue, nil
}
func convertGithubPullRequest(issue *IssuesResponse, repoId int) (*models.GithubPullRequest, error) {
	convertedCreatedAt, err := utils.ConvertStringToTime(issue.GithubCreatedAt)
	convertedClosedAt := utils.ConvertStringToSqlNullTime(issue.ClosedAt)
	if err != nil {
		return nil, err
	}
	githubPull := &models.GithubPullRequest{
		GithubId:        issue.GithubId,
		RepositoryId:    repoId,
		Number:          issue.Number,
		State:           issue.State,
		Title:           issue.Title,
		GithubCreatedAt: *convertedCreatedAt,
		ClosedAt:        *convertedClosedAt,
	}
	return githubPull, nil
}

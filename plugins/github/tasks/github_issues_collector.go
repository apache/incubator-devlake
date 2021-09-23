package tasks

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/github/utils"
	"gorm.io/gorm/clause"
)

type ApiIssuesResponse []IssuesResponse
type IssuesResponse struct {
	GithubId        int `json:"id"`
	Number          int
	State           string
	Title           string
	Body            string
	ClosedAt        string `json:"closed_at"`
	GithubCreatedAt string `json:"created_at"`
	GithubUpdatedAt string `json:"updated_at"`
}

func CollectIssues(owner string, repositoryName string, repositoryId int) error {
	githubApiClient := CreateApiClient()
	getUrl := fmt.Sprintf("repos/%v/%v/issues", owner, repositoryName)
	return githubApiClient.FetchWithPaginationAnts(getUrl, 100,
		func(res *http.Response) error {
			githubApiResponse := &ApiIssuesResponse{}
			err := core.UnmarshalResponse(res, githubApiResponse)
			if err != nil {
				logger.Error("Error: ", err)
				return err
			}

			for _, issue := range *githubApiResponse {
				githubIssue := &models.GithubIssue{
					GithubId:        issue.GithubId,
					Number:          issue.Number,
					State:           issue.State,
					Title:           issue.Title,
					Body:            issue.Body,
					ClosedAt:        utils.ConvertStringToTime(issue.ClosedAt),
					GithubCreatedAt: utils.ConvertStringToTime(issue.GithubCreatedAt),
					GithubUpdatedAt: utils.ConvertStringToTime(issue.GithubUpdatedAt),
				}

				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&githubIssue).Error
				if err != nil {
					logger.Error("Could not upsert: ", err)
				}
			}
			return nil
		})
}

package tasks

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

type ApiPullsResponse []Pull

type Pull struct {
	GithubId        int `json:"id"`
	State           string
	Title           string
	Number          int
	CommentsUrl     string `json:"comments_url"`
	CommitsUrl      string `json:"commits_url"`
	HTMLUrl         string `json:"html_url"`
	MergedAt        string `json:"merged_at"`
	GithubCreatedAt string `json:"created_at"`
	ClosedAt        string `json:"closed_at"`
	Labels          []Label
}

type Label struct {
	Id          int
	Name        string
	Description string
	Url         string
}

func CollectPullRequests(owner string, repositoryName string, repositoryId int) error {
	githubApiClient := CreateApiClient()
	getUrl := fmt.Sprintf("repos/%v/%v/pulls?state=all", owner, repositoryName)
	return githubApiClient.FetchWithPaginationAnts(getUrl, 100,
		func(res *http.Response) error {
			githubApiResponse := &ApiPullsResponse{}
			err := core.UnmarshalResponse(res, githubApiResponse)
			if err != nil {
				logger.Error("Error: ", err)
				return err
			}
			for _, pull := range *githubApiResponse {
				githubPull := &models.GithubPullRequest{
					GithubId:        pull.GithubId,
					RepositoryId:    repositoryId,
					Number:          pull.Number,
					State:           pull.State,
					Title:           pull.Title,
					CommentsUrl:     pull.CommentsUrl,
					CommitsUrl:      pull.CommitsUrl,
					HTMLUrl:         pull.HTMLUrl,
					MergedAt:        pull.MergedAt,
					GithubCreatedAt: pull.GithubCreatedAt,
					ClosedAt:        pull.ClosedAt,
				}
				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&githubPull).Error
				if err != nil {
					logger.Error("Could not upsert pull request: ", err)
				}
				for _, label := range pull.Labels {
					githubLabel := &models.GithubPullRequestLabel{
						GithubId:      label.Id,
						Name:          label.Name,
						Description:   label.Description,
						Url:           label.Url,
						PullRequestId: pull.GithubId,
					}
					err = lakeModels.Db.Clauses(clause.OnConflict{
						UpdateAll: true,
					}).Create(&githubLabel).Error
					if err != nil {
						logger.Error("Could not upsert label: ", err)
					}
				}
			}
			return nil
		})
}

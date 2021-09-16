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

type ApiCommitsResponse struct {
	GithubId     string `json:"id"`
	RepositoryId int    `json:"repository_id"`
	Title        string
}

func CollectCommits(owner string, repositoryName string) error {
	githubApiClient := CreateApiClient()
	getUrl := fmt.Sprintf("repos/%v/%v/commits", owner, repositoryName)
	return githubApiClient.FetchWithPagination(getUrl, 100,
		func(res *http.Response) error {
			githubApiResponse := &ApiCommitsResponse{}
			err := core.UnmarshalResponse(res, githubApiResponse)
			if err != nil {
				logger.Error("Error: ", err)
				return err
			}
			githubCommits := &models.GithubCommit{
				GithubId:     githubApiResponse.GithubId,
				RepositoryId: githubApiResponse.RepositoryId,
				Title:        githubApiResponse.Title,
			}
			err = lakeModels.Db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&githubCommits).Error
			if err != nil {
				logger.Error("Could not upsert: ", err)
			}
			return nil

		})
}

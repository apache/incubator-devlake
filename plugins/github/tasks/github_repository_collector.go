package tasks

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

type ApiRepositoryResponse struct {
	Name     string `json:"name"`
	GithubId int    `json:"id"`
	HTMLUrl  string `json:"html_url"`
}

func CollectRepository(owner string, repositoryName string) (int, error) {
	githubApiClient := CreateApiClient()
	getUrl := fmt.Sprintf("repos/%v/%v", owner, repositoryName)
	res, err := githubApiClient.Get(getUrl, nil, nil)
	if err != nil {
		logger.Error("Error: ", err)
		return 0, err
	}
	githubApiResponse := &ApiRepositoryResponse{}
	err = core.UnmarshalResponse(res, githubApiResponse)
	if err != nil {
		logger.Error("Error: ", err)
		return 0, err
	}
	githubRepository := &models.GithubRepository{
		Name:     githubApiResponse.Name,
		GithubId: githubApiResponse.GithubId,
		HTMLUrl:  githubApiResponse.HTMLUrl,
	}
	err = lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&githubRepository).Error
	if err != nil {
		logger.Error("Could not upsert: ", err)
	}
	return githubRepository.GithubId, nil
}

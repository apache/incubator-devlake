package tasks

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

type GithubApiRepo struct {
	Name        string `json:"name"`
	GithubId    int    `json:"id"`
	HTMLUrl     string `json:"html_url"`
	Language    string `json:"language"`
	Description string `json:"description"`
	Owner       models.GithubUser
	Parent      *GithubApiRepo    `json:"parent"`
	CreatedAt   core.Iso8601Time  `json:"created_at"`
	UpdatedAt   *core.Iso8601Time `json:"updated_at"`
}

type ApiRepositoryResponse GithubApiRepo

func CollectRepository(owner string, repositoryName string, githubApiClient *GithubApiClient) (int, error) {
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
		GithubId:    githubApiResponse.GithubId,
		Name:        githubApiResponse.Name,
		HTMLUrl:     githubApiResponse.HTMLUrl,
		Description: githubApiResponse.Description,
		OwnerId:     githubApiResponse.Owner.Id,
		OwnerLogin:  githubApiResponse.Owner.Login,
		Language:    githubApiResponse.Language,
		CreatedDate: githubApiResponse.CreatedAt.ToTime(),
		UpdatedDate: core.Iso8601TimeToTime(githubApiResponse.UpdatedAt),
	}
	if githubApiResponse.Parent != nil {
		githubRepository.ParentGithubId = githubApiResponse.Parent.GithubId
		githubRepository.ParentHTMLUrl = githubApiResponse.Parent.HTMLUrl
	}
	err = lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&githubRepository).Error
	if err != nil {
		logger.Error("Could not upsert: ", err)
	}
	err = lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&githubApiResponse.Owner).Error
	if err != nil {
		logger.Error("Could not upsert: ", err)
	}
	return githubRepository.GithubId, nil
}

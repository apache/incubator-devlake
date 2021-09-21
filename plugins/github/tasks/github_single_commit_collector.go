package tasks

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
)

type ApiSingleCommitResponse struct {
	Stats struct {
		Additions int
		Deletions int
	}
}

func CollectCommit(owner string, repositoryName string, repositoryId int, commit *models.GithubCommit) error {
	githubApiClient := CreateApiClient()
	getUrl := fmt.Sprintf("repos/%v/%v/commits/%v", owner, repositoryName, commit.Sha)
	res, getErr := githubApiClient.Get(getUrl, nil, nil)
	if getErr != nil {
		logger.Error("GET Error: ", getErr)
		return getErr
	}

	githubApiResponse := &ApiSingleCommitResponse{}
	unmarshalErr := core.UnmarshalResponse(res, githubApiResponse)
	if unmarshalErr != nil {
		logger.Error("Error: ", unmarshalErr)
		return unmarshalErr
	}
	dbErr := lakeModels.Db.Debug().Model(&commit).Updates(models.GithubCommit{
		Additions: githubApiResponse.Stats.Additions,
		Deletions: githubApiResponse.Stats.Deletions,
	}).Error
	if dbErr != nil {
		logger.Error("Could not update: ", dbErr)
	}
	return nil
}

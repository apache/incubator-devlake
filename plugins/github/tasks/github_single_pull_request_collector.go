package tasks

import (
	"fmt"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/utils"
)

type ApiSinglePullResponse struct {
	Additions      int
	Deletions      int
	Comments       int
	Commits        int
	ReviewComments int `json:"review_comments"`
	Merged         bool
	MergedAt       string `json:"merged_at"`
}

func CollectPullRequest(owner string, repositoryName string, repositoryId int, pr *models.GithubPullRequest) error {
	githubApiClient := CreateApiClient()
	getUrl := fmt.Sprintf("repos/%v/%v/pulls/%v?state=all", owner, repositoryName, pr.Number)
	res, getErr := githubApiClient.Get(getUrl, nil, nil)
	if getErr != nil {
		logger.Error("GET Error: ", getErr)
		return getErr
	}

	githubApiResponse := &ApiSinglePullResponse{}
	unmarshalErr := core.UnmarshalResponse(res, githubApiResponse)
	if unmarshalErr != nil {
		logger.Error("Error: ", unmarshalErr)
		return unmarshalErr
	}
	dbErr := lakeModels.Db.Model(&pr).Updates(models.GithubPullRequest{
		Additions:      githubApiResponse.Additions,
		Deletions:      githubApiResponse.Deletions,
		Comments:       githubApiResponse.Comments,
		Commits:        githubApiResponse.Commits,
		ReviewComments: githubApiResponse.ReviewComments,
		Merged:         githubApiResponse.Merged,
		MergedAt:       utils.ConvertStringToSqlNullTime(githubApiResponse.MergedAt),
	}).Error
	if dbErr != nil {
		logger.Error("Could not update: ", dbErr)
	}
	return nil
}

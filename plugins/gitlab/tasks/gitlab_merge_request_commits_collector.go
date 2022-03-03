package tasks

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"gorm.io/gorm/clause"
)

type ApiMergeRequestCommitResponse []GitlabApiCommit

func CollectMergeRequestCommits(projectId int, mr *models.GitlabMergeRequest, gitlabApiClient *GitlabApiClient) error {

	getUrl := fmt.Sprintf("projects/%v/merge_requests/%v/commits", projectId, mr.Iid)
	return gitlabApiClient.FetchWithPagination(getUrl, nil, 100,
		func(res *http.Response) error {
			gitlabApiResponse := &ApiMergeRequestCommitResponse{}
			err := core.UnmarshalResponse(res, gitlabApiResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}
			for _, commit := range *gitlabApiResponse {
				gitlabCommit, err := ConvertCommit(&commit)
				if err != nil {
					return err
				}
				result := lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&gitlabCommit)

				if result.Error != nil {
					logger.Error("Could not upsert: ", result.Error)
				}
				GitlabMergeRequestCommitMergeRequest := &models.GitlabMergeRequestCommit{
					CommitSha:      commit.GitlabId,
					MergeRequestId: mr.GitlabId,
				}
				result = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&GitlabMergeRequestCommitMergeRequest)

				if result.Error != nil {
					logger.Error("Could not upsert: ", result.Error)
				}
			}

			return nil
		})
}

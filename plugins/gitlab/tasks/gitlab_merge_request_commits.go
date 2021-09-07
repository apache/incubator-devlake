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

// This is just a relationship table between Merge Requests and Commits

type ApiMergeRequestCommitResponse []struct {
	CommitId string `json:"id"`
}

func CollectMergeRequestCommits(projectId int, mr *models.GitlabMergeRequest) error {
	gitlabApiClient := CreateApiClient()

	getUrl := fmt.Sprintf("projects/%v/merge_requests/%v/commits", projectId, mr.Iid)
	return gitlabApiClient.FetchWithPaginationAnts(getUrl, "100",
		func(res *http.Response) error {
			gitlabApiResponse := &ApiMergeRequestCommitResponse{}
			err := core.UnmarshalResponse(res, gitlabApiResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}
			for _, commit := range *gitlabApiResponse {
				gitlabMergeRequestCommit := &models.GitlabMergeRequestCommit{
					MergeRequestId: mr.GitlabId,
					CommitId:       commit.CommitId,
				}
				result := lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&gitlabMergeRequestCommit)

				if result.Error != nil {
					logger.Error("Could not upsert: ", result.Error)
				}
			}

			return nil
		})
}

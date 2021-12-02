package tasks

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/ae/models"
	"github.com/merico-dev/lake/plugins/core"
	"gorm.io/gorm/clause"
)

type ApiCommitResponse []AEApiCommit

type AEApiCommit struct {
	HexSha      string `json:"hexsha"`
	AnalysisId  string `json:"analysis_id"`
	AuthorEmail string `json:"author_email"`
	DevEq       int    `json:"dev_eq"`
}

func CollectCommits(projectId int) error {
	aeApiClient := CreateApiClient()
	relativePath := fmt.Sprintf("projects/%v/commits", projectId)
	pageSize := 2000
	return aeApiClient.FetchWithPagination(relativePath, pageSize,
		func(res *http.Response) error {

			aeApiResponse := &ApiCommitResponse{}
			err := core.UnmarshalResponse(res, aeApiResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}
			for _, aeApiCommit := range *aeApiResponse {
				aeCommit, err := convertCommit(&aeApiCommit, projectId)
				if err != nil {
					return err
				}

				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&aeCommit).Error

				if err != nil {
					logger.Error("Could not upsert: ", err)
				}
			}

			return nil
		})
}

// Convert the API response to our DB model instance
func convertCommit(commit *AEApiCommit, projectId int) (*models.AECommit, error) {
	aeCommit := &models.AECommit{
		HexSha:      commit.HexSha,
		AnalysisId:  commit.AnalysisId,
		AuthorEmail: commit.AuthorEmail,
		DevEq:       commit.DevEq,
		AEProjectId: projectId,
	}
	return aeCommit, nil
}

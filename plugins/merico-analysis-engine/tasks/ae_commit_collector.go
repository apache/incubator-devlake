package tasks

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/merico-dev/lake/config"
	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/merico-analysis-engine/models"
	"github.com/merico-dev/lake/utils"
	"gorm.io/gorm/clause"
)

type ApiCommitResponse []AEApiCommit

type AEApiCommit struct {
	HexSha      string `json:"hexsha"`
	AnalysisId  string `json:"analysis_id"`
	AuthorEmail string `json:"author_email"`
	DevEq       int    `json:"dev_eq"`
}

func CollectCommits(projectId int, scheduler *utils.WorkerScheduler) error {
	aeApiClient := CreateApiClient()
	relativePath := fmt.Sprintf("projects/%v/commits", projectId)
	queryParams := &url.Values{}
	queryParams.Set("app_id", config.V.GetString("AE_APP_ID"))
	queryParams.Set("sign", config.V.GetString("AE_SIGN"))
	return aeApiClient.FetchWithPaginationAnts(scheduler, relativePath, queryParams, 100,
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

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

type ApiPipelineResponse []struct {
	GitlabId        int    `json:"id"`
	ProjectId       int    `json:"project_id"`
	GitlabCreatedAt string `json:"created_at"`
	Status          string
}

func CollectPipelines(projectId int) error {
	gitlabApiClient := CreateApiClient()

	return gitlabApiClient.FetchWithPaginationAnts(fmt.Sprintf("projects/%v/pipelines", projectId), 100,
		func(res *http.Response) error {

			apiPipelineResponse := &ApiPipelineResponse{}
			err := core.UnmarshalResponse(res, apiPipelineResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}

			for _, value := range *apiPipelineResponse {
				gitlabPipeline := &models.GitlabPipeline{
					GitlabId:        value.GitlabId,
					ProjectId:       value.ProjectId,
					GitlabCreatedAt: value.GitlabCreatedAt,
					Status:          value.Status,
				}

				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&gitlabPipeline).Error

				if err != nil {
					logger.Error("Could not upsert: ", err)
				}
			}

			return nil
		})
}

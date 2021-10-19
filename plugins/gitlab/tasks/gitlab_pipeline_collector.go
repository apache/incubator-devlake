package tasks

import (
	"fmt"
	"net/http"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	gitlabModels "github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/utils"
	"gorm.io/gorm/clause"
)

type ApiPipelineResponse []struct {
	GitlabId        int    `json:"id"`
	ProjectId       int    `json:"project_id"`
	GitlabCreatedAt string `json:"created_at"`
	Ref             string
	Sha             string
	WebUrl          string `json:"web_url"`
	Status          string
}

type ApiSinglePipelineResponse struct {
	GitlabId        int    `json:"id"`
	ProjectId       int    `json:"project_id"`
	GitlabCreatedAt string `json:"created_at"`
	Ref             string
	Sha             string
	WebUrl          string `json:"web_url"`
	Duration        int
	StartedAt       string `json:"started_at"`
	FinishedAt      string `json:"finished_at"`
	Coverage        string
	Status          string
}

func CollectAllPipelines(projectId int, scheduler *utils.WorkerScheduler) error {
	gitlabApiClient := CreateApiClient()

	return gitlabApiClient.FetchWithPaginationAnts(scheduler, fmt.Sprintf("projects/%v/pipelines?order_by=updated_at&sort=desc", projectId), 100,
		func(res *http.Response) error {

			apiPipelineResponse := &ApiPipelineResponse{}
			err := core.UnmarshalResponse(res, apiPipelineResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}

			for _, value := range *apiPipelineResponse {
				gitlabPipeline := &gitlabModels.GitlabPipeline{
					GitlabId:        value.GitlabId,
					ProjectId:       value.ProjectId,
					GitlabCreatedAt: utils.ConvertStringToTime(value.GitlabCreatedAt),
					Ref:             value.Ref,
					Sha:             value.Sha,
					WebUrl:          value.WebUrl,
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

func CollectChildrenOnPipelines(projectIdInt int, scheduler *utils.WorkerScheduler) {
	gitlabApiClient := CreateApiClient()

	var pipelines []gitlabModels.GitlabPipeline
	lakeModels.Db.Find(&pipelines)

	for i := 0; i < len(pipelines); i++ {
		pipeline := (pipelines)[i]
		schedulerErr := scheduler.Submit(func() error {

			getUrl := fmt.Sprintf("projects/%v/pipelines/%v", projectIdInt, pipeline.GitlabId)
			res, err := gitlabApiClient.Get(getUrl, nil, nil)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}

			pipelineRes := &ApiSinglePipelineResponse{}
			err2 := core.UnmarshalResponse(res, pipelineRes)

			if err2 != nil {
				logger.Error("Error: ", err2)
				return nil
			}

			gitlabPipeline := &gitlabModels.GitlabPipeline{
				GitlabId:        pipelineRes.GitlabId,
				ProjectId:       pipelineRes.ProjectId,
				GitlabCreatedAt: utils.ConvertStringToTime(pipelineRes.GitlabCreatedAt),
				Ref:             pipelineRes.Ref,
				Sha:             pipelineRes.Sha,
				WebUrl:          pipelineRes.WebUrl,
				Duration:        pipelineRes.Duration,
				StartedAt:       utils.ConvertStringToSqlNullTime(pipelineRes.StartedAt),
				FinishedAt:      utils.ConvertStringToSqlNullTime(pipelineRes.FinishedAt),
				Coverage:        pipelineRes.Coverage,
				Status:          pipelineRes.Status,
			}

			err = lakeModels.Db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&gitlabPipeline).Error

			if err != nil {
				logger.Error("Could not upsert: ", err)
			}
			return nil
		})

		if schedulerErr != nil {
			logger.Error("Error: ", schedulerErr)
		}

	}
	scheduler.WaitUntilFinish()
}

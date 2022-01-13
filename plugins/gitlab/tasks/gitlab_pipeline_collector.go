package tasks

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/merico-dev/lake/logger"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	gitlabModels "github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/utils"
	"gorm.io/gorm/clause"
)

type ApiPipelineResponse []ApiPipeline

type ApiPipeline struct {
	GitlabId        int              `json:"id"`
	ProjectId       int              `json:"project_id"`
	GitlabCreatedAt core.Iso8601Time `json:"created_at"`
	Ref             string
	Sha             string
	WebUrl          string `json:"web_url"`
	Status          string
}

type ApiSinglePipelineResponse struct {
	GitlabId        int              `json:"id"`
	ProjectId       int              `json:"project_id"`
	GitlabCreatedAt core.Iso8601Time `json:"created_at"`
	Ref             string
	Sha             string
	WebUrl          string `json:"web_url"`
	Duration        int
	StartedAt       *core.Iso8601Time `json:"started_at"`
	FinishedAt      *core.Iso8601Time `json:"finished_at"`
	Coverage        string
	Status          string
}

var pipelinesSlice = []models.GitlabPipeline{}

func CollectAllPipelines(projectId int, scheduler *utils.WorkerScheduler) error {
	gitlabApiClient := CreateApiClient()
	queryParams := &url.Values{}
	queryParams.Set("order_by", "updated_at")
	queryParams.Set("sort", "desc")
	return gitlabApiClient.FetchWithPaginationAnts(scheduler,
		fmt.Sprintf("projects/%v/pipelines", projectId),
		queryParams,
		100,
		func(res *http.Response) error {

			apiPipelineResponse := &ApiPipelineResponse{}
			err := core.UnmarshalResponse(res, apiPipelineResponse)

			if err != nil {
				logger.Error("Error: ", err)
				return nil
			}

			for _, pipeline := range *apiPipelineResponse {

				gitlabPipeline, err := convertPipeline(&pipeline)
				if err != nil {
					return err
				}
				pipelinesSlice = append(pipelinesSlice, *gitlabPipeline)
			}
			err = savePipelinesInBatches()
			if err != nil {
				logger.Error("Error: ", err)
				return err
			}
			return nil
		})
}

func savePipelinesInBatches() error {
	err := lakeModels.Db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&pipelinesSlice).Error
	if err != nil {
		return err
	}
	return nil
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

			gitlabPipeline, err := convertSinglePipeline(pipelineRes)
			if err != nil {
				return err
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

func convertSinglePipeline(pipeline *ApiSinglePipelineResponse) (*models.GitlabPipeline, error) {

	gitlabPipeline := &gitlabModels.GitlabPipeline{
		GitlabId:        pipeline.GitlabId,
		ProjectId:       pipeline.ProjectId,
		GitlabCreatedAt: pipeline.GitlabCreatedAt.ToTime(),
		Ref:             pipeline.Ref,
		Sha:             pipeline.Sha,
		WebUrl:          pipeline.WebUrl,
		Duration:        pipeline.Duration,
		StartedAt:       core.Iso8601TimeToTime(pipeline.StartedAt),
		FinishedAt:      core.Iso8601TimeToTime(pipeline.FinishedAt),
		Coverage:        pipeline.Coverage,
		Status:          pipeline.Status,
	}
	return gitlabPipeline, nil
}

func convertPipeline(pipeline *ApiPipeline) (*models.GitlabPipeline, error) {
	gitlabPipeline := &gitlabModels.GitlabPipeline{
		GitlabId:        pipeline.GitlabId,
		ProjectId:       pipeline.ProjectId,
		GitlabCreatedAt: pipeline.GitlabCreatedAt.ToTime(),
		Ref:             pipeline.Ref,
		Sha:             pipeline.Sha,
		WebUrl:          pipeline.WebUrl,
		Status:          pipeline.Status,
	}
	return gitlabPipeline, nil
}

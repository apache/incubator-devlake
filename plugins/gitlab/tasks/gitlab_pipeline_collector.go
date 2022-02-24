package tasks

import (
	"fmt"
	"net/http"
	"net/url"

	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	gitlabModels "github.com/merico-dev/lake/plugins/gitlab/models"
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

func CollectAllPipelines(projectId int, gitlabApiClient *GitlabApiClient) error {

	queryParams := &url.Values{}
	queryParams.Set("order_by", "updated_at")
	queryParams.Set("sort", "desc")
	return gitlabApiClient.FetchWithPaginationAnts(
		fmt.Sprintf("projects/%v/pipelines", projectId),
		queryParams,
		100,
		func(res *http.Response) error {

			apiPipelineResponse := &ApiPipelineResponse{}
			err := core.UnmarshalResponse(res, apiPipelineResponse)

			if err != nil {
				return err
			}

			for _, pipeline := range *apiPipelineResponse {

				gitlabPipeline, err := convertPipeline(&pipeline)
				if err != nil {
					return err
				}
				err = lakeModels.Db.Clauses(clause.OnConflict{
					UpdateAll: true,
				}).Create(&gitlabPipeline).Error

				if err != nil {
					return err

				}
			}

			return nil
		})
}

func CollectChildrenOnPipelines(projectIdInt int, gitlabApiClient *GitlabApiClient) error {

	var pipelines []gitlabModels.GitlabPipeline

	//Find all piplines associated with the current projectid
	lakeModels.Db.Where("project_id=?", projectIdInt).Find(&pipelines)

	for i := 0; i < len(pipelines); i++ {
		pipeline := (pipelines)[i]
		getUrl := fmt.Sprintf("projects/%v/pipelines/%v", projectIdInt, pipeline.GitlabId)

		err := gitlabApiClient.GetAsync(getUrl, nil, func(res *http.Response) error {
			// Check the StatusCode of the HTTP response
			if res.StatusCode != 200 {
				return fmt.Errorf("got a bad response StatusCode [%d] when requesting [%s]", res.StatusCode, getUrl)
			}

			pipelineRes := &ApiSinglePipelineResponse{}
			err := core.UnmarshalResponse(res, pipelineRes)

			if err != nil {
				return err
			}

			gitlabPipeline, err := convertSinglePipeline(pipelineRes)
			if err != nil {
				return err
			}

			// use projectIdInt to set the value of ProjectId for it
			gitlabPipeline.ProjectId = projectIdInt

			err = lakeModels.Db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(&gitlabPipeline).Error

			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			return err
		}

	}
	gitlabApiClient.WaitOtherGoroutines()
	return nil
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

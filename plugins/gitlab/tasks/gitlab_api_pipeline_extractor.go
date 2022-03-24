package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/plugins/helper"
)

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

var ExtractApiPipelinesMeta = core.SubTaskMeta{
	Name:             "extractApiPipelines",
	EntryPoint:       ExtractApiPipelines,
	EnabledByDefault: true,
	Description:      "Extract raw pipelines data into tool layer table GitlabPipeline",
}

var ExtractApiChildrenOnPipelinesMeta = core.SubTaskMeta{
	Name:             "extractApiChildrenOnPipelines",
	EntryPoint:       ExtractApiChildrenOnPipelines,
	EnabledByDefault: true,
	Description:      "Extract raw pipelines data into tool layer table GitlabPipeline",
}

func ExtractApiPipelines(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			// create gitlab commit
			gitlabApiPipeline := &ApiPipeline{}
			err := json.Unmarshal(row.Data, gitlabApiPipeline)
			if err != nil {
				return nil, err
			}
			gitlabPipeline, err := convertPipeline(gitlabApiPipeline)
			if err != nil {
				return nil, err
			}

			// use data.Options.ProjectId to set the value of ProjectId for it
			gitlabPipeline.ProjectId = data.Options.ProjectId

			results := make([]interface{}, 0, 1)
			results = append(results, gitlabPipeline)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func ExtractApiChildrenOnPipelines(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_CHILDREN_ON_PIPELINE_TABLE)

	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			pipelineRes := &ApiSinglePipelineResponse{}
			err := json.Unmarshal(row.Data, pipelineRes)
			if err != nil {
				return nil, err
			}

			gitlabPipeline, err := convertSinglePipeline(pipelineRes)
			if err != nil {
				return nil, err
			}

			// use data.Options.ProjectId to set the value of ProjectId for it
			gitlabPipeline.ProjectId = data.Options.ProjectId

			results := make([]interface{}, 0, 1)
			results = append(results, gitlabPipeline)

			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

func convertSinglePipeline(pipeline *ApiSinglePipelineResponse) (*models.GitlabPipeline, error) {
	gitlabPipeline := &models.GitlabPipeline{
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
	gitlabPipeline := &models.GitlabPipeline{
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

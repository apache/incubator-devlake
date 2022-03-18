package tasks

import (
	"encoding/json"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

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

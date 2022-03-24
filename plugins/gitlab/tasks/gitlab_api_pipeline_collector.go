package tasks

import (
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

const RAW_PIPELINE_TABLE = "gitlab_api_pipeline"
const RAW_CHILDREN_ON_PIPELINE_TABLE = "gitlab_api_children_on_pipeline"

var CollectApiPipelinesMeta = core.SubTaskMeta{
	Name:             "collectApiPipelines",
	EntryPoint:       CollectApiPipelines,
	EnabledByDefault: true,
	Description:      "Collect pipeline data from gitlab api",
}

var CollectApiChildrenOnPipelinesMeta = core.SubTaskMeta{
	Name:             "collectApiChildrenOnPipelines",
	EntryPoint:       CollectApiChildrenOnPipelines,
	EnabledByDefault: true,
	Description:      "Collect pipline child data from gitlab api",
}

func CollectApiPipelines(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PIPELINE_TABLE)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		Incremental:        false,
		UrlTemplate:        "projects/{{ .Params.ProjectId }}/pipelines",
		Query:              GetQueryOrder,
		GetTotalPages:      GetTotalPagesFromResponse,
		ResponseParser:     GetRawMessageFromResponse,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

func CollectApiChildrenOnPipelines(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_CHILDREN_ON_PIPELINE_TABLE)

	iterator, err := GetPipelinesIterator(taskCtx)
	if err != nil {
		return err
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		Incremental:        false,
		Input:              iterator,
		UrlTemplate:        "projects/{{ .Params.ProjectId }}/pipelines/{{ .Input.GitlabId }}",
		Query:              GetQueryOrder,
		ResponseParser:     helper.GetRawMessageDirectFromResponse,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

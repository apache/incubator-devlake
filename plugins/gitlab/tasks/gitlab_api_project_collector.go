package tasks

import (
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

const RAW_PROJECT_TABLE = "gitlab_api_project"

func CollectApiProject(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_PROJECT_TABLE)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		Incremental:        false,
		UrlTemplate:        "projects/{{ .Params.ProjectId }}",
		Query:              GetQuery,
		ResponseParser:     helper.GetRawMessageDirectFromResponse,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

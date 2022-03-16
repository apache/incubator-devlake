package tasks

import (
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

const RAW_MERGE_REQUEST_TABLE = "gitlab_api_merge_requests"

func CollectApiMergeRequests(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GitlabTaskData)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GitlabApiParams{
				ProjectId: data.Options.ProjectId,
			},
			Table: RAW_MERGE_REQUEST_TABLE,
		},
		ApiClient:      data.ApiClient,
		PageSize:       100,
		Incremental:    false,
		UrlTemplate:    "projects/{{ .Params.ProjectId }}/merge_requests",
		Query:          GetQuery,
		GetTotalPages:  GetTotalPagesFromResponse,
		ResponseParser: GetRawMessageFromResponse,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

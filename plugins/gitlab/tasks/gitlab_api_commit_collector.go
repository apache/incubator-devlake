package tasks

import (
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

const RAW_COMMIT_TABLE = "gitlab_api_commit"

func CollectApiCommits(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*GitlabTaskData)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GitlabApiParams{
				ProjectId: data.Options.ProjectId,
			},
			Table: RAW_COMMIT_TABLE,
		},
		ApiClient:      data.ApiClient,
		PageSize:       100,
		Incremental:    false,
		UrlTemplate:    "projects/{{ .Params.ProjectId }}/repository/commits",
		Query:          GetQuery,
		ResponseParser: GetRawMessageFromResponse,
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

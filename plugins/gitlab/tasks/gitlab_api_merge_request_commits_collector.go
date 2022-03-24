package tasks

import (
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

const RAW_MERGE_REQUEST_COMMITS_TABLE = "gitlab_api_merge_request_commits"

var CollectApiMergeRequestsCommitsMeta = core.SubTaskMeta{
	Name:             "collectApiMergeRequestsCommits",
	EntryPoint:       CollectApiMergeRequestsCommits,
	EnabledByDefault: true,
	Description:      "Collect merge requests commits data from gitlab api",
}

func CollectApiMergeRequestsCommits(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_COMMITS_TABLE)

	iterator, err := GetMergeRequestsIterator(taskCtx)
	if err != nil {
		return err
	}

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		ApiClient:          data.ApiClient,
		PageSize:           100,
		Incremental:        false,
		Input:              iterator,
		UrlTemplate:        "projects/{{ .Params.ProjectId }}/merge_requests/{{ .Input.Iid }}/commits",
		Query:              GetQuery,
		GetTotalPages:      GetTotalPagesFromResponse,
		ResponseParser:     GetRawMessageFromResponse,
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

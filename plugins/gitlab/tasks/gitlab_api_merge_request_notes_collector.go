package tasks

import (
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

const RAW_MERGE_REQUEST_NOTES_TABLE = "gitlab_api_merge_request_notes"

var CollectApiMergeRequestsNotesMeta = core.SubTaskMeta{
	Name:             "collectApiMergeRequestsNotes",
	EntryPoint:       CollectApiMergeRequestsNotes,
	EnabledByDefault: true,
	Description:      "Collect merge requests notes data from gitlab api",
}

func CollectApiMergeRequestsNotes(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_NOTES_TABLE)

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
		UrlTemplate:        "projects/{{ .Params.ProjectId }}/merge_requests/{{ .Input.Iid }}/notes?system=false",
		Query:              GetQuery,
		GetTotalPages:      GetTotalPagesFromResponse,
		ResponseParser:     GetRawMessageFromResponse,
	})
	if err != nil {
		return err
	}

	return collector.Execute()
}

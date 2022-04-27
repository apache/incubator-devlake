package tasks

import (
	"fmt"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"net/url"
)

const RAW_BUG_STATUS_TABLE = "tapd_api_bug_status"

var _ core.SubTaskEntryPoint = CollectBugStatus

func CollectBugStatus(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect bugStatus")

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_BUG_STATUS_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    100,
		UrlTemplate: "workflows/status_map",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceID))
			query.Set("system", "bug")
			return query, nil
		},
		ResponseParser: GetRawMessageDirectFromResponse,
	})
	if err != nil {
		logger.Error("collect bugStatus error:", err)
		return err
	}
	return collector.Execute()
}

var CollectBugStatusMeta = core.SubTaskMeta{
	Name:        "collectBugStatus",
	EntryPoint:  CollectBugStatus,
	Required:    true,
	Description: "collect Tapd bugStatus",
}

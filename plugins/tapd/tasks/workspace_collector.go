package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"net/http"
	"net/url"
)

const RAW_WORKSPACE_TABLE = "tapd_api_workspaces"

var _ core.SubTaskEntryPoint = CollectWorkspaces

type TapdApiParams struct {
	ConnectionId uint64
	CompanyId    uint64
	WorkspaceID  uint64
}

func CollectWorkspaces(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect workspaces")
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_WORKSPACE_TABLE,
		},
		ApiClient: data.ApiClient,
		//PageSize:    100,
		UrlTemplate: "workspaces/sub_workspaces",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceID))
			//query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			//query.Set("limit", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		//GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Workspaces []json.RawMessage `json:"data"`
			}
			err := helper.UnmarshalResponse(res, &data)
			return data.Workspaces, err
		},
	})
	if err != nil {
		logger.Error("collect workspace error:", err)
		return err
	}
	return collector.Execute()
}

var CollectWorkspaceMeta = core.SubTaskMeta{
	Name:        "collectWorkspaces",
	EntryPoint:  CollectWorkspaces,
	Required:    true,
	Description: "collect Tapd workspaces",
}

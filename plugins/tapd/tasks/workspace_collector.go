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
	SourceId    uint64
	CompanyId   uint64
	WorkspaceId uint64
}

func CollectWorkspaces(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect workspaces")
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId:  data.Source.ID,
				CompanyId: data.Source.CompanyId,
			},
			Table: RAW_WORKSPACE_TABLE,
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: "workspaces/projects",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("company_id", fmt.Sprintf("%v", data.Source.CompanyId))
			return query, nil
		},
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

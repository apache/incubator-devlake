package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"net/http"
	"net/url"
)

const RAW_USER_TABLE = "tapd_api_users"

var _ core.SubTaskEntryPoint = CollectUsers

func CollectUsers(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect users")
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_USER_TABLE,
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: "workspaces/users",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Users []json.RawMessage `json:"data"`
			}
			err := helper.UnmarshalResponse(res, &data)
			return data.Users, err
		},
	})
	if err != nil {
		logger.Error("collect user error:", err)
		return err
	}
	return collector.Execute()
}

var CollectUserMeta = core.SubTaskMeta{
	Name:        "collectUsers",
	EntryPoint:  CollectUsers,
	Required:    true,
	Description: "collect Tapd users",
}

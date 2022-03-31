package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"net/http"
	"net/url"
)

const RAW_USER_ROLE_TABLE = "tapd_api_user_roles"

var _ core.SubTaskEntryPoint = CollectUserRoles

func CollectUserRoles(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect userRoles")
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_USER_ROLE_TABLE,
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: "roles",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				UserRoles []json.RawMessage `json:"data"`
			}
			err := helper.UnmarshalResponse(res, &data)
			return data.UserRoles, err
		},
	})
	if err != nil {
		logger.Error("collect userRole error:", err)
		return err
	}
	return collector.Execute()
}

var CollectUserRoleMeta = core.SubTaskMeta{
	Name:        "collectUserRoles",
	EntryPoint:  CollectUserRoles,
	Required:    true,
	Description: "collect Tapd userRoles",
}

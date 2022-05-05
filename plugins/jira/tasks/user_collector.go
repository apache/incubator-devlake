package tasks

import (
	"encoding/json"
	"net/http"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

const RAW_USERS_TABLE = "jira_api_users"

func CollectUsers(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect user")
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_USERS_TABLE,
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: "api/3/users/search",
		PageSize:    50,
		Concurrency: 10,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var result []json.RawMessage
			err := helper.UnmarshalResponse(res, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		},
	})
	if err != nil {
		logger.Error("collect user error:", err)
		return err
	}

	return collector.Execute()
}

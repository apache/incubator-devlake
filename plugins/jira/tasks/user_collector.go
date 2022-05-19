package tasks

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/jira/models"
)

const RAW_USERS_TABLE = "jira_api_users"

func CollectUsers(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	db := taskCtx.GetDb()
	logger := taskCtx.GetLogger()
	logger.Info("collect user")
	cursor, err := db.Model(&models.JiraUser{}).Where("connection_id = ?", data.Options.ConnectionId).Rows()
	if err != nil {
		return err
	}
	iterator, err := helper.NewCursorIterator(db, cursor, reflect.TypeOf(models.JiraUser{}))
	if err != nil {
		return err
	}
	queryKey := "accountId"
	if data.JiraServerInfo.DeploymentType == models.DeploymentServer {
		queryKey = "username"
	}
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
		Input:       iterator,
		UrlTemplate: "api/2/user",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			user := reqData.Input.(*models.JiraUser)
			query := url.Values{}
			query.Set(queryKey, user.AccountId)
			return query, nil
		},
		Concurrency: 10,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var result json.RawMessage
			err := helper.UnmarshalResponse(res, &result)
			if err != nil {
				return nil, err
			}
			return []json.RawMessage{result}, nil
		},
	})
	if err != nil {
		logger.Error("collect user error:", err)
		return err
	}

	return collector.Execute()
}

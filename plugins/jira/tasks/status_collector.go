package tasks

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_STATUS_TABLE = "jira_api_status"

func CollectStatus(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)

	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_STATUS_TABLE,
		},
		ApiClient:     data.ApiClient,
		PageSize:      100,
		Incremental:   false,
		UrlTemplate:   "rest/api/2/status",
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data []json.RawMessage
			blob, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(blob, &data)
			if err != nil {
				return nil, err
			}
			return data, nil
		},
	})

	if err != nil {
		return err
	}

	return collector.Execute()
}

package tasks

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_BOARD_TABLE = "jira_api_boards"

var _ core.SubTaskEntryPoint = CollectBoard

func CollectBoard(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect board:%d", data.Options.BoardId)
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				ConnectionId: data.Connection.ID,
				BoardId:      data.Options.BoardId,
			},
			Table: RAW_BOARD_TABLE,
		},
		ApiClient:     data.ApiClient,
		UrlTemplate:   "agile/1.0/board/{{ .Params.BoardId }}",
		GetTotalPages: GetTotalPagesFromResponse,
		Concurrency:   10,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			blob, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}
			res.Body.Close()
			return []json.RawMessage{blob}, nil
		},
	})
	if err != nil {
		logger.Error("collect board error:", err)
		return err
	}

	return collector.Execute()
}

package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"net/http"
	"net/url"
)

const RAW_ITERATION_TABLE = "tapd_api_iterations"

var _ core.SubTaskEntryPoint = CollectIterations

func CollectIterations(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect iterations")
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_ITERATION_TABLE,
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: "iterations",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("workspace_id", fmt.Sprintf("%v", data.Options.WorkspaceId))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var data struct {
				Iterations []json.RawMessage `json:"data"`
			}
			err := helper.UnmarshalResponse(res, &data)
			return data.Iterations, err
		},
	})
	if err != nil {
		logger.Error("collect iteration error:", err)
		return err
	}
	return collector.Execute()
}

var CollectIterationMeta = core.SubTaskMeta{
	Name:        "collectIterations",
	EntryPoint:  CollectIterations,
	Required:    true,
	Description: "collect Tapd iterations",
}

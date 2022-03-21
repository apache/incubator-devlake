package tasks

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

const RAW_PROJECT_TABLE = "jira_api_projects"

var _ core.SubTaskEntryPoint = CollectProjects

func CollectProjects(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*JiraTaskData)
	logger := taskCtx.GetLogger()
	logger.Info("collect projects")
	jql := "ORDER BY created ASC"
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: JiraApiParams{
				SourceId: data.Source.ID,
			},
			Table: RAW_PROJECT_TABLE,
		},
		ApiClient:   data.ApiClient,
		UrlTemplate: "api/2/project",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("jql", jql)
			return query, nil
		},
		GetTotalPages: GetTotalPagesFromResponse,
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			var result []json.RawMessage
			err := helper.UnmarshalResponse(res, &result)
			return result, err
		},
	})
	if err != nil {
		logger.Error("collect project error:", err)
		return err
	}
	return collector.Execute()
}

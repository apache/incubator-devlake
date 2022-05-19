package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
)

const RAW_COMMITS_TABLE = "ae_commits"

func CollectCommits(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*AeTaskData)
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: AeApiParams{
				ProjectId: data.Options.ProjectId,
			},
			Table: RAW_COMMITS_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    2000,
		UrlTemplate: "projects/{{ .Params.ProjectId }}/commits",
		Query: func(reqData *helper.RequestData) (url.Values, error) {
			query := url.Values{}
			query.Set("page", fmt.Sprintf("%v", reqData.Pager.Page))
			query.Set("per_page", fmt.Sprintf("%v", reqData.Pager.Size))
			return query, nil
		},
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}
			var results []json.RawMessage
			err = json.Unmarshal(body, &results)
			return results, err
		},
	})

	if err != nil {
		return err
	}
	return collector.Execute()
}

var CollectCommitsMeta = core.SubTaskMeta{
	Name:             "collectCommits",
	EntryPoint:       CollectCommits,
	EnabledByDefault: true,
	Description:      "Collect commit analysis data from AE api",
}

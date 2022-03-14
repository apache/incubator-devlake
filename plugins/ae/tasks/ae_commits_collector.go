package tasks

import (
	"fmt"
	"net/url"

	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

func CollectCommits(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*AeTaskData)
	collector, err := helper.NewApiCollector(helper.ApiCollectorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: AeApiParams{
				ProjectId: data.Options.ProjectId,
			},
			Table: RAW_PROJECT_TABLE,
		},
		ApiClient:   data.ApiClient,
		PageSize:    2000,
		UrlTemplate: "projects/%v/commits",
		Query: func(pager *helper.Pager) (url.Values, error) {
			query := url.Values{}
			query.Set("page", fmt.Sprintf("%v", pager.Page))
			query.Set("per_page", fmt.Sprintf("%v", pager.Size))
			return query, nil
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

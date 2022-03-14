package tasks

import (
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

const RAW_PROJECT_TABLE = "ae_project"

func CollectProject(taskCtx core.SubTaskContext) error {
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
		PageSize:    100,
		UrlTemplate: "projects/{{ .Params.ProjectId }}",
	})

	if err != nil {
		return err
	}
	return collector.Execute()
}

var CollectProjectMeta = core.SubTaskMeta{
	Name:             "collectProject",
	EntryPoint:       CollectProject,
	EnabledByDefault: true,
	Description:      "Collect analysis project data from AE api",
}

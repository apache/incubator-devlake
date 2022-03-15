package tasks

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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
		UrlTemplate: "projects/{{ .Params.ProjectId }}",
		ResponseParser: func(res *http.Response) ([]json.RawMessage, error) {
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}
			res.Body.Close()
			return []json.RawMessage{
				json.RawMessage(body),
			}, nil
		},
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

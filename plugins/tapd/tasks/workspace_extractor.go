package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"gorm.io/datatypes"
)

var _ core.SubTaskEntryPoint = ExtractWorkspaces

var ExtractWorkspacesMeta = core.SubTaskMeta{
	Name:             "extractWorkspaces",
	EntryPoint:       ExtractWorkspaces,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table tapd_workspaces",
}

type TapdWorkspaceRes struct {
	Workspace datatypes.JSON
}

func ExtractWorkspaces(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId:  data.Source.ID,
				CompanyId: data.Options.CompanyId,
			},
			Table: RAW_WORKSPACE_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var workspaceRes TapdWorkspaceRes
			err := json.Unmarshal(row.Data, &workspaceRes)
			if err != nil {
				return nil, err
			}
			var workspace models.TapdWorkspace
			err = json.Unmarshal(workspaceRes.Workspace, &workspace)
			if err != nil {
				return nil, err
			}
			results := make([]interface{}, 0, 1)
			results = append(results, &workspace)
			return results, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

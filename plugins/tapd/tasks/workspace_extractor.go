package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractWorkspaces

var ExtractWorkspaceMeta = core.SubTaskMeta{
	Name:             "extractWorkspaces",
	EntryPoint:       ExtractWorkspaces,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_workspaces",
}

type TapdWorkspaceRes struct {
	Workspace models.TapdWorkspaceApiRes
}

func ExtractWorkspaces(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_WORKSPACE_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var workspaceRes TapdWorkspaceRes
			err := json.Unmarshal(row.Data, &workspaceRes)
			if err != nil {
				return nil, err
			}

			wsRes := workspaceRes.Workspace

			i, err := VoToDTO(&wsRes, &models.TapdWorkspace{})
			if err != nil {
				return nil, err
			}
			ws := i.(*models.TapdWorkspace)
			ws.SourceId = data.Source.ID
			return []interface{}{
				ws,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

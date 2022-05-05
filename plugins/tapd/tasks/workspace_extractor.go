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
	Workspace models.TapdWorkspace
}

func ExtractWorkspaces(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_WORKSPACE_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var workspaceRes TapdWorkspaceRes
			err := json.Unmarshal(row.Data, &workspaceRes)
			if err != nil {
				return nil, err
			}

			ws := workspaceRes.Workspace

			ws.ConnectionId = data.Connection.ID
			return []interface{}{
				&ws,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

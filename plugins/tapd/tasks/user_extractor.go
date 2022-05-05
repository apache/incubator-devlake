package tasks

import (
	"encoding/json"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
)

var _ core.SubTaskEntryPoint = ExtractUsers

var ExtractUserMeta = core.SubTaskMeta{
	Name:             "extractUsers",
	EntryPoint:       ExtractUsers,
	EnabledByDefault: true,
	Description:      "Extract raw workspace data into tool layer table _tool_tapd_users",
}

type TapdUserRes struct {
	UserWorkspace models.TapdUser
}

func ExtractUsers(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	extractor, err := helper.NewApiExtractor(helper.ApiExtractorArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,
				//CompanyId: data.Options.CompanyId,
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_USER_TABLE,
		},
		Extract: func(row *helper.RawData) ([]interface{}, error) {
			var userRes TapdUserRes
			err := json.Unmarshal(row.Data, &userRes)
			if err != nil {
				return nil, err
			}
			toolL := models.TapdUser{
				ConnectionId: data.Connection.ID,
				WorkspaceID:  data.Options.WorkspaceID,
				Name:         userRes.UserWorkspace.Name,
				User:         userRes.UserWorkspace.User,
			}
			return []interface{}{
				&toolL,
			}, nil
		},
	})

	if err != nil {
		return err
	}

	return extractor.Execute()
}

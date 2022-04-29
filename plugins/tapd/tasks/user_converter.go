package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/user"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"reflect"
)

func ConvertUser(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("convert board:%d", data.Options.WorkspaceID)
	cursor, err := db.Model(&models.TapdUser{}).Where("connection_id = ? AND workspace_id = ?", data.Connection.ID, data.Options.WorkspaceID).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				ConnectionId: data.Connection.ID,

				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_USER_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdUser{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			userTool := inputRow.(*models.TapdUser)
			issue := &user.User{
				DomainEntity: domainlayer.DomainEntity{
					Id: UserIdGen.Generate(userTool.ConnectionId, userTool.WorkspaceID, userTool.User),
				},
				Name: userTool.Name,
			}

			return []interface{}{
				issue,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

var ConvertUserMeta = core.SubTaskMeta{
	Name:             "convertUser",
	EntryPoint:       ConvertUser,
	EnabledByDefault: true,
	Description:      "convert Tapd User",
}

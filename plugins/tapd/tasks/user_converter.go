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
	logger.Info("convert board:%d", data.Options.WorkspaceId)
	cursor, err := db.Model(&models.TapdUser{}).Where("source_id = ? AND workspace_id = ?", data.Source.ID, data.Options.WorkspaceId).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: TapdApiParams{
				SourceId: data.Source.ID,
				//CompanyId:   data.Source.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: RAW_USER_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdUser{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			userTool := inputRow.(*models.TapdUser)
			issue := &user.User{
				DomainEntity: domainlayer.DomainEntity{
					Id: UserIdGen.Generate(userTool.SourceId, userTool.WorkspaceId, userTool.Name),
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

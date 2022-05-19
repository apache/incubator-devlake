package tasks

import (
	"fmt"
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"reflect"
)

func ConvertWorkspace(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("collect board:%d", data.Options.WorkspaceID)
	cursor, err := db.Model(&models.TapdWorkspace{}).Where("connection_id = ? AND id = ?", data.Connection.ID, data.Options.WorkspaceID).Rows()
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
			Table: RAW_WORKSPACE_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdWorkspace{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			workspace := inputRow.(*models.TapdWorkspace)
			domainBoard := &ticket.Board{
				DomainEntity: domainlayer.DomainEntity{
					Id: WorkspaceIdGen.Generate(workspace.ConnectionId, workspace.ID),
				},
				Name: workspace.Name,
				Url:  fmt.Sprintf("%s/%d", "https://tapd.cn", workspace.ID),
			}
			return []interface{}{
				domainBoard,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

var ConvertWorkspaceMeta = core.SubTaskMeta{
	Name:             "convertWorkspace",
	EntryPoint:       ConvertWorkspace,
	EnabledByDefault: true,
	Description:      "convert Tapd workspace",
}

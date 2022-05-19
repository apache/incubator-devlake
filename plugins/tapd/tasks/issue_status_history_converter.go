package tasks

//import (
//	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
//	"github.com/apache/incubator-devlake/plugins/core"
//	"github.com/apache/incubator-devlake/plugins/helper"
//	"github.com/apache/incubator-devlake/plugins/tapd/models"
//	"reflect"
//)
//
//func ConvertIssueStatusHistory(taskCtx core.SubTaskContext) error {
//	data := taskCtx.GetData().(*TapdTaskData)
//	logger := taskCtx.GetLogger()
//	db := taskCtx.GetDb()
//	logger.Info("convert board:%d", data.Options.WorkspaceID)
//	cursor, err := db.Model(&models.TapdIssueStatusHistory{}).Where("connection_id = ? AND workspace_id = ?", data.Connection.ID, data.Options.WorkspaceID).Rows()
//	if err != nil {
//		return err
//	}
//	defer cursor.Close()
//	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
//		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
//			Ctx: taskCtx,
//			Params: TapdApiParams{
//				ConnectionId: data.Connection.ID,
//
//				WorkspaceID: data.Options.WorkspaceID,
//			},
//			Table: "tapd_api_%",
//		},
//		InputRowType: reflect.TypeOf(models.TapdIssueStatusHistory{}),
//		Input:        cursor,
//		Convert: func(inputRow interface{}) ([]interface{}, error) {
//			toolL := inputRow.(*models.TapdIssueStatusHistory)
//			domainL := &ticket.IssueStatusHistory{
//				IssueId:        IssueIdGen.Generate(data.Connection.ID, toolL.IssueId),
//				OriginalStatus: toolL.OriginalStatus,
//				StartDate:      toolL.StartDate,
//				EndDate:        &toolL.EndDate,
//			}
//			return []interface{}{
//				domainL,
//			}, nil
//		},
//	})
//	if err != nil {
//		return err
//	}
//
//	return converter.Execute()
//}
//
//var ConvertIssueStatusHistoryMeta = core.SubTaskMeta{
//	Name:             "convertIssueStatusHistory",
//	EntryPoint:       ConvertIssueStatusHistory,
//	EnabledByDefault: true,
//	Description:      "convert Tapd IssueStatusHistory",
//}

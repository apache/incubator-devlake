package tasks

//import (
//	"github.com/merico-dev/lake/models/domainlayer/didgen"
//	"github.com/merico-dev/lake/models/domainlayer/ticket"
//	"github.com/merico-dev/lake/plugins/core"
//	"github.com/merico-dev/lake/plugins/helper"
//	"github.com/merico-dev/lake/plugins/tapd/models"
//	"reflect"
//)
//
//func ConvertIssueSprintsHistory(taskCtx core.SubTaskContext) error {
//	data := taskCtx.GetData().(*TapdTaskData)
//	logger := taskCtx.GetLogger()
//	db := taskCtx.GetDb()
//	logger.Info("convert board:%d", data.Options.WorkspaceID)
//	iterIdGen := didgen.NewDomainIdGenerator(&models.TapdIteration{})
//	cursor, err := db.Model(&models.TapdIssueSprintHistory{}).Where("source_id = ? AND workspace_id = ?", data.Source.ID, data.Options.WorkspaceID).Rows()
//	if err != nil {
//		return err
//	}
//	defer cursor.Close()
//	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
//		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
//			Ctx: taskCtx,
//			Params: TapdApiParams{
//				SourceId: data.Source.ID,
//				//CompanyId:   data.Source.CompanyId,
//				WorkspaceID: data.Options.WorkspaceID,
//			},
//			Table: "tapd_api_%",
//		},
//		InputRowType: reflect.TypeOf(models.TapdIssueSprintHistory{}),
//		Input:        cursor,
//		Convert: func(inputRow interface{}) ([]interface{}, error) {
//			toolL := inputRow.(*models.TapdIssueSprintHistory)
//			domainL := &ticket.IssueSprintsHistory{
//				IssueId:   IssueIdGen.Generate(data.Source.ID, toolL.IssueId),
//				SprintId:  iterIdGen.Generate(data.Source.ID, toolL.SprintId),
//				StartDate: toolL.StartDate,
//				EndDate:   &toolL.EndDate,
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
//var ConvertIssueSprintsHistoryMeta = core.SubTaskMeta{
//	Name:             "convertIssueSprintsHistory",
//	EntryPoint:       ConvertIssueSprintsHistory,
//	EnabledByDefault: true,
//	Description:      "convert Tapd IssueSprintsHistory",
//}

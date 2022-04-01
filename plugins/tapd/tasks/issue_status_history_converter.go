package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"reflect"
)

func ConvertIssueStatusHistory(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("convert board:%d", data.Options.WorkspaceId)
	cursor, err := db.Model(&models.TapdIssueStatusHistory{}).Where("source_id = ? AND workspace_id = ?", data.Source.ID, data.Options.WorkspaceId).Rows()
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
			Table: "tapd_api_%",
		},
		InputRowType: reflect.TypeOf(models.TapdIssueStatusHistory{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			toolL := inputRow.(*models.TapdIssueStatusHistory)
			domainL := &ticket.IssueStatusHistory{
				IssueId:        IssueIdGen.Generate(data.Source.ID, toolL.IssueId),
				OriginalStatus: toolL.OriginalStatus,
				StartDate:      toolL.StartDate,
				EndDate:        &toolL.EndDate,
			}
			return []interface{}{
				domainL,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

var ConvertIssueStatusHistoryMeta = core.SubTaskMeta{
	Name:             "convertIssueStatusHistory",
	EntryPoint:       ConvertIssueStatusHistory,
	EnabledByDefault: true,
	Description:      "convert Tapd IssueStatusHistory",
}

package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"reflect"
)

func ConvertIterationIssue(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("convert board:%d", data.Options.WorkspaceId)
	iterIdGen := didgen.NewDomainIdGenerator(&models.TapdIteration{})

	cursor, err := db.Model(&models.TapdIterationIssue{}).Where("source_id = ? AND workspace_id = ?", data.Source.ID, data.Options.WorkspaceId).Rows()
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
			Table: "_tool_tapd_api_%",
		},
		InputRowType: reflect.TypeOf(models.TapdIterationIssue{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			toolL := inputRow.(*models.TapdIterationIssue)
			domainL := &ticket.SprintIssue{
				SprintId:      iterIdGen.Generate(data.Source.ID, toolL.IterationId),
				IssueId:       IssueIdGen.Generate(data.Source.ID, toolL.IssueId),
				AddedDate:     nil,
				RemovedDate:   nil,
				AddedStage:    nil,
				ResolvedStage: nil,
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

var ConvertIterationIssueMeta = core.SubTaskMeta{
	Name:             "convertIterationIssue",
	EntryPoint:       ConvertIterationIssue,
	EnabledByDefault: true,
	Description:      "convert Tapd IterationIssue",
}

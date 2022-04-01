package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"reflect"
	"strconv"
)

func ConvertTask(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("convert board:%d", data.Options.WorkspaceId)

	cursor, err := db.Model(&models.TapdTask{}).Where("source_id = ? AND workspace_id = ?", data.Source.ID, data.Options.WorkspaceId).Rows()
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
			Table: RAW_TASK_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdTask{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			toolL := inputRow.(*models.TapdTask)
			issue := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: IssueIdGen.Generate(toolL.SourceId, toolL.ID),
				},
				Url:            toolL.Url,
				Key:            strconv.FormatUint(toolL.ID, 10),
				Title:          toolL.Name,
				Summary:        toolL.Description,
				EpicKey:        toolL.EpicKey,
				Type:           "TASK",
				Status:         toolL.Status,
				ResolutionDate: toolL.Completed,
				CreatedDate:    toolL.Created,
				UpdatedDate:    toolL.Modified,
				ParentIssueId:  IssueIdGen.Generate(toolL.SourceId, toolL.StoryID),
				Priority:       toolL.Priority,
				CreatorId:      UserIdGen.Generate(data.Options.SourceId, toolL.WorkspaceId, toolL.Creator),
				AssigneeId:     UserIdGen.Generate(data.Options.SourceId, toolL.WorkspaceId, toolL.Owner),
				AssigneeName:   toolL.Owner,
			}
			if issue.ResolutionDate != nil && issue.CreatedDate != nil {
				issue.LeadTimeMinutes = uint(int64(issue.ResolutionDate.Minute() - issue.CreatedDate.Minute()))
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

var ConvertTaskMeta = core.SubTaskMeta{
	Name:             "convertTask",
	EntryPoint:       ConvertTask,
	EnabledByDefault: true,
	Description:      "convert Tapd Task",
}

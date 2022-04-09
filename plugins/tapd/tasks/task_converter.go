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
			domainL := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: IssueIdGen.Generate(toolL.SourceId, toolL.ID),
				},
				Url:            toolL.Url,
				Number:         strconv.FormatUint(toolL.ID, 10),
				Title:          toolL.Name,
				Description:    toolL.Description,
				EpicKey:        toolL.EpicKey,
				Type:           toolL.StdType,
				Status:         toolL.StdStatus,
				ResolutionDate: toolL.Completed,
				CreatedDate:    toolL.Created,
				UpdatedDate:    toolL.Modified,
				ParentIssueId:  IssueIdGen.Generate(toolL.SourceId, toolL.StoryID),
				Priority:       toolL.Priority,
				CreatorId:      UserIdGen.Generate(data.Options.SourceId, toolL.WorkspaceId, toolL.Creator),
				AssigneeId:     UserIdGen.Generate(data.Options.SourceId, toolL.WorkspaceId, toolL.Owner),
				AssigneeName:   toolL.Owner,
			}
			if domainL.ResolutionDate != nil && domainL.CreatedDate != nil {
				domainL.LeadTimeMinutes = uint(int64(domainL.ResolutionDate.Minute() - domainL.CreatedDate.Minute()))
			}
			results := make([]interface{}, 0, 2)
			boardIssue := &ticket.BoardIssue{
				BoardId: WorkspaceIdGen.Generate(data.Options.WorkspaceId),
				IssueId: domainL.Id,
			}
			results = append(results, domainL, boardIssue)
			return results, nil
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

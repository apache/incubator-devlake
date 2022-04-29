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
	logger.Info("convert board:%d", data.Options.WorkspaceID)

	cursor, err := db.Model(&models.TapdTask{}).Where("connection_id = ? AND workspace_id = ?", data.Connection.ID, data.Options.WorkspaceID).Rows()
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
			Table: RAW_TASK_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdTask{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			toolL := inputRow.(*models.TapdTask)
			domainL := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: IssueIdGen.Generate(toolL.ConnectionId, toolL.ID),
				},
				Url:            toolL.Url,
				Number:         strconv.FormatUint(toolL.ID, 10),
				Title:          toolL.Name,
				Description:    toolL.Description,
				Type:           toolL.StdType,
				Status:         toolL.StdStatus,
				OriginalStatus: toolL.Status,
				ResolutionDate: toolL.Completed.ToNullableTime(),
				CreatedDate:    toolL.Created.ToNullableTime(),
				UpdatedDate:    toolL.Modified.ToNullableTime(),
				ParentIssueId:  IssueIdGen.Generate(toolL.ConnectionId, toolL.StoryID),
				Priority:       toolL.Priority,
				CreatorId:      UserIdGen.Generate(data.Connection.ID, toolL.WorkspaceID, toolL.Creator),
				AssigneeId:     UserIdGen.Generate(data.Connection.ID, toolL.WorkspaceID, toolL.Owner),
				AssigneeName:   toolL.Owner,
			}
			if domainL.ResolutionDate != nil && domainL.CreatedDate != nil {
				domainL.LeadTimeMinutes = uint(int64(domainL.ResolutionDate.Minute() - domainL.CreatedDate.Minute()))
			}
			results := make([]interface{}, 0, 2)
			boardIssue := &ticket.BoardIssue{
				BoardId: WorkspaceIdGen.Generate(toolL.WorkspaceID),
				IssueId: domainL.Id,
			}
			sprintIssue := &ticket.SprintIssue{
				SprintId: IterIdGen.Generate(data.Connection.ID, toolL.IterationID),
				IssueId:  domainL.Id,
			}
			results = append(results, domainL, boardIssue, sprintIssue)
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

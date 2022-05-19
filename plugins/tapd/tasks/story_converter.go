package tasks

import (
	"github.com/apache/incubator-devlake/models/domainlayer"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/helper"
	"github.com/apache/incubator-devlake/plugins/tapd/models"
	"reflect"
	"strconv"
	"time"
)

func ConvertStory(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("convert board:%d", data.Options.WorkspaceID)

	cursor, err := db.Model(&models.TapdStory{}).Where("connection_id = ? AND workspace_id = ?", data.Connection.ID, data.Options.WorkspaceID).Rows()
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
			Table: RAW_STORY_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdStory{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			toolL := inputRow.(*models.TapdStory)
			domainL := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: IssueIdGen.Generate(toolL.ConnectionId, toolL.ID),
				},
				Url:                  toolL.Url,
				Number:               strconv.FormatUint(toolL.ID, 10),
				Title:                toolL.Name,
				Type:                 toolL.StdType,
				Status:               toolL.StdStatus,
				StoryPoint:           uint(toolL.Size),
				OriginalStatus:       toolL.Status,
				ResolutionDate:       (*time.Time)(toolL.Completed),
				CreatedDate:          (*time.Time)(toolL.Created),
				UpdatedDate:          (*time.Time)(toolL.Modified),
				ParentIssueId:        IssueIdGen.Generate(toolL.ConnectionId, toolL.ParentID),
				Priority:             toolL.Priority,
				TimeRemainingMinutes: int64(toolL.Remain),
				CreatorId:            UserIdGen.Generate(data.Connection.ID, toolL.WorkspaceID, toolL.Creator),
				AssigneeId:           UserIdGen.Generate(data.Connection.ID, toolL.WorkspaceID, toolL.Owner),
				AssigneeName:         toolL.Owner,
				Severity:             "",
				Component:            toolL.Feature,
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

var ConvertStoryMeta = core.SubTaskMeta{
	Name:             "convertStory",
	EntryPoint:       ConvertStory,
	EnabledByDefault: true,
	Description:      "convert Tapd story",
}

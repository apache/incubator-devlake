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

func ConvertStory(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("convert board:%d", data.Options.WorkspaceId)

	cursor, err := db.Model(&models.TapdStory{}).Where("source_id = ? AND workspace_id = ?", data.Source.ID, data.Options.WorkspaceId).Rows()
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
			Table: RAW_STORY_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdStory{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			toolL := inputRow.(*models.TapdStory)
			issue := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: IssueIdGen.Generate(toolL.SourceId, toolL.ID),
				},
				Url:                  toolL.Url,
				Key:                  strconv.FormatUint(toolL.ID, 10),
				Title:                toolL.Name,
				Summary:              toolL.Name,
				EpicKey:              toolL.EpicKey,
				Type:                 toolL.Type,
				Status:               toolL.Status,
				StoryPoint:           uint(toolL.Size),
				ResolutionDate:       toolL.Completed,
				CreatedDate:          toolL.Created,
				UpdatedDate:          toolL.Modified,
				ParentIssueId:        IssueIdGen.Generate(toolL.SourceId, toolL.ParentID),
				Priority:             toolL.Priority,
				TimeRemainingMinutes: int64(toolL.Remain),
				CreatorId:            UserIdGen.Generate(data.Options.SourceId, toolL.WorkspaceId, toolL.Creator),
				AssigneeId:           UserIdGen.Generate(data.Options.SourceId, toolL.WorkspaceId, toolL.Owner),
				AssigneeName:         toolL.Owner,
				Severity:             "",
				Component:            toolL.Feature,
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

var ConvertStoryMeta = core.SubTaskMeta{
	Name:             "convertStory",
	EntryPoint:       ConvertStory,
	EnabledByDefault: true,
	Description:      "convert Tapd story",
}

package tasks

import (
	"fmt"
	"github.com/merico-dev/lake/models/domainlayer"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
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
	issueIdGen := didgen.NewDomainIdGenerator(&models.TapdStory{})
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
			Table: RAW_WORKSPACE_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdStory{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			story := inputRow.(*models.TapdStory)
			issue := &ticket.Issue{
				DomainEntity: domainlayer.DomainEntity{
					Id: issueIdGen.Generate(story.SourceId, story.ID),
				},
				Url:                  fmt.Sprintf("https://www.tapd.cn/%d/prong/stories/view/%d", story.WorkspaceID, story.ID),
				Key:                  strconv.FormatUint(story.ID, 10),
				Title:                story.Name,
				Summary:              story.Name,
				EpicKey:              story.EpicKey,
				Type:                 story.Type,
				Status:               story.Status,
				StoryPoint:           uint(story.Size),
				ResolutionDate:       story.Completed,
				CreatedDate:          story.Created,
				UpdatedDate:          story.Modified,
				LeadTimeMinutes:      uint(story.EffortCompleted),
				ParentIssueId:        issueIdGen.Generate(story.SourceId, story.ParentID),
				Priority:             story.Priority,
				TimeRemainingMinutes: int64(story.Remain),
				CreatorId:            UserIdGen.Generate(data.Options.SourceId, story.Creator),
				AssigneeId:           UserIdGen.Generate(data.Options.SourceId, story.Owner),
				AssigneeName:         story.Owner,
				Severity:             "",
				Component:            story.Feature,
			}
			if issue.ResolutionDate != nil && issue.CreatedDate != nil {
				issue.TimeSpentMinutes = int64(issue.ResolutionDate.Minute() - issue.CreatedDate.Minute())
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

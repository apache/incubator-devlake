package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer/crossdomain"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
	"github.com/merico-dev/lake/plugins/tapd/models"
	"reflect"
)

func ConvertStoryCommit(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*TapdTaskData)
	logger := taskCtx.GetLogger()
	db := taskCtx.GetDb()
	logger.Info("convert board:%d", data.Options.WorkspaceID)

	cursor, err := db.Model(&models.TapdStoryCommit{}).Where("source_id = ? AND workspace_id = ?", data.Source.ID, data.Options.WorkspaceID).Rows()
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
				WorkspaceID: data.Options.WorkspaceID,
			},
			Table: RAW_STORY_COMMIT_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdStoryCommit{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			toolL := inputRow.(*models.TapdStoryCommit)
			domainL := &crossdomain.IssueCommit{
				IssueId:   IssueIdGen.Generate(models.Uint64s(data.Source.ID), models.Uint64s(toolL.StoryId)),
				CommitSha: toolL.CommitID,
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

var ConvertStoryCommitMeta = core.SubTaskMeta{
	Name:             "convertStoryCommit",
	EntryPoint:       ConvertStoryCommit,
	EnabledByDefault: true,
	Description:      "convert Tapd StoryCommit",
}

package tasks

import (
	"github.com/merico-dev/lake/plugins/tapd/models"
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

var ConvertStoryLabelsMeta = core.SubTaskMeta{
	Name:             "convertStoryLabels",
	EntryPoint:       ConvertStoryLabels,
	EnabledByDefault: true,
	Description:      "Convert tool layer table tapd_issue_labels into  domain layer table issue_labels",
}

func ConvertStoryLabels(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*TapdTaskData)

	cursor, err := db.Model(&models.TapdStoryLabel{}).
		Joins(`left join _tool_tapd_workspace_issues on _tool_tapd_workspace_issues.issue_id = _tool_tapd_story_labels.story_id`).
		Where("_tool_tapd_workspace_issues.workspace_id = ?", data.Options.WorkspaceID).
		Order("issue_id ASC").
		Rows()
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
			Table: RAW_BUG_TABLE,
		},
		InputRowType: reflect.TypeOf(models.TapdStoryLabel{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			issueLabel := inputRow.(*models.TapdStoryLabel)
			domainStoryLabel := &ticket.IssueLabel{
				IssueId:   IssueIdGen.Generate(issueLabel.StoryId),
				LabelName: issueLabel.LabelName,
			}
			return []interface{}{
				domainStoryLabel,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

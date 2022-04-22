package tasks

import (
	"github.com/merico-dev/lake/plugins/tapd/models"
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

var ConvertIssueLabelsMeta = core.SubTaskMeta{
	Name:             "convertIssueLabels",
	EntryPoint:       ConvertIssueLabels,
	EnabledByDefault: true,
	Description:      "Convert tool layer table tapd_issue_labels into  domain layer table issue_labels",
}

func ConvertIssueLabels(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*TapdTaskData)

	err := db.Model(&models.TapdIssueLabel{}).Where("_raw_data_table like ?", "_raw_tapd_api_%").
		Update("_raw_data_table", "_raw_tapd_api_issues").Error
	if err != nil {
		return err
	}
	cursor, err := db.Model(&models.TapdIssueLabel{}).
		Joins(`left join _tool_tapd_workspace_issues on _tool_tapd_workspace_issues.issue_id = _tool_tapd_issue_labels.issue_id`).
		Where("_tool_tapd_workspace_issues.workspace_id = ?", data.Source.WorkspaceId).
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
				SourceId: data.Source.ID,
				//CompanyId:   data.Source.CompanyId,
				WorkspaceId: data.Options.WorkspaceId,
			},
			Table: "tapd_api_issues",
		},
		InputRowType: reflect.TypeOf(models.TapdIssueLabel{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			issueLabel := inputRow.(*models.TapdIssueLabel)
			domainIssueLabel := &ticket.IssueLabel{
				IssueId:   IssueIdGen.Generate(issueLabel.IssueId),
				LabelName: issueLabel.LabelName,
			}
			return []interface{}{
				domainIssueLabel,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

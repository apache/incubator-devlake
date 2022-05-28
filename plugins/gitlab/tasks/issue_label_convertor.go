package tasks

import (
	"reflect"

	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/models/domainlayer/ticket"
	"github.com/apache/incubator-devlake/plugins/core"
	gitlabModels "github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
)

var ConvertIssueLabelsMeta = core.SubTaskMeta{
	Name:             "convertIssueLabels",
	EntryPoint:       ConvertIssueLabels,
	EnabledByDefault: true,
	Description:      "Convert tool layer table gitlab_issue_labels into  domain layer table issue_labels",
}

func ConvertIssueLabels(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GitlabTaskData)
	projectId := data.Options.ProjectId

	cursor, err := db.Model(&gitlabModels.GitlabIssueLabel{}).
		Joins(`left join _tool_gitlab_issues on _tool_gitlab_issues.gitlab_id = _tool_gitlab_issue_labels.issue_id`).
		Where("_tool_gitlab_issues.project_id = ?", projectId).
		Order("issue_id ASC").
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	issueIdGen := didgen.NewDomainIdGenerator(&gitlabModels.GitlabIssue{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GitlabApiParams{
				ProjectId: projectId,
			},
			Table: RAW_ISSUE_TABLE,
		},
		InputRowType: reflect.TypeOf(gitlabModels.GitlabIssueLabel{}),
		Input:        cursor,
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			issueLabel := inputRow.(*gitlabModels.GitlabIssueLabel)
			domainIssueLabel := &ticket.IssueLabel{
				IssueId:   issueIdGen.Generate(issueLabel.IssueId),
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

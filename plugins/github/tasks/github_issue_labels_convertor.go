package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/models/domainlayer/ticket"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
	"reflect"
)

func ConvertIssueLabels(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)

	cursor, err := db.Model(&githubModels.GithubIssueLabel{}).
		Where("_raw_data_params = ?", data.Options.ParamString).
		Order("issue_id ASC").
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	issueIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubIssue{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		Ctx:          taskCtx,
		InputRowType: reflect.TypeOf(githubModels.GithubIssueLabel{}),
		Input:        cursor,
		BatchSelectors: map[reflect.Type]helper.BatchSelector{
			reflect.TypeOf(&ticket.IssueLabel{}): {
				Query: "_raw_data_params = ?",
				Parameters: []interface{}{
					data.Options.ParamString,
				},
			},
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			issueLabel := inputRow.(*githubModels.GithubIssueLabel)
			domainIssueLabel := &ticket.IssueLabel{
				IssueId:   issueIdGen.Generate(issueLabel.IssueId),
				LabelName: issueLabel.LabelName,
			}
			domainIssueLabel.RawDataOrigin = issueLabel.RawDataOrigin
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

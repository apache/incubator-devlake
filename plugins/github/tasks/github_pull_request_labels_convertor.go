package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
	"reflect"
)

func ConvertPullRequestLabels(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)

	cursor, err := db.Model(&githubModels.GithubPullRequestLabel{}).
		Where("github_pull_request_labels._raw_data_params = ?", data.Options.ParamString).
		Order("pull_id ASC").
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	prIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubPullRequest{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		Ctx:          taskCtx,
		InputRowType: reflect.TypeOf(githubModels.GithubPullRequestLabel{}),
		Input:        cursor,
		BatchSelectors: map[reflect.Type]helper.BatchSelector{
			reflect.TypeOf(&code.PullRequestLabel{}): {
				Query: "_raw_data_params = ?",
				Parameters: []interface{}{
					data.Options.ParamString,
				},
			},
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			prLabel := inputRow.(*githubModels.GithubPullRequestLabel)
			domainPrLabel := &code.PullRequestLabel{
				PullRequestId: prIdGen.Generate(prLabel.PullId),
				LabelName:     prLabel.LabelName,
			}
			domainPrLabel.RawDataOrigin = prLabel.RawDataOrigin

			return []interface{}{
				domainPrLabel,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

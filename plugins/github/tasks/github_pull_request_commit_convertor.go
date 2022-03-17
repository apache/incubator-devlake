package tasks

import (
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
	"reflect"
)

func ConvertPullRequestCommits(taskCtx core.SubTaskContext) (err error) {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)

	pullIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubPullRequest{})

	cursor, err := db.Model(&githubModels.GithubPullRequestCommit{}).
		Where("_raw_data_params = ?", data.Options.ParamString).
		Order("pull_request_id ASC").Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		Ctx:          taskCtx,
		InputRowType: reflect.TypeOf(githubModels.GithubPullRequestCommit{}),
		Input:        cursor,
		BatchSelectors: map[reflect.Type]helper.BatchSelector{
			reflect.TypeOf(&code.PullRequestCommit{}): {
				Query: "_raw_data_params = ?",
				Parameters: []interface{}{
					data.Options.ParamString,
				},
			},
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			githubPullRequestCommit := inputRow.(*githubModels.GithubPullRequestCommit)
			domainPrCommit := &code.PullRequestCommit{
				CommitSha:     githubPullRequestCommit.CommitSha,
				PullRequestId: pullIdGen.Generate(githubPullRequestCommit.PullRequestId),
			}
			domainPrCommit.RawDataOrigin = githubPullRequestCommit.RawDataOrigin
			return []interface{}{
				domainPrCommit,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

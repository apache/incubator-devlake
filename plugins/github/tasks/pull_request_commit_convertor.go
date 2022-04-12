package tasks

import (
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var ConvertPullRequestCommitsMeta = core.SubTaskMeta{
	Name:             "convertPullRequestCommits",
	EntryPoint:       ConvertPullRequestCommits,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_pull_request_commits into  domain layer table pull_request_commits",
}

func ConvertPullRequestCommits(taskCtx core.SubTaskContext) (err error) {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	pullIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubPullRequest{})

	cursor, err := db.Model(&githubModels.GithubPullRequestCommit{}).
		Joins(`left join _tool_github_pull_requests on _tool_github_pull_requests.github_id = _tool_github_pull_request_commits.pull_request_id`).
		Where("_tool_github_pull_requests.repo_id = ?", repoId).
		Order("pull_request_id ASC").Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(githubModels.GithubPullRequestCommit{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			Table: RAW_PULL_REQUEST_COMMIT_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			githubPullRequestCommit := inputRow.(*githubModels.GithubPullRequestCommit)
			domainPrCommit := &code.PullRequestCommit{
				CommitSha:     githubPullRequestCommit.CommitSha,
				PullRequestId: pullIdGen.Generate(githubPullRequestCommit.PullRequestId),
			}
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

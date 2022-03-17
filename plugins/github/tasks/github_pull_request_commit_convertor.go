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
	repoId := data.Repo.GithubId

	pullIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubPullRequest{})

	cursor, err := db.Model(&githubModels.GithubPullRequestCommit{}).
		Joins(`left join github_pull_requests on github_pull_requests.github_id = github_pull_request_commits.pull_request_id`).
		Where("github_pull_requests.repo_id = ?", repoId).
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
				Query: "pull_request_id like ?",
				Parameters: []interface{}{
					pullIdGen.Generate(data.Repo.GithubId, didgen.WILDCARD),
				},
			},
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

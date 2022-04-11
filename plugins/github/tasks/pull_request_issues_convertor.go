package tasks

import (
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer/crossdomain"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var ConvertPullRequestIssuesMeta = core.SubTaskMeta{
	Name:             "convertPullRequestIssues",
	EntryPoint:       ConvertPullRequestIssues,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_pull_request_issues into  domain layer table pull_request_issues",
}

func ConvertPullRequestIssues(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	cursor, err := db.Model(&githubModels.GithubPullRequestIssue{}).
		Joins(`left join _tool_github_pull_requests on _tool_github_pull_requests.github_id = _tool_github_pull_request_issues.pull_request_id`).
		Where("_tool_github_pull_requests.repo_id = ?", repoId).
		Order("pull_request_id ASC").
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	prIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubPullRequest{})
	issueIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubIssue{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(githubModels.GithubPullRequestIssue{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx: taskCtx,
			Params: GithubApiParams{
				Owner: data.Options.Owner,
				Repo:  data.Options.Repo,
			},
			Table: RAW_PULL_REQUEST_TABLE,
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			githubPrIssue := inputRow.(*githubModels.GithubPullRequestIssue)
			pullRequestIssue := &crossdomain.PullRequestIssue{
				PullRequestId: prIdGen.Generate(githubPrIssue.PullRequestId),
				IssueId:       issueIdGen.Generate(githubPrIssue.IssueId),
				IssueNumber:   githubPrIssue.IssueNumber,
				PullNumber:    githubPrIssue.PullNumber,
			}
			return []interface{}{
				pullRequestIssue,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

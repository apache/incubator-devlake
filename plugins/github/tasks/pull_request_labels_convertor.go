package tasks

import (
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	githubModels "github.com/merico-dev/lake/plugins/github/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var ConvertPullRequestLabelsMeta = core.SubTaskMeta{
	Name:             "convertPullRequestLabels",
	EntryPoint:       ConvertPullRequestLabels,
	EnabledByDefault: true,
	Description:      "Convert tool layer table github_pull_request_labels into  domain layer table pull_request_labels",
}

func ConvertPullRequestLabels(taskCtx core.SubTaskContext) error {
	db := taskCtx.GetDb()
	data := taskCtx.GetData().(*GithubTaskData)
	repoId := data.Repo.GithubId

	cursor, err := db.Model(&githubModels.GithubPullRequestLabel{}).
		Joins(`left join github_pull_requests on github_pull_requests.github_id = github_pull_request_labels.pull_id`).
		Where("github_pull_requests.repo_id = ?", repoId).
		Order("pull_id ASC").
		Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	prIdGen := didgen.NewDomainIdGenerator(&githubModels.GithubPullRequest{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(githubModels.GithubPullRequestLabel{}),
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
			prLabel := inputRow.(*githubModels.GithubPullRequestLabel)
			domainPrLabel := &code.PullRequestLabel{
				PullRequestId: prIdGen.Generate(prLabel.PullId),
				LabelName:     prLabel.LabelName,
			}
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

package tasks

import (
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/didgen"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/gitlab/models"
	"github.com/merico-dev/lake/plugins/helper"
)

var ConvertApiMergeRequestsCommitsMeta = core.SubTaskMeta{
	Name:             "convertApiMergeRequestsCommits",
	EntryPoint:       ConvertApiMergeRequestsCommits,
	EnabledByDefault: true,
	Description:      "Update domain layer PullRequestCommit according to GitlabMergeRequestCommit",
}

func ConvertApiMergeRequestsCommits(taskCtx core.SubTaskContext) error {
	rawDataSubTaskArgs, data := CreateRawDataSubTaskArgs(taskCtx, RAW_MERGE_REQUEST_COMMITS_TABLE)
	db := taskCtx.GetDb()

	cursor, err := db.Model(&models.GitlabMergeRequestCommit{}).
		Joins(`left join gitlab_merge_requests on gitlab_merge_requests.gitlab_id = gitlab_merge_request_commits.merge_request_id`).
		Where("gitlab_merge_requests.project_id = ?", data.Options.ProjectId).
		Order("merge_request_id ASC").Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	// TODO: adopt batch indate operation
	domainIdGenerator := didgen.NewDomainIdGenerator(&models.GitlabMergeRequest{})

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		RawDataSubTaskArgs: *rawDataSubTaskArgs,
		InputRowType:       reflect.TypeOf(models.GitlabMergeRequestCommit{}),
		Input:              cursor,

		Convert: func(inputRow interface{}) ([]interface{}, error) {
			gitlabMergeRequestCommit := inputRow.(*models.GitlabMergeRequestCommit)
			domainPrcommit := &code.PullRequestCommit{
				CommitSha:     gitlabMergeRequestCommit.CommitSha,
				PullRequestId: domainIdGenerator.Generate(gitlabMergeRequestCommit.MergeRequestId),
			}
			return []interface{}{
				domainPrcommit,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

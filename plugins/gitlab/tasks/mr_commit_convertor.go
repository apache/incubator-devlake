package tasks

import (
	"reflect"

	"github.com/apache/incubator-devlake/models/domainlayer/code"
	"github.com/apache/incubator-devlake/models/domainlayer/didgen"
	"github.com/apache/incubator-devlake/plugins/core"
	"github.com/apache/incubator-devlake/plugins/gitlab/models"
	"github.com/apache/incubator-devlake/plugins/helper"
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
		Joins(`left join _tool_gitlab_merge_requests on _tool_gitlab_merge_requests.gitlab_id = _tool_gitlab_merge_request_commits.merge_request_id`).
		Where("_tool_gitlab_merge_requests.project_id = ?", data.Options.ProjectId).
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

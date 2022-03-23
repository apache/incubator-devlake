package tasks

import (
	"fmt"
	"reflect"

	"github.com/merico-dev/lake/models/domainlayer/crossdomain"
	"github.com/merico-dev/lake/plugins/core"
	"github.com/merico-dev/lake/plugins/helper"
)

func CalculateIssuesDiff(taskCtx core.SubTaskContext) error {
	data := taskCtx.GetData().(*RefdiffTaskData)
	repoId := data.Options.RepoId
	pairs := data.Options.Pairs
	db := taskCtx.GetDb()
	// use to calculate progress
	pairList := make([][2]string, len(pairs))
	for i, pair := range pairs {
		pairList[i] = [2]string{fmt.Sprintf("%s:%s", repoId, pair.NewRef), fmt.Sprintf("%s:%s", repoId, pair.OldRef)}
	}
	cursor, err := db.Table("refs_commits_diffs").
		Joins("left join pull_requests on pull_requests.merge_commit_sha = refs_commits_diffs.commit_sha").
		Joins("left join pull_request_issues on pull_request_issues.pull_request_id = pull_requests.id").
		Joins("left join refs on refs.commit_sha = refs_commits_diffs.new_ref_commit_sha").
		Order("refs_commits_diffs.new_ref_name ASC").
		Where("refs.repo_id = ? and pull_request_issues.issue_number > 0 and (refs_commits_diffs.new_ref_name, refs_commits_diffs.old_ref_name) in ?",
			repoId, pairList).
		Select("refs_commits_diffs.new_ref_commit_sha as new_ref_commit_sha, refs_commits_diffs.old_ref_commit_sha as old_ref_commit_sha, " +
			"pull_request_issues.issue_id as issue_id, pull_request_issues.issue_number as issue_number, " +
			"refs_commits_diffs.new_ref_name as new_ref_name, refs_commits_diffs.old_ref_name as old_ref_name").Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()

	converter, err := helper.NewDataConverter(helper.DataConverterArgs{
		InputRowType: reflect.TypeOf(crossdomain.RefsIssuesDiffs{}),
		Input:        cursor,
		RawDataSubTaskArgs: helper.RawDataSubTaskArgs{
			Ctx:   taskCtx,
			Table: "refs_commits_diffs",
		},
		Convert: func(inputRow interface{}) ([]interface{}, error) {
			refPairIssue := inputRow.(*crossdomain.RefsIssuesDiffs)
			return []interface{}{
				refPairIssue,
			}, nil
		},
	})
	if err != nil {
		return err
	}

	return converter.Execute()
}

var CalculateIssuesDiffMeta = core.SubTaskMeta{
	Name:             "calculateIssuesDiff",
	EntryPoint:       CalculateIssuesDiff,
	EnabledByDefault: true,
	Description:      "Calculate diff issues between refs",
}

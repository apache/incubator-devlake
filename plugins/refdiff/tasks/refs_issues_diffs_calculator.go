package tasks

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/errors"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/crossdomain"
	"gorm.io/gorm/clause"
)

func CalculateIssuesDiff(ctx context.Context, pairs []RefPair, progress chan<- float32, repoId string) error {
	// use to calculate progress
	var numOfPairs int64
	index := int64(1)
	pairList := make([][2]string, len(pairs))
	for i, pair := range pairs {
		pairList[i] = [2]string{fmt.Sprintf("%s:%s", repoId, pair.NewRef), fmt.Sprintf("%s:%s", repoId, pair.OldRef)}
	}
	db := lakeModels.Db.Table("refs_commits_diffs").
		Joins("left join (  "+
			"select pull_request_id as id, commit_sha from pull_request_commits "+
			"left join pull_requests p on pull_request_commits.pull_request_id = p.id "+
			"where p.repo_id = '"+repoId+
			"' union  "+
			"select id, merge_commit_sha as commit_sha from pull_requests where repo_id = '"+repoId+"') _combine_pr "+
			"on _combine_pr.commit_sha = refs_commits_diffs.commit_sha").
		Joins("left join pull_request_issues on pull_request_issues.pull_request_id = _combine_pr.id").
		Joins("left join refs on refs.commit_sha = refs_commits_diffs.new_ref_commit_sha").
		Order("refs_commits_diffs.new_ref_name ASC").
		Where("refs.repo_id = ? and pull_request_issues.issue_number > 0 and (refs_commits_diffs.new_ref_name, refs_commits_diffs.old_ref_name) in ?",
			repoId, pairList).
		Select("refs_commits_diffs.new_ref_commit_sha as new_ref_commit_sha, refs_commits_diffs.old_ref_commit_sha as old_ref_commit_sha, " +
			"pull_request_issues.issue_id as issue_id, pull_request_issues.issue_number as issue_number, " +
			"refs_commits_diffs.new_ref_name as new_ref_name, refs_commits_diffs.old_ref_name as old_ref_name")
	err := db.Count(&numOfPairs).Error
	if err != nil {
		return err
	}
	// we iterate the whole refCommitsDiff table, and convert the commit sha to issue
	refPairIssue := &crossdomain.RefsIssuesDiffs{}
	cursor, err := db.Rows()
	if err != nil {
		return err
	}

	defer cursor.Close()
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.TaskCanceled
		default:
		}

		err = lakeModels.Db.ScanRows(cursor, refPairIssue)
		if err != nil {
			return err
		}
		err = lakeModels.Db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(refPairIssue).Error
		if err != nil {
			return err
		}
		progress <- float32(index) / float32(numOfPairs)
		index++
	}

	return nil
}

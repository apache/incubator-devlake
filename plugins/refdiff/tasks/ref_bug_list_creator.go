package tasks

import (
	"context"
	"github.com/merico-dev/lake/errors"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/crossdomain"
	"gorm.io/gorm/clause"
)

func CreatRefBugStats(ctx context.Context, progress chan<- float32, repoId string) error {
	// use to calculate progress
	var numOfPairs int64
	index := int64(1)
	db := lakeModels.Db.Table("refs_commits_diffs").
		Joins("left join pull_requests on pull_requests.merge_commit_sha = refs_commits_diffs.commit_sha").
		Joins("left join pull_request_issues on pull_request_issues.pull_request_id = pull_requests.id").
		Joins("left join refs on refs.commit_sha = refs_commits_diffs.new_ref_commit_sha").
		Order("refs_commits_diffs.new_ref_name ASC").
		Where("refs.repo_id = ? and pull_request_issues.issue_number > 0", repoId).
		Select("refs_commits_diffs.new_ref_commit_sha as new_ref_commit_sha, refs_commits_diffs.old_ref_commit_sha as old_ref_commit_sha, " +
			"pull_request_issues.issue_id as issue_id, pull_request_issues.issue_number as issue_number")
	err := db.Count(&numOfPairs).Error
	if err != nil {
		return err
	}
	// we iterate the whole refCommitsDiff table, and convert the commit sha to issue
	refPairIssue := &crossdomain.RefPairIssue{}
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

package tasks

import (
	"context"
	"fmt"
	"github.com/merico-dev/lake/errors"
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/models/domainlayer/crossdomain"
	"gorm.io/gorm/clause"
)

func CreateRepoBugList(ctx context.Context, repoId int) (err error) {

	refCommitsDiff := &code.RefsCommitsDiff{}
	// we iterate the whole refCommitsDiff table, and convert the commit sha to issue
	cursor, err := lakeModels.Db.Model(&refCommitsDiff).
		Order("new_ref_name ASC").
		Rows()
	if err != nil {
		return err
	}
	bugListStr := ""
	newRefName := refCommitsDiff.NewRefName
	oldRefName := refCommitsDiff.OldRefName
	lastBugNumber := ""
	count := 0
	defer cursor.Close()
	// iterate all rows
	for cursor.Next() {
		select {
		case <-ctx.Done():
			return errors.TaskCanceled
		default:
		}
		err = lakeModels.Db.ScanRows(cursor, refCommitsDiff)
		if err != nil {
			return err
		}

		if newRefName == "" {
			newRefName = refCommitsDiff.NewRefName
			oldRefName = refCommitsDiff.OldRefName

		}
		//if we are going to convert a new pair, save the last pair
		if newRefName != refCommitsDiff.NewRefName {
			bugList := &crossdomain.RefBugStats{
				NewRefName:  newRefName,
				OldRefName:  oldRefName,
				IssueNumber: bugListStr,
				IssueCount:  count,
			}
			err = lakeModels.Db.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(bugList).Error
			if err != nil {
				return err
			}
			newRefName = refCommitsDiff.NewRefName
			oldRefName = refCommitsDiff.OldRefName
			bugListStr = ""
			count = 0
		}

		var bugNumberList []string
		err = lakeModels.Db.Table("github_pull_requests").
			Joins("left join github_pull_request_issues on github_pull_requests.github_id = github_pull_request_issues.pull_request_id").
			Where("`github_pull_request_issues`.issue_number is not null and github_pull_requests.merge_commit_sha = ?", refCommitsDiff.CommitSha).
			Pluck("`github_pull_request_issues`.issue_number", &bugNumberList).Error
		if err != nil {
			return err
		}

		for _, bugNumber := range bugNumberList {
			//if this is the first issueNumber
			if bugListStr == "" {
				count++
				bugListStr = bugNumber
				continue
			}
			//same number, just continue
			if lastBugNumber == bugNumber {
				continue
			}
			count++
			lastBugNumber = bugNumber
			bugListStr = fmt.Sprintf("%s, %s", bugListStr, bugNumber)
		}

	}
	bugList := &crossdomain.RefBugStats{
		NewRefName:  newRefName,
		OldRefName:  oldRefName,
		IssueNumber: bugListStr,
	}
	err = lakeModels.Db.Save(bugList).Error
	if err != nil {
		return err
	}

	return nil
}

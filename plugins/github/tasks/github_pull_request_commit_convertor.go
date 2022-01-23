package tasks

import (
	lakeModels "github.com/merico-dev/lake/models"
	"github.com/merico-dev/lake/models/domainlayer/code"
	"github.com/merico-dev/lake/plugins/github/models"
	"gorm.io/gorm/clause"
)

func PrCommitConvertor() (err error) {
	githubPullRequestCommit := &models.GithubPullRequestCommit{}

	cursor, err := lakeModels.Db.Model(&githubPullRequestCommit).Rows()
	if err != nil {
		return err
	}
	defer cursor.Close()
	// iterate all rows
	for cursor.Next() {
		err = lakeModels.Db.ScanRows(cursor, githubPullRequestCommit)
		if err != nil {
			return err
		}
		err = lakeModels.Db.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&code.PullRequestCommit{
			CommitSha:     githubPullRequestCommit.CommitSha,
			PullRequestId: githubPullRequestCommit.PullRequestId,
		}).Error
		if err != nil {
			return err
		}
	}
	return nil
}

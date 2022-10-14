package migrationscripts

import (
	"context"
	"time"

	"github.com/apache/incubator-devlake/errors"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type Task20221014 struct {
	ConnectionId        uint64 `gorm:"primaryKey"`
	BitbucketId         string `gorm:"primaryKey"`
	Status              string `gorm:"type:varchar(100)"`
	Result              string `gorm:"type:varchar(100)"`
	RefName             string `gorm:"type:varchar(255)"`
	RepoId              string `gorm:"type:varchar(255)"`
	CommitSha           string `gorm:"type:varchar(255)"`
	WebUrl              string `gorm:"type:varchar(255)"`
	DurationInSeconds   uint64
	BitbucketCreatedOn  *time.Time
	BitbucketCompleteOn *time.Time
	archived.NoPKModel
}

func (Task20221014) TableName() string {
	return "_tool_bitbucket_pipelines"
}

type addRepoIdAndCommitShaField struct{}

func (*addRepoIdAndCommitShaField) Up(ctx context.Context, db *gorm.DB) errors.Error {
	err := db.Migrator().AddColumn(Task20221014{}, "repo_id")
	if err != nil {
		return errors.Convert(err)
	}
	err = db.Migrator().AddColumn(Task20221014{}, "commit_sha")
	if err != nil {
		return errors.Convert(err)
	}
	return nil
}

func (*addRepoIdAndCommitShaField) Version() uint64 {
	return 20221014114623
}

func (*addRepoIdAndCommitShaField) Name() string {
	return "add column `repo_id` and `commit_sha` at _tool_bitbucket_pipelines"
}

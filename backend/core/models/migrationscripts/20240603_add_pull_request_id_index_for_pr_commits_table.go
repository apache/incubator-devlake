package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type addPullRequestIdIndexToPullRequestCommits struct{}

type pullRequestCommits20240602 struct {
	PullRequestId string `gorm:"index"`
}

func (pullRequestCommits20240602) TableName() string {
	return "pull_request_commits"
}

func (u *addPullRequestIdIndexToPullRequestCommits) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &pullRequestCommits20240602{})
}

func (*addPullRequestIdIndexToPullRequestCommits) Version() uint64 {
	return 20240602103400
}

func (*addPullRequestIdIndexToPullRequestCommits) Name() string {
	return "add pull_request_id index for pull_request_commits"
}

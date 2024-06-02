package migrationscripts

import (
	"github.com/apache/incubator-devlake/core/context"
	"github.com/apache/incubator-devlake/core/errors"
	"github.com/apache/incubator-devlake/helpers/migrationhelper"
)

type addPullRequestIdIndexToPullRequestComments struct{}

type pullRequestComments20240602 struct {
	PullRequestId string `gorm:"index"`
}

func (pullRequestComments20240602) TableName() string {
	return "pull_request_comments"
}

func (u *addPullRequestIdIndexToPullRequestComments) Up(basicRes context.BasicRes) errors.Error {
	return migrationhelper.AutoMigrateTables(basicRes, &pullRequestComments20240602{})
}

func (*addPullRequestIdIndexToPullRequestComments) Version() uint64 {
	return 20240602103401
}

func (*addPullRequestIdIndexToPullRequestComments) Name() string {
	return "add pull_request_id index for pull_request_comments"
}

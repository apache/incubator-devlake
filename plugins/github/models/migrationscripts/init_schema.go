package migrationscripts

import (
	"context"

	"github.com/apache/incubator-devlake/plugins/github/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type InitSchemas struct{}

func (*InitSchemas) Up(ctx context.Context, db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		&archived.GithubRepo{},
		&archived.GithubCommit{},
		&archived.GithubRepoCommit{},
		&archived.GithubPullRequest{},
		&archived.GithubReviewer{},
		&archived.GithubPullRequestComment{},
		&archived.GithubPullRequestCommit{},
		&archived.GithubPullRequestLabel{},
		&archived.GithubIssue{},
		&archived.GithubIssueComment{},
		&archived.GithubIssueEvent{},
		&archived.GithubIssueLabel{},
		&archived.GithubUser{},
		&archived.GithubPullRequestIssue{},
		&archived.GithubCommitStat{})
}

func (*InitSchemas) Version() uint64 {
	return 20220407201135
}

func (*InitSchemas) Name() string {
	return "Github init schemas"
}

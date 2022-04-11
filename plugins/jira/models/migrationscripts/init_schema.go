package migrationscripts

import (
	"context"

	"github.com/merico-dev/lake/plugins/jira/models/migrationscripts/archived"
	"gorm.io/gorm"
)

const (
	Owner = "Jira"
)

type InitSchemas struct{}

func (*InitSchemas) Up(ctx context.Context, db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		&archived.JiraProject{},
		&archived.JiraUser{},
		&archived.JiraIssue{},
		&archived.JiraBoard{},
		&archived.JiraBoardIssue{},
		&archived.JiraChangelog{},
		&archived.JiraChangelogItem{},
		&archived.JiraRemotelink{},
		&archived.JiraIssueCommit{},
		&archived.JiraSource{},
		&archived.JiraIssueTypeMapping{},
		&archived.JiraIssueStatusMapping{},
		&archived.JiraSprint{},
		&archived.JiraBoardSprint{},
		&archived.JiraSprintIssue{},
		&archived.JiraWorklog{},
	)
}

func (*InitSchemas) Version() uint64 {
	return 20220407201138
}

func (*InitSchemas) Owner() string {
	return Owner
}

func (*InitSchemas) Name() string {
	return "create init schemas"
}

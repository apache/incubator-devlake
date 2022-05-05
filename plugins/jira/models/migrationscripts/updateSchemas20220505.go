package migrationscripts

import (
	"context"
	"github.com/merico-dev/lake/plugins/jira/models"
	"github.com/merico-dev/lake/plugins/jira/models/migrationscripts/archived"

	"gorm.io/gorm"
)

type UpdateSchemas20220505 struct{}

func (*UpdateSchemas20220505) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().RenameTable(archived.JiraSource{}, models.JiraConnection{})
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraBoard{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraProject{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraUser{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraIssue{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraBoardIssue{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraChangelog{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraChangelogItem{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraRemotelink{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraIssueCommit{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraIssueTypeMapping{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraIssueStatusMapping{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraSprint{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraBoardSprint{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraSprintIssue{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	err = db.Migrator().RenameColumn(archived.JiraWorklog{}, "source_id", "connection_id")
	if err != nil {
		return err
	}
	return nil
}

func (*UpdateSchemas20220505) Version() uint64 {
	return 20220505212344
}

func (*UpdateSchemas20220505) Owner() string {
	return "Jira"
}

func (*UpdateSchemas20220505) Name() string {
	return "Rename source to connection "
}
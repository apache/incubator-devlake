package migrationscripts

import (
	"context"
	"github.com/merico-dev/lake/plugins/tapd/models/migrationscripts/archived"

	"gorm.io/gorm"
)

type InitSchemas struct{}

func (*InitSchemas) Up(ctx context.Context, db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		&archived.TapdWorkspace{},
		&archived.TapdWorkSpaceIssue{},
		&archived.TapdWorklog{},
		&archived.TapdWorkspaceIteration{},
		&archived.TapdUser{},
		&archived.TapdChangelog{},
		&archived.TapdChangelogItem{},
		&archived.TapdIssue{},
		&archived.TapdIssueCommit{},
		&archived.TapdIssueTypeMapping{},
		&archived.TapdIssueStatusMapping{},
		&archived.TapdIssueSprintHistory{},
		&archived.TapdIssueStatusHistory{},
		&archived.TapdIssueAssigneeHistory{},
		&archived.TapdIteration{},
		&archived.TapdIterationIssue{},
		&archived.TapdSource{},
		&archived.TapdBug{},
		&archived.TapdStory{},
		&archived.TapdTask{},
		&archived.TapdIssueLabel{},
	)
}

func (*InitSchemas) Version() uint64 {
	return 20220420231138
}

func (*InitSchemas) Name() string {
	return "Tapd init schemas"
}

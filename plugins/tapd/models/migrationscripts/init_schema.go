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
		&archived.TapdWorklog{},
		&archived.TapdWorkspaceIteration{},
		&archived.TapdUser{},
		&archived.TapdBugChangelog{},
		&archived.TapdBugChangelogItem{},
		&archived.TapdStoryChangelog{},
		&archived.TapdStoryChangelogItem{},
		&archived.TapdTaskChangelog{},
		&archived.TapdTaskChangelogItem{},
		&archived.TapdIssue{},
		&archived.TapdIssueSprintHistory{},
		&archived.TapdIssueStatusHistory{},
		&archived.TapdIssueAssigneeHistory{},
		&archived.TapdIteration{},
		&archived.TapdSource{},
		&archived.TapdBug{},
		&archived.TapdStory{},
		&archived.TapdTask{},
		&archived.TapdTaskLabel{},
		&archived.TapdBugLabel{},
		&archived.TapdStoryLabel{},
		&archived.TapdBugStatus{},
		&archived.TapdStoryStatus{},
		&archived.TapdBugCommit{},
		&archived.TapdStoryCommit{},
		&archived.TapdTaskCommit{},
		&archived.TapdWorkSpaceBug{},
		&archived.TapdWorkSpaceStory{},
		&archived.TapdWorkSpaceTask{},
		&archived.TapdIterationBug{},
		&archived.TapdIterationStory{},
		&archived.TapdIterationTask{},
	)
}

func (*InitSchemas) Version() uint64 {
	return 20220420231138
}

func (*InitSchemas) Name() string {
	return "Tapd init schemas"
}

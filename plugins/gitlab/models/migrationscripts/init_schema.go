package migrationscripts

import (
	"context"

	"github.com/merico-dev/lake/plugins/gitlab/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type InitSchemas struct{}

func (*InitSchemas) Up(ctx context.Context, db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		&archived.GitlabProject{},
		&archived.GitlabMergeRequest{},
		&archived.GitlabCommit{},
		&archived.GitlabTag{},
		&archived.GitlabProjectCommit{},
		&archived.GitlabPipeline{},
		&archived.GitlabReviewer{},
		&archived.GitlabMergeRequestNote{},
		&archived.GitlabMergeRequestCommit{},
		&archived.GitlabMergeRequestComment{},
		&archived.GitlabUser{},
	)
}

func (*InitSchemas) Version() uint64 {
	return 20220407201136
}

func (*InitSchemas) Name() string {
	return "Gitlab init schemas"
}

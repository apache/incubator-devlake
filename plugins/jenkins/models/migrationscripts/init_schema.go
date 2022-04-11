package migrationscripts

import (
	"context"

	"github.com/merico-dev/lake/plugins/jenkins/models/migrationscripts/archived"
	"gorm.io/gorm"
)

const (
	Owner = "Jenkins"
)

type InitSchemas struct{}

func (*InitSchemas) Up(ctx context.Context, db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		&archived.JenkinsJob{},
		&archived.JenkinsBuild{},
	)
}

func (*InitSchemas) Version() uint64 {
	return 20220407201137
}

func (*InitSchemas) Owner() string {
	return Owner
}

func (*InitSchemas) Name() string {
	return "create init schemas"
}

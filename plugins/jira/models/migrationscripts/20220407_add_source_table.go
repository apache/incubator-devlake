package migrationscripts

import (
	"context"
	"github.com/apache/incubator-devlake/plugins/jira/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type InitSchemas struct{}

func (*InitSchemas) Up(ctx context.Context, db *gorm.DB) error {
	m := db.Migrator()
	if m.HasTable(&archived.JiraConnection{}) {
		return nil
	}
	return db.Migrator().AutoMigrate(
		&archived.JiraSource{},
	)
}

func (*InitSchemas) Version() uint64 {
	return 20220407201138
}

func (*InitSchemas) Name() string {
	return "Jira init schemas"
}

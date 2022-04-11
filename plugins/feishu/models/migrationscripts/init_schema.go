package migrationscripts

import (
	"context"

	"github.com/merico-dev/lake/plugins/feishu/models/migrationscripts/archived"
	"gorm.io/gorm"
)

const (
	Owner = "Lark"
)

type InitSchemas struct{}

func (*InitSchemas) Up(ctx context.Context, db *gorm.DB) error {
	return db.Migrator().AutoMigrate(
		&archived.FeishuMeetingTopUserItem{},
	)
}

func (*InitSchemas) Version() uint64 {
	return 20220407201134
}

func (*InitSchemas) Owner() string {
	return Owner
}

func (*InitSchemas) Name() string {
	return "create init schemas"
}

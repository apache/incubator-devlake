package migrationscripts

import (
	"context"
	"github.com/merico-dev/lake/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type updateSchemas20220510 struct{}

func (*updateSchemas20220510) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().RenameColumn(archived.Note{}, "system", "is_system")
	if err != nil {
		return err
	}

	return nil
}

func (*updateSchemas20220510) Version() uint64 {
	return 20220510212399
}

func (*updateSchemas20220510) Owner() string {
	return "Framework"
}

func (*updateSchemas20220510) Name() string {
	return "Change key word system to is_system"
}

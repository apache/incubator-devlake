package migrationscripts

import (
	"context"
	"github.com/apache/incubator-devlake/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type updateSchemas20220505 struct{}

func (*updateSchemas20220505) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().RenameColumn(archived.Pipeline{}, "step", "stage")
	if err != nil {
		return err
	}
	return nil
}

func (*updateSchemas20220505) Version() uint64 {
	return 20220505212344
}

func (*updateSchemas20220505) Name() string {
	return "Rename step to stage "
}

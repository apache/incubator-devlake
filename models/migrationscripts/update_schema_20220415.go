package migrationscripts

import (
	"context"
	"gorm.io/gorm"
)

type updateSchema0415 struct{}

func (*updateSchema0415) Up(ctx context.Context, db *gorm.DB) error {
	//err := db.Migrator().AddColumn(&archived.Board{}, "test varchar(255)")
	return nil
}

func (*updateSchema0415) Version() uint64 {
	return 20220415212344
}

func (*updateSchema0415) Owner() string {
	return "Framework"
}

func (*updateSchema0415) Name() string {
	return "create init schemas"
}

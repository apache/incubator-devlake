package migrationscripts

import (
	"context"
	"github.com/apache/incubator-devlake/plugins/gitlab/models/migrationscripts/archived"

	"gorm.io/gorm"
)

type UpdateSchemas20220510 struct{}

func (*UpdateSchemas20220510) Up(ctx context.Context, db *gorm.DB) error {
	err := db.Migrator().RenameColumn(archived.GitlabMergeRequestNote{}, "system", "is_system")
	if err != nil {
		return err
	}

	return nil
}

func (*UpdateSchemas20220510) Version() uint64 {
	return 20220510212344
}

func (*UpdateSchemas20220510) Owner() string {
	return "Gitlab"
}

func (*UpdateSchemas20220510) Name() string {
	return "Change key word system to is_system"
}

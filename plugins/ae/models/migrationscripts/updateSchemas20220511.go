package migrationscripts

import (
	"context"
	"time"

	"github.com/merico-dev/lake/models/migrationscripts/archived"
	"gorm.io/gorm"
)

type AEProject20220511 struct {
	Id           int    `gorm:"primaryKey"` // The update
	GitUrl       string `gorm:"type:varchar(255);comment:url of the repo in github"`
	Priority     int
	AECreateTime *time.Time
	AEUpdateTime *time.Time
	archived.NoPKModel
}

func (AEProject20220511) TableName() string {
	return "_tool_ae_projects"
}

type UpdateSchemas20220511 struct{}

func (*UpdateSchemas20220511) Up(ctx context.Context, db *gorm.DB) error {
	// if db.Name() == "postgres" || db.Name() == "pg" || db.Name() == "postgresql" {
	// 	err := db.Exec("ALTER TABLE _tool_ae_projects ALTER COLUMN id TYPE bigint USING id::bigint").Error
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }
	err := db.Exec("ALTER TABLE _tool_ae_projects ALTER COLUMN id TYPE bigint USING id::bigint").Error
	if err != nil {
		return err
	}
	err = db.Migrator().AlterColumn(AEProject20220511{}, "id")
	if err != nil {
		return err
	}
	return nil
}

func (*UpdateSchemas20220511) Version() uint64 {
	return 20220511091919
}

func (*UpdateSchemas20220511) Name() string {
	return "Alter _tool_ae_projects column `id` to int"
}

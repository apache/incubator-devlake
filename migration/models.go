package migration

import (
	"time"
)

const (
	tableName = "_devlake_migration_history"
)

type MigrationHistory struct {
	CreatedAt     time.Time
	ScriptVersion uint64 `gorm:"primarykey"`
	ScriptName    string `gorm:"primarykey;type:varchar(255)"`
	Comment       string
}

func (MigrationHistory) TableName() string {
	return tableName
}

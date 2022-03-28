package models

import "github.com/robfig/cron/v3"

type CronEntry struct {
	EntryId     cron.EntryID
	Enable      bool
	BlueprintId uint64 `gorm:"primaryKey""`
}

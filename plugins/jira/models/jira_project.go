package models

import (
	"github.com/merico-dev/lake/models/common"
)

type JiraProject struct {
	common.NoPKModel

	// collected fields
	SourceId uint64 `gorm:"primarykey"`
	Id       string `gorm:"primaryKey"`
	Key      string
	Name     string
}

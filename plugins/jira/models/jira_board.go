package models

import (
	"github.com/merico-dev/lake/models/common"
)

type JiraBoard struct {
	common.NoPKModel
	SourceId  uint64 `gorm:"primaryKey"`
	BoardId   uint64 `gorm:"primaryKey"`
	ProjectId uint
	Name      string
	Self      string
	Type      string
}

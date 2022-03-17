package models

import "github.com/merico-dev/lake/models/common"

type JiraBoardIssue struct {
	SourceId uint64 `gorm:"primaryKey"`
	BoardId  uint64 `gorm:"primaryKey"`
	IssueId  uint64 `gorm:"primaryKey"`
	common.NoPKModel
}

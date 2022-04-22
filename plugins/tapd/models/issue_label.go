package models

import (
	"github.com/merico-dev/lake/models/common"
)

type TapdIssueLabel struct {
	IssueId   uint64 `gorm:"primaryKey;autoIncrement:false"`
	LabelName string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}

func (TapdIssueLabel) TableName() string {
	return "_tool_tapd_issue_labels"
}

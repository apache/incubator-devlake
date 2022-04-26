package models

import (
	"github.com/merico-dev/lake/models/common"
)

type TapdBugLabel struct {
	BugId     Uint64s `gorm:"primaryKey;autoIncrement:false"`
	LabelName string  `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}

func (TapdBugLabel) TableName() string {
	return "_tool_tapd_bug_labels"
}

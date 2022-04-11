package archived

import (
	"github.com/merico-dev/lake/models/common"
)

type JiraBoard struct {
	common.NoPKModel
	SourceId  uint64 `gorm:"primaryKey"`
	BoardId   uint64 `gorm:"primaryKey"`
	ProjectId uint
	Name      string `gorm:"type:varchar(255)"`
	Self      string `gorm:"type:varchar(255)"`
	Type      string `gorm:"type:varchar(100)"`
}

type JiraBoardIssue struct {
	SourceId uint64 `gorm:"primaryKey"`
	BoardId  uint64 `gorm:"primaryKey"`
	IssueId  uint64 `gorm:"primaryKey"`
	common.NoPKModel
}

func (JiraBoard) TableName() string {
	return "_tool_jira_boards"
}

func (JiraBoardIssue) TableName() string {
	return "_tool_jira_board_issues"
}

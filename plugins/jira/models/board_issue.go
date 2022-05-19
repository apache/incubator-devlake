package models

import "github.com/apache/incubator-devlake/models/common"

type JiraBoardIssue struct {
	ConnectionId uint64 `gorm:"primaryKey"`
	BoardId      uint64 `gorm:"primaryKey"`
	IssueId      uint64 `gorm:"primaryKey"`
	common.NoPKModel
}

func (JiraBoardIssue) TableName() string {
	return "_tool_jira_board_issues"
}

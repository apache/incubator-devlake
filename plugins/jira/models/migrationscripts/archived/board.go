package archived

import "github.com/merico-dev/lake/models/migrationscripts/archived"

type JiraBoard struct {
	archived.NoPKModel
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
	archived.NoPKModel
}

func (JiraBoard) TableName() string {
	return "_tool_jira_boards"
}

func (JiraBoardIssue) TableName() string {
	return "_tool_jira_board_issues"
}

package models

type JiraBoardIssue struct {
	BoardId uint64 `gorm:"primary_key"`
	IssueId uint64 `gorm:"primary_key"`
}

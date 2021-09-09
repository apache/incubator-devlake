package models

type JiraBoardIssue struct {
	BoardId uint64 `gorm:"primaryKey"`
	IssueId uint64 `gorm:"primaryKey"`
}

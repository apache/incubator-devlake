package models

type JiraBoardIssue struct {
	SourceId uint64 `gorm:"primaryKey"`
	BoardId  uint64 `gorm:"primaryKey"`
	IssueId  uint64 `gorm:"primaryKey"`
}

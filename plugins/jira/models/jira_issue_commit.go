package models

type JiraIssueCommit struct {
	SourceId  uint64 `gorm:"primaryKey"`
	IssueId   uint64 `gorm:"primaryKey"`
	CommitSha string `gorm:"primaryKey;type:char(40)"`
}

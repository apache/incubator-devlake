package models

type GitlabProjectCommit struct {
	GitlabProjectId int    `gorm:"primaryKey"`
	CommitSha       string `gorm:"primaryKey;type:char(40)"`
}

package models

type GithubRepoCommit struct {
	GithubRepoId int    `gorm:"primaryKey"`
	CommitSha    string `gorm:"primaryKey;type:char(40)"`
}

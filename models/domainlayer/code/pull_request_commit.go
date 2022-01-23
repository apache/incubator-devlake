package code

type PullRequestCommit struct {
	CommitSha     string `gorm:"primaryKey"`
	PullRequestId int    `gorm:"primaryKey;autoIncrement:false"`
}

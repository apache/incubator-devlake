package archived

type IssueRepoCommit struct {
	NoPKModel
	IssueId   string `gorm:"primaryKey;type:varchar(255)"`
	RepoUrl   string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha string `gorm:"primaryKey;type:varchar(255)"`
}

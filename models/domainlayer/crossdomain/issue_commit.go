package crossdomain

type IssueCommit struct {
	IssueId   string `gorm:"primaryKey;type:varchar(255)"`
	CommitSha string `gorm:"primaryKey;type:varchar(255)"`
}

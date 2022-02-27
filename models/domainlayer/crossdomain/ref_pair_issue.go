package crossdomain

type RefPairIssue struct {
	NewRefCommitSha string `gorm:"primaryKey;type:char(40)"`
	OldRefCommitSha string `gorm:"primaryKey;type:char(40)"`
	IssueNumber     string `gorm:"primaryKey;type:varchar(255)"`
	IssueId         string
}

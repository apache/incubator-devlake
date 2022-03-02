package crossdomain

type RefsIssuesDiffs struct {
	NewRefName      string `gorm:"type:varchar(255)"`
	OldRefName      string `gorm:"type:varchar(255)"`
	NewRefCommitSha string `gorm:"type:char(40)"`
	OldRefCommitSha string `gorm:"type:char(40)"`
	IssueNumber     string `gorm:"type:varchar(255)"`
	IssueId         string `gorm:";type:varchar(255)"`
}

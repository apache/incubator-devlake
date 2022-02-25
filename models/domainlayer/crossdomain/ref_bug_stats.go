package crossdomain

type RefBugStats struct {
	NewRefName  string `gorm:"primaryKey;type:varchar(255)"`
	OldRefName  string `gorm:"primaryKey;type:varchar(255)"`
	IssueCount  int
	IssueNumber string
}

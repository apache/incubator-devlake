package ticket

type BoardIssue struct {
	BoardId string `gorm:"primaryKey;type:varchar(255)"`
	IssueId string `gorm:"primaryKey;type:varchar(255)"`
}

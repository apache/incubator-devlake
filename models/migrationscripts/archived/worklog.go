package archived

import (
	"time"
)

type IssueWorklog struct {
	DomainEntity
	AuthorId         string `gorm:"type:varchar(255)"`
	Comment          string
	TimeSpentMinutes int
	LoggedDate       *time.Time
	StartedDate      *time.Time
	IssueId          string `gorm:"index;type:varchar(255)"`
}

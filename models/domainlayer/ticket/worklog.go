package ticket

import "time"

type Worklog struct {
	Id               string `gorm:"primaryKey"`
	AuthorId         string
	Comment          string
	TimeSpentMinutes int
	LoggedDate       time.Time
	StartedDate      time.Time
	IssueId          string `gorm:"index"`
}

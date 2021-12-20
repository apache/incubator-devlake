package ticket

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type Sprint struct {
	domainlayer.DomainEntity

	// collected fields
	BoardId      string `gorm:"index"`
	Url          string
	State        string
	Name         string
	StartDate    *time.Time
	EndDate      *time.Time
	CompleteDate *time.Time
}

type SprintIssue struct {
	SprintId  string `gorm:"primaryKey"`
	IssueId   string `gorm:"primaryKey"`
	AddedAt   *time.Time
	RemovedAt *time.Time
}

type SprintIssueBurndown struct {
	SprintId  string `gorm:"primaryKey"`
	EndedHour int    `gorm:"primaryKey"`
	StartedAt time.Time
	EndedAt   time.Time

	Added     int
	Removed   int
	Remaining int

	AddedRequirements     int
	RemovedRequirements   int
	RemainingRequirements int

	AddedBugs     int
	RemovedBugs   int
	RemainingBugs int

	AddedIncidents     int
	RemovedIncidents   int
	RemainingIncidents int

	AddedOtherIssues     int
	RemovedOtherIssues   int
	RemainingOtherIssues int

	AddedStoryPoints     int
	RemovedStoryPoints   int
	RemainingStoryPoints int
}

func (SprintIssueBurndown) TableName() string {
	return "sprint_issue_burndown"
}

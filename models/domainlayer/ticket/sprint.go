package ticket

import (
	"time"
)

type Sprint struct {
	Id            string `gorm:"primaryKey"`
	Name          string
	Url           string
	Status        string
	Title         string
	StartedDate   *time.Time
	EndedDate     *time.Time
	CompletedDate *time.Time
}

type SprintIssue struct {
	SprintId    string `gorm:"primaryKey"`
	IssueId     string `gorm:"primaryKey"`
	Status      bool
	AddedDate   *time.Time
	RemovedDate *time.Time
	AddedStage  string
}

type SprintIssueBurndown struct {
	SprintId  string `gorm:"primaryKey"`
	EndedHour   int    `gorm:"primaryKey"`
	StartedDate time.Time
	EndedDate   time.Time

	Added     int
	Removed   int
	Remaining int
	Resolved int

	AddedRequirements     int
	RemovedRequirements   int
	RemainingRequirements int
	ResolvedRequirements int

	AddedBugs     int
	RemovedBugs   int
	RemainingBugs int
	ResolvedBugs int

	AddedIncidents     int
	RemovedIncidents   int
	RemainingIncidents int
	ResolvedIncidents int

	AddedOtherIssues     int
	RemovedOtherIssues   int
	RemainingOtherIssues int
	ResolvedOtherIssues int

	AddedStoryPoints     int
	RemovedStoryPoints   int
	RemainingStoryPoints int
	ResolvedStoryPoints int
}

func (SprintIssueBurndown) TableName() string {
	return "sprint_issue_burndown"
}

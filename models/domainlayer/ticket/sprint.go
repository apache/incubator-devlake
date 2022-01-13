package ticket

import (
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

const (
	BeforeSprint = "BEFORE_SPRINT"
	DuringSprint = "DURING_SPRINT"
	AfterSprint  = "AFTER_SPRINT"
)

type Sprint struct {
	domainlayer.DomainEntity
	Name          string
	Url           string
	Status        string
	Title         string
	StartedDate   *time.Time
	EndedDate     *time.Time
	CompletedDate *time.Time
}

type SprintIssue struct {
	SprintId      string `gorm:"primaryKey"`
	IssueId       string `gorm:"primaryKey"`
	IsRemoved     bool
	AddedDate     *time.Time
	RemovedDate   *time.Time
	AddedStage    string
	ResolvedStage string
}
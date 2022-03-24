package ticket

import (
	"github.com/merico-dev/lake/models/common"
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

var (
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
	OriginBoardID string
}

type SprintIssue struct {
	common.NoPKModel
	SprintId      string `gorm:"primaryKey"`
	IssueId       string `gorm:"primaryKey"`
	IsRemoved     bool
	AddedDate     *time.Time
	RemovedDate   *time.Time
	AddedStage    *string `gorm:"type:varchar(255)"`
	ResolvedStage *string `gorm:"type:varchar(255)"`
}

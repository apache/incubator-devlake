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
	Name            string `gorm:"type:varchar(255)"`
	Url             string `gorm:"type:varchar(255)"`
	Status          string `gorm:"type:varchar(100)"`
	StartedDate     *time.Time
	EndedDate       *time.Time
	CompletedDate   *time.Time
	OriginalBoardID string `gorm:"type:varchar(255)"`
}

type SprintIssue struct {
	common.NoPKModel
	SprintId      string `gorm:"primaryKey;type:varchar(255)"`
	IssueId       string `gorm:"primaryKey;type:varchar(255)"`
	IsRemoved     bool
	AddedDate     *time.Time
	RemovedDate   *time.Time
	AddedStage    *string `gorm:"type:varchar(255)"`
	ResolvedStage *string `gorm:"type:varchar(255)"`
}

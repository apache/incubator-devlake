package archived

import (
	"time"

	"github.com/merico-dev/lake/models/common"
	"github.com/merico-dev/lake/models/domainlayer"
)

type Board struct {
	domainlayer.DomainEntity
	Name        string `gorm:"type:char(255)"`
	Description string
	Url         string `gorm:"type:char(255)"`
	CreatedDate *time.Time
}

type BoardSprint struct {
	common.NoPKModel
	BoardId  string `gorm:"primaryKey;type:varchar(255)"`
	SprintId string `gorm:"primaryKey;type:varchar(255)"`
}

type BoardIssue struct {
	BoardId string `gorm:"primaryKey;type:varchar(255)"`
	IssueId string `gorm:"primaryKey;type:varchar(255)"`
	common.NoPKModel
}

type BoardRepo struct {
	BoardId string `gorm:"primaryKey;type:varchar(255)"`
	RepoId  string `gorm:"primaryKey;type:varchar(255)"`
}

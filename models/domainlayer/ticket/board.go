package ticket

import (
	"github.com/merico-dev/lake/models/common"
	"time"

	"github.com/merico-dev/lake/models/domainlayer"
)

type Board struct {
	domainlayer.DomainEntity
	Name        string
	Description string
	Url         string
	CreatedDate *time.Time
}

type BoardSprint struct {
	common.NoPKModel
	BoardId  string `gorm:"primaryKey"`
	SprintId string `gorm:"primaryKey"`
}

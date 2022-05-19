package ticket

import (
	"github.com/apache/incubator-devlake/models/common"
	"time"

	"github.com/apache/incubator-devlake/models/domainlayer"
)

type Board struct {
	domainlayer.DomainEntity
	Name        string `gorm:"type:varchar(255)"`
	Description string
	Url         string `gorm:"type:varchar(255)"`
	CreatedDate *time.Time
}

type BoardSprint struct {
	common.NoPKModel
	BoardId  string `gorm:"primaryKey;type:varchar(255)"`
	SprintId string `gorm:"primaryKey;type:varchar(255)"`
}

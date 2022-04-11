package archived

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type AEProject struct {
	Id           int    `gorm:"primaryKey;type:varchar(255)"`
	GitUrl       string `gorm:"type:varchar(255);comment:url of the repo in github"`
	Priority     int
	AECreateTime *time.Time
	AEUpdateTime *time.Time
	common.NoPKModel
}

func (AEProject) TableName() string {
	return "_tool_ae_projects"
}

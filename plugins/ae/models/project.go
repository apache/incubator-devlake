package models

import (
	"time"

	"github.com/merico-dev/lake/models/common"
)

type AEProject struct {
	Id           int    `gorm:"primaryKey"`
	GitUrl       string `gorm:"type:varchar(255);comment:url of the repo in github"`
	Priority     int
	AECreateTime *time.Time
	AEUpdateTime *time.Time
	common.NoPKModel
}

func (AEProject) TableName() string {
	return "_tool_ae_projects"
}

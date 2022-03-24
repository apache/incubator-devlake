package models

import (
	"github.com/merico-dev/lake/models/common"
	"time"
)

type AEProject struct {
	Id           int    `gorm:"primaryKey"`
	GitUrl       string `gorm:"comment:url of the repo in github"`
	Priority     int
	AECreateTime *time.Time
	AEUpdateTime *time.Time
	common.NoPKModel
}

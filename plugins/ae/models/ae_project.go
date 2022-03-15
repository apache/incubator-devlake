package models

import (
	"time"

	"github.com/merico-dev/lake/plugins/helper"
)

type AEProject struct {
	Id           int    `gorm:"primaryKey"`
	GitUrl       string `gorm:"comment:url of the repo in github"`
	Priority     int
	AECreateTime *time.Time
	AEUpdateTime *time.Time

	helper.RawDataOrigin
}

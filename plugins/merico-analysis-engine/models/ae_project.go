package models

import (
	"time"

	"github.com/merico-dev/lake/models"
)

type AEProject struct {
	Id           int    `gorm:"primaryKey"`
	GitUrl       string `gorm:"comment:url of the repo in github"`
	Priority     int
	AECreateTime *time.Time
	AEUpdateTime *time.Time

	models.NoPKModel
}

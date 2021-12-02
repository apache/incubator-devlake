package models

import (
	"time"
)

type AEProject struct {
	Id           int    `gorm:"primaryKey"`
	GitUrl       string `gorm:"comment:url of the repo in github"`
	Priority     int
	AECreateTime *time.Time
	AEUpdateTime *time.Time
}

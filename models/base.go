package models

import "time"

type Model struct {
	ID        uint64    `gorm:"primary_key"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time `gorm:"<-:update"`
}

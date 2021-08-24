package models

import "time"

type Model struct {
	ID        uint `gorm:"primary_key"`
	Status    int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
